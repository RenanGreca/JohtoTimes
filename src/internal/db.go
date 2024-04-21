package internal

import (
	"log"
	"os"

	"github.com/a-h/templ"
)

type DB struct {
	Posts      []*Post
	Issues     []*Issue
	Categories []Category
	Tags       []Category
}

type Post struct {
	Contents    templ.Component
	String      string
	Slug        string
	Title       string
	Category    string
	Tags        []string
	Img         string
	Description string
}

type Issue struct {
	Volume      int
	Issue       int
	Title       string
	News        templ.Component
	Post        *Post
	Mailbag     templ.Component
	Description string
}

type Category struct {
	Name   string
	Plural string
	Slug   string
	Posts  []*Post
}

// Returns list of posts that match given category slug
func (db *DB) PostsInCategory(categorySlug string) []*Post {
	posts := []*Post{}

	for _, post := range db.Posts {
		if post.Category == categorySlug {
			posts = append(posts, post)
		}
	}

	return posts
}

func (db *DB) FillDB(postsPath string) {
	entries, err := os.ReadDir(postsPath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, e := range entries {
		post := ParseMarkdown(postsPath + e.Name())
		db.Posts = append(db.Posts, &post)
	}
}

func (db *DB) GetIssues(issuesPath string) {
	// entries, err := os.ReadDir(issuesPath)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, e := range entries {
	// parse issue markdown
	// search for news file
	// search for post file
	// search for mailbag file
	// }
}
