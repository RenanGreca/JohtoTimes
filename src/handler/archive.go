package handler

import (
	"net/http"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/database"
	"johtotimes.com/src/templates"
)

func ArchiveHandler(w http.ResponseWriter, req *http.Request) {
	assert.LogDebug("Handling request to " + req.URL.Path)
	body := templates.ArchiveTemplate("Archive")
	render(body, isHTMX(req), "Archive", w)
}

func ArchiveIssuesHandler(w http.ResponseWriter, req *http.Request) {
	db := database.Connect()
	defer db.Close()

	posts := db.Posts.GetPage('I', 1, 1000)
	assert.LogDebug("Found %d issues", len(posts))

	body := templates.ArchivePostsTemplate("Issues Archive", posts)

	render(body, isHTMX(req), "Archive", w)
}
