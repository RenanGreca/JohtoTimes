package post

import (
	"log"
	"os"
	"time"
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

func GetFromDirectory(postsDir string) []Post {
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
