package handler

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/database"
	"johtotimes.com/src/templates"
)

func IssuesHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	listPage("Issues", 'I', page).Render(req.Context(), w)
}

func PostsHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	listPage("", 'P', page).Render(req.Context(), w)
}

func MailbagHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	listPage("Mailbag", 'M', page).Render(req.Context(), w)
}

func NewsHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	listPage("News", 'N', page).Render(req.Context(), w)
}

func listPage(title string, postType byte, page int) templ.Component {
	log.Printf("Rendering content of type " + string(postType))
	db := database.Connect()
	defer db.Close()
	assert.NotNil(page, "Page number cannot be nil")
	posts := db.Posts.GetPage(postType, page, 10)

	log.Printf("Found %d posts", len(posts))
	description := renderHTML("")

	list := templates.ListTemplate(title, title, description, posts)

	return templates.Base("Johto Times", list)
}
