package post

import (
	"log"
	"os"
	"time"

	"johtotimes.com/src/internal"
)

type Post struct {
	ID          int64
	Title       string
	String      string
	Slug        string
	Category    string
	Tags        []string
	Img         string
	Description string
	Date        time.Time
}

func getFromDirectory(postsDir string) []Post {
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatalln(err)
	}

	var posts []Post
	for _, e := range entries {
		post, _ := parseHeaders(postsDir + "/" + e.Name())
		// posts[i] = post
		posts = append(posts, post)
	}
	return posts
}

// Received the path to a markdown file and returns a Post element
func parseHeaders(fileName string) (Post, string) {
	md := internal.ReadFile(fileName)

	metadata, buf := internal.ParseMarkdown(md)

	return Post{
		// Contents: content,
		Slug:  internal.ExtractSlug(fileName),
		Title: metadata["Title"].(string),
		// Category:    metadata["Category"].(string),
		Img:         metadata["Header"].(string),
		Description: metadata["Description"].(string),
		Date:        internal.ExtractDate(fileName),
		// Tags:        extractTags(metadata),
	}, buf.String()
}
