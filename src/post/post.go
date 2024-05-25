package post

import (
	"log"
	"os"
	"time"

	"johtotimes.com/src/category"
)

type Post struct {
	ID          int64
	Title       string
	String      string
	Slug        string
	FileName    string
	Category    *category.Category
	Tags        []*category.Category
	Img         string
	Description string
	Type        byte // 'P' for post, 'N' for news, 'M' for mailbag
	Date        time.Time
}

func getFromDirectory(postsDir string) []Markdown {
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatalln(err)
	}

	var posts []Markdown
	for _, e := range entries {
		fileName := postsDir + "/" + e.Name()
		post := ParseHeaders(fileName)
		// posts[i] = post
		posts = append(posts, post)
	}
	return posts
}
