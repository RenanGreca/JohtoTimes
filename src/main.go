package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/go-mail/mail"

	"johtotimes.com/src/handler"
	"johtotimes.com/src/internal"
	T "johtotimes.com/src/templates"
)

func main() {
	// emailSender()
	httpHandler()
}

func httpHandler() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Listening on port " + port)

	mux := http.NewServeMux()
	// mux.Handle("/", templ.Handler(singlePage()))
	mux.Handle("/", templ.Handler(handler.IndexPage()))
	prefix := "/" + internal.AssetPath + "/"
	assets := http.FileServer(http.Dir(internal.AssetPath))
	fmt.Println(prefix)
	mux.Handle("GET "+prefix, http.StripPrefix(prefix, assets))
	// for _, dir := range [...]string{"fonts", "img", "scripts", "styles"} {
	// 	pattern := prefix + dir
	// 	fmt.Println(pattern)
	// 	mux.Handle(pattern, http.StripPrefix(prefix, assets))
	// }

	mux.HandleFunc("GET /category/{category}", handler.ListHandler)
	mux.HandleFunc("GET /posts/{category}/{slug}", handler.PostHandler)
	// for _, category := range internal.Categories {
	// 	slug := strings.ToLower(category.Name)
	// 	mux.Handle("/"+slug, templ.Handler(handler.ListPage(slug)))
	// }

	http.ListenAndServe(":"+port, mux)
}

func emailSender() {
	fileName := "./web/posts/2024-02-22-pokemon-legends-celebi-a-concept.md"
	post := internal.ParseMarkdown(fileName)
	f, err := os.Create("head.html")
	if err != nil {
		log.Fatalln(err)
	}
	err = T.Head(post.Title).Render(context.Background(), f)
	if err != nil {
		log.Fatalln(err)
	}
	head := internal.ReadFile("head.html")

	pass := os.Getenv("GOEMAILPASS")

	m := mail.NewMessage()
	m.SetHeader("From", "newsletter@johtotimes.com")
	m.SetHeader("To", "renangreca@icloud.com")
	m.SetHeader("Subject", post.Title)
	m.SetBody("text/html", head+post.String)

	d := mail.NewDialer("smtp.gmail.com", 587, "renangreca@gmail.com", pass)
	if err := d.DialAndSend(m); err != nil {
		log.Fatalln(err)
	}
}
