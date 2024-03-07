package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/internal"
	T "johtotimes.com/templates"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello posts")
	fmt.Println(r.URL.Path)
	fileName := "web" + r.URL.Path + ".md"
	fmt.Println(fileName)

	singlePage(fileName).Render(r.Context(), w)
}

func singlePage(fileName string) templ.Component {

	post := internal.ParseMarkdown(fileName)

	return T.Base(post.Title, post.Contents)
}
