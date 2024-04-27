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

func Populate() []Post {
	entries, err := os.ReadDir("web/posts")
	if err != nil {
		log.Fatalln(err)
	}

	var posts []Post
	for _, e := range entries {
		post, _ := parseHeaders("web/posts/" + e.Name())
		// posts[i] = post
		posts = append(posts, post)
	}
	return posts
}
