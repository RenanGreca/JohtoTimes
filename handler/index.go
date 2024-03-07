package handler

import (
	"log"
	"os"

	"github.com/a-h/templ"
	"johtotimes.com/internal"
	T "johtotimes.com/templates"
)

func IndexPage() templ.Component {
	// log.Println("Index page")
	entries, err := os.ReadDir("web/posts")
	if err != nil {
		log.Fatalln(err)
	}

	// for _, e := range entries {
	posts := []internal.Post{}

	for _, e := range entries {
		post := internal.ParseMarkdown("web/posts/" + e.Name())
		// posts[i] = post
		posts = append(posts, post)
	}
	list := T.List("", posts)

	return T.Base("Johto Times", list)

	// }
}
