package handler

import (
	"log"
	"os"
	"strings"

	"github.com/a-h/templ"
	"johtotimes.com/internal"
	T "johtotimes.com/templates"
)

func ListPage(slug string) templ.Component {
	log.Println("List page: ", slug)
	entries, err := os.ReadDir("web/posts")
	if err != nil {
		log.Fatalln(err)
	}

	posts := []internal.Post{}
	for _, e := range entries {
		post := internal.ParseMarkdown("web/posts/" + e.Name())
		if strings.ToLower(post.Category) == slug {
			posts = append(posts, post)
		}
	}
	list := T.List("Category: "+slug, posts)

	return T.Base("Johto Times", list)
}
