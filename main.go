package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	T "johtotimes.com/templates"

	"johtotimes.com/internal/globals"
	"johtotimes.com/internal/handler"
	"johtotimes.com/internal/types"
	"johtotimes.com/internal/utils"
)

func listPage(slug string) templ.Component {
	log.Println("List page")
	entries, err := os.ReadDir("web/posts")
	if err != nil {
		log.Fatalln(err)
	}

	posts := []types.Post{}
	for _, e := range entries {
		post := utils.ParseMarkdown("web/posts/" + e.Name())
		category := post.Metadata["Category"].(string)
		if strings.ToLower(category) == slug {
			posts = append(posts, post)
		}
	}
	list := T.List("Category: "+slug, posts)

	return T.Base("Johto Times", list)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello posts")
	fmt.Println(r.URL.Path)
	fileName := "web" + r.URL.Path + ".md"
	fmt.Println(fileName)

	singlePage(fileName).Render(r.Context(), w)
}

func singlePage(fileName string) templ.Component {

	post := utils.ParseMarkdown(fileName)

	return T.Base(post.Metadata["Title"].(string), post.Contents)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening on port " + port)

	assets := http.FileServer(http.Dir(globals.AssetPath))

	mux := http.NewServeMux()
	// mux.Handle("/", templ.Handler(singlePage()))
	mux.Handle("/", templ.Handler(handler.IndexPage()))
	mux.HandleFunc("/posts/", postHandler)

	for _, category := range globals.Categories {
		slug := strings.ToLower(category.Name)
		mux.Handle("/"+slug, templ.Handler(listPage(slug)))
	}

	prefix := "/" + globals.AssetPath + "/"
	mux.Handle(prefix, http.StripPrefix(prefix, assets))
	http.ListenAndServe(":"+port, mux)
}
