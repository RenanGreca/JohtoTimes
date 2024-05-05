package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"

	"johtotimes.com/src/database"
	"johtotimes.com/src/handler"
	"johtotimes.com/src/internal"
	"johtotimes.com/src/list"
	"johtotimes.com/src/post"
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

	// Index page
	mux.Handle("/", templ.Handler(handler.IndexPage()))

	// Assets directory
	prefix := "/" + internal.AssetPath + "/"
	assets := http.FileServer(http.Dir(internal.AssetPath))
	mux.Handle("GET "+prefix, http.StripPrefix(prefix, assets))

	// for _, dir := range [...]string{"fonts", "img", "scripts", "styles"} {
	// 	pattern := prefix + dir
	// 	fmt.Println(pattern)
	// 	mux.Handle(pattern, http.StripPrefix(prefix, assets))
	// }

	// Category lists
	mux.HandleFunc("GET /category/{category}", list.Handler)

	// Post pages
	// Handle direct link to post
	mux.HandleFunc("GET /posts/{slug}", post.Handler)
	// Handle category-type link
	mux.HandleFunc("GET /posts/{category}/{slug}", post.Handler)
	// for _, category := range internal.Categories {
	// 	slug := strings.ToLower(category.Name)
	// 	mux.Handle("/"+slug, templ.Handler(handler.ListPage(slug)))
	// }

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
