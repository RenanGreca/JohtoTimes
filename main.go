package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	T "johtotimes.com/templates"

	Globals "johtotimes.com/internal/globals"
	Types "johtotimes.com/internal/types"
)

func indexPage() templ.Component {
	log.Println("Index page")
	entries, err := os.ReadDir("web/posts")
	if err != nil {
		log.Fatalln(err)
	}

	// for _, e := range entries {
	posts := []Types.Post{}

	for _, e := range entries {
		post := parseMarkdown("web/posts/" + e.Name())
		// posts[i] = post
		posts = append(posts, post)
	}
	list := T.List(posts)

	return T.Base("Johto Times", list)

	// }
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello posts")
	fmt.Println(r.URL.Path)
	fileName := "web" + r.URL.Path + ".md"
	fmt.Println(fileName)

	singlePage(fileName).Render(r.Context(), w)
}

func singlePage(fileName string) templ.Component {

	post := parseMarkdown(fileName)

	return T.Base(post.Metadata["Title"].(string), post.Contents)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening on port " + port)

	assets := http.FileServer(http.Dir(Globals.AssetPath))

	mux := http.NewServeMux()
	// mux.Handle("/", templ.Handler(singlePage()))
	mux.Handle("/", templ.Handler(indexPage()))
	mux.HandleFunc("/posts/", postHandler)
	prefix := "/" + Globals.AssetPath + "/"
	mux.Handle(prefix, http.StripPrefix(prefix, assets))
	http.ListenAndServe(":"+port, mux)
}

func readFile(fileName string) string {
	log.Println("Opening file: " + fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

func parseMarkdown(fileName string) Types.Post {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	md := readFile(fileName)
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, ".")
	slug := split2[0]
	log.Println("Post slug: " + slug)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(md), &buf, parser.WithContext(context)); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	content := unsafe(buf.String())
	metadata := meta.Get(context)
	title := metadata["Title"]
	fmt.Println("Title: ", title)

	return Types.Post{
		Contents: content,
		Metadata: metadata,
		Slug:     slug,
	}
}

func unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
