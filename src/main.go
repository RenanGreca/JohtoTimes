package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"johtotimes.com/src/database"
	"johtotimes.com/src/handler"
	"johtotimes.com/src/internal"
)

func main() {
	database.NewDB()

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

	// Index page and tab bar items
	mux.HandleFunc("GET /", handler.PostsHandler)
	mux.HandleFunc("GET /archive", handler.ArchiveHandler)
	// mux.HandleFunc("GET /search", search.Handler)
	// mux.HandleFunc("GET /search/{query}", search.Handler)
	// mux.HandleFunc("GET /community", community.Handler)
	// mux.HandleFunc("GET /about", about.Handler)

	// Assets directory
	prefix := "/" + internal.AssetPath + "/"
	assets := http.FileServer(http.Dir(internal.AssetPath))
	mux.Handle("GET "+prefix, http.StripPrefix(prefix, assets))

	// Category/type lists
	mux.HandleFunc("GET /category/{category}", handler.CategoryHandler)
	mux.HandleFunc("GET /category/{category}/{page}", handler.CategoryHandler)
	mux.HandleFunc("GET /mailbag", handler.MailbagHandler)
	mux.HandleFunc("GET /mailbag/{page}", handler.MailbagHandler)
	mux.HandleFunc("GET /news", handler.NewsHandler)
	mux.HandleFunc("GET /news/{page}", handler.NewsHandler)

	// Post pages
	// Handle direct link to post
	mux.HandleFunc("GET /posts/{slug}", handler.PostHandler)
	// Handle category-type link
	mux.HandleFunc("GET /posts/{category}/{slug}", handler.PostHandler)

	// // Handle direct link to issue
	// mux.HandleFunc("GET /issue/{slug}", issue.Handler)

	http.ListenAndServe(":"+port, mux)
}

func emailSender() {
	// fileName := "./web/posts/2024-02-22-pokemon-legends-celebi-a-concept.md"
	// post := internal.ParseMarkdown(fileName)
	// f, err := os.Create("head.html")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// err = T.Head(post.Title).Render(context.Background(), f)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// head := internal.ReadFile("head.html")
	//
	// pass := os.Getenv("GOEMAILPASS")
	//
	// m := mail.NewMessage()
	// m.SetHeader("From", "newsletter@johtotimes.com")
	// m.SetHeader("To", "renangreca@icloud.com")
	// m.SetHeader("Subject", post.Title)
	// m.SetBody("text/html", head+post.String)
	//
	// d := mail.NewDialer("smtp.gmail.com", 587, "renangreca@gmail.com", pass)
	// if err := d.DialAndSend(m); err != nil {
	// 	log.Fatalln(err)
	// }
}
