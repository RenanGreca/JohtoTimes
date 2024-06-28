package model

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"time"

	"github.com/gosimple/slug"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/file"
	"johtotimes.com/src/markdown"
)

type Post struct {
	ID int64
	// Title is the post's title
	Title string
	// Slug is the slug seen in the post's URL
	Slug string
	// FileName is the path of the Markdown file that originated this post.
	FileName string
	// Category is the category this post belongs to.
	Category Category
	// Tags are the post's associated tags.
	Tags []Category
	// Img is the preview image used for lists.
	Img string
	// Description is the summary that appears below the title.
	Description string
	// Issue represents the year the post appeared in.
	Issue int
	// Volume represents which week the post appeared in.
	Volume int
	// The permalink URL to the post.
	Permalink string
	// Type is the type of post:
	// 'P' for post, 'N' for news, 'M' for mailbag, 'I' for issue
	Type byte
	// CreatedAt is the date the post was originally published (extracted from the filename).
	CreatedAt time.Time
	// ModifiedAt is the date the post was modified, if ever.
	ModifiedAt time.Time
	// Hash is the MD5 hash of the post's contents, used for tracking changes.
	Hash string
}

func GetPostsFromDirectory(postsDir string) []Post {
	entries, err := os.ReadDir(postsDir)
	assert.NoError(err, "PostRepository: Error reading directory: %s", postsDir)

	var posts []Post
	for _, e := range entries {
		fileName := postsDir + "/" + e.Name()
		post := NewPostFromMarkdown(fileName)
		posts = append(posts, post)
	}
	return posts
}

/*
NewPostFromMarkdown receives the path to a Markdown file and returns a Post element.

	The file contents and header are parsed and placed in the Post struct.
	Initially, the Post's Category and Tags contain only their respective Slugs.
*/
func NewPostFromMarkdown(fileName string) Post {
	md := file.ReadFile(fileName)

	metadata, buf := markdown.ParseMarkdown(md)
	hash := md5.Sum([]byte(buf.String()))

	post := Post{
		FileName:  fileName,
		Slug:      markdown.ExtractSlug(fileName),
		CreatedAt: markdown.ExtractDate(fileName),
		Hash:      hex.EncodeToString(hash[:]),
	}
	post.extractMetadata(metadata)

	return post
}

func (post *Post) Content() string {
	md := file.ReadFile(post.FileName)

	_, buf := markdown.ParseMarkdown(md)

	return buf.String()
}

func (post *Post) extractMetadata(metadata map[string]interface{}) {
	assert.NotNil(metadata, "PostRepository: Error extracting metadata")

	if metadata["Title"] != nil {
		post.Title = metadata["Title"].(string)
	}

	if metadata["Header"] != nil {
		post.Img = metadata["Header"].(string)
	}

	if metadata["Category"] != nil {
		post.Category = Category{
			Singular: metadata["Category"].(string),
			Slug:     slug.Make(metadata["Category"].(string)),
		}
	}

	if metadata["Description"] != nil {
		post.Description = metadata["Description"].(string)
	}

	if metadata["Tags"] != nil {
		post.Tags = extractTags(metadata["Tags"].([]interface{}))
	}

	if metadata["Issue"] != nil {
		post.Issue = int(metadata["Issue"].(int))
	}

	if metadata["Volume"] != nil {
		post.Volume = int(metadata["Volume"].(int))
	}
}

func (post *Post) SetPermalink() {
	assert.NotNil(post.Type, "PostRepository: Post Type must be defined before Permalink.")

	switch post.Type {
	case 'P':
		assert.NotNil(post.Category.Slug, "PostRepository: Post Category must be defined before Permalink.")
		post.Permalink = "/posts/" + post.Category.Slug + "/" + post.Slug
	case 'M':
		post.Permalink = "/mailbag/" + post.Slug
	case 'N':
		post.Permalink = "/news/" + post.Slug
	case 'I':
		post.Permalink = "/issues/" + post.Slug
	}

	assert.NotNil(post.Permalink, "PostRepository: Error creating permalink")
}
