package handler

import (
	"log"
	"net/http"

	"johtotimes.com/src/templates"
)

func ArchiveHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling request to " + req.URL.Path)
	body := templates.ArchiveTemplate("Archive")

	templates.Base("Archive", body).Render(req.Context(), w)
}
