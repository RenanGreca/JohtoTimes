package handler

import (
	"log"
	"net/http"

	"johtotimes.com/src/templates"
)

func ArchiveHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	body := templates.ArchiveTemplate("Archive")
	render(body, isHTMX(req), "Archive", w)
}
