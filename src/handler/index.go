package handler

import (
	"github.com/a-h/templ"
	"johtotimes.com/src/list"
	"johtotimes.com/src/post"
	T "johtotimes.com/src/templates"
)

func IndexPage() templ.Component {
	// log.Println("Index page")
	// entries, err := os.ReadDir("web/posts")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, e := range entries {
	posts := []post.Post{}

	// for _, e := range entries {
	// post := internal.ParseMarkdown("web/posts/" + e.Name())
	// // posts[i] = post
	// posts = append(posts, post)
	// }
	list := list.List("", posts)

	return T.Base("Johto Times", list)

	// }
}
