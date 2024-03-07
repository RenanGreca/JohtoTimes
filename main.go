package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"

	"johtotimes.com/handler"
	"johtotimes.com/internal"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening on port " + port)

	assets := http.FileServer(http.Dir(internal.AssetPath))

	mux := http.NewServeMux()
	// mux.Handle("/", templ.Handler(singlePage()))
	mux.Handle("/", templ.Handler(handler.IndexPage()))
	mux.HandleFunc("/posts/", handler.PostHandler)

	for _, category := range internal.Categories {
		slug := strings.ToLower(category.Name)
		mux.Handle("/"+slug, templ.Handler(handler.ListPage(slug)))
	}

	prefix := "/" + internal.AssetPath + "/"
	mux.Handle(prefix, http.StripPrefix(prefix, assets))
	http.ListenAndServe(":"+port, mux)
}
