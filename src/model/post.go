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

	// assert.NotNil(metadata["Title"],
	// 	"PostRepository: Error extracting Title from markdown metadata")
	if metadata["Title"] != nil {
		post.Title = metadata["Title"].(string)
	}

	assert.NotNil(metadata["Header"],
		"PostRepository: Error extracting Header from markdown metadata")
	post.Img = metadata["Header"].(string)

	assert.NotNil(metadata["Category"],
		"PostRepository: Error extracting Category from markdown metadata")
	post.Category = Category{
		Singular: metadata["Category"].(string),
		Slug:     slug.Make(metadata["Category"].(string)),
	}

	assert.NotNil(metadata["Description"],
		"PostRepository: Error extracting Description from markdown metadata")
	post.Description = metadata["Description"].(string)

	assert.NotNil(metadata["Tags"],
		"PostRepository: Error extracting Tags from markdown metadata")
	post.Tags = extractTags(metadata["Tags"].([]interface{}))

	assert.NotNil(metadata["Issue"],
		"PostRepository: Error extracting Issue from markdown metadata")
	post.Issue = int(metadata["Issue"].(int))

	assert.NotNil(metadata["Volume"],
		"PostRepository: Error extracting Volume from markdown metadata")
	post.Volume = int(metadata["Volume"].(int))
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
