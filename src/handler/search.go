package handler

import (
	"log"
	"net/http"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/database"
	"johtotimes.com/src/model"
	"johtotimes.com/src/templates"
)

func SearchHandler(w http.ResponseWriter, req *http.Request) {
	htmx := isHTMX(req)
	query := req.PathValue("query")
	searchTemplate := templates.SearchTemplate(query)
	render(searchTemplate, htmx, "Search", w)
	return
}

func SearchResultsHandler(w http.ResponseWriter, req *http.Request) {
	htmx := isHTMX(req)
	query := req.FormValue("query")
	posts := search(query, 0)
	resultsTemplate := templates.SearchResultsTemplate(posts)
	render(resultsTemplate, htmx, "Search", w)
	return
}

func search(query string, page int) []model.Post {
	log.Printf("Searching for %s", query)
	db := database.Connect()
	defer db.Close()
	assert.NotNil(query, "SearchHandler: query cannot be nil")

	posts := db.Posts.Search(query, page, 10)

	log.Printf("Found %d posts", len(posts))

	return posts
}
