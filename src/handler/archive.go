package handler

import (
	"net/http"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/templates"
)

func ArchiveHandler(w http.ResponseWriter, req *http.Request) {
	assert.LogDebug("Handling request to " + req.URL.Path)
	body := templates.ArchiveTemplate("Archive")
	render(body, isHTMX(req), "Archive", w)
}
