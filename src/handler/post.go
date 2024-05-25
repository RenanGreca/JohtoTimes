package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"johtotimes.com/src/database"
	"johtotimes.com/src/internal"
	"johtotimes.com/src/post"
	"johtotimes.com/src/templates"
)

func PostHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	slug := req.PathValue("slug")
	db := database.Connect()
	defer db.Close()
	post, err := db.Posts.GetBySlug(slug)
	if err != nil {
		log.Fatal(err)
	}
	filePath := internal.PostsPath
	if post.Type == 'M' {
		filePath = internal.MailbagPath
	}
	fileName := filePath + "/" + slug + ".md"
	fmt.Println(fileName)

	if _, err := os.Stat(fileName); err == nil {
		postPage(post, fileName).Render(req.Context(), w)
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Error trying to open file " + fileName)
	}

}

func postPage(p *post.Post, fileName string) templ.Component {
	headers := post.ParseHeaders(fileName)
	postBody := templates.PostTemplate(p, unsafe(headers.Contents))

	return templates.Base(p.Title, postBody)
}
