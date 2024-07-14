package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/templates"
)

// getPageNumber extracts the page number from the GET request
func getPageNumber(req *http.Request) int {
	pagestr := req.URL.Query().Get("page")
	if len(pagestr) > 0 {
		page, err := strconv.Atoi(pagestr)
		assert.NoError(err, "Handler: Error converting %q to int\n", req.URL.Query().Get("page"))
		return page
	}
	return 0
}

// renderHTML produces the templ component from raw HTML
func renderHTML(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

// render renders the component to the response writer
// If isHTMX is false, it renders the component to the base template.
func render(component templ.Component, isHTMX bool, title string, w http.ResponseWriter) {
	if isHTMX {
		log.Println("Rendering HTMX")
		component.Render(context.Background(), http.ResponseWriter(w))
	} else {
		log.Println("Rendering HTML")
		templates.Base(title, component).Render(context.Background(), http.ResponseWriter(w))
	}
}

// isHTMX returns true if the request contains "htmx=true"
func isHTMX(req *http.Request) bool {
	htmx, exists := req.Header["Hx-Request"]
	if exists && len(htmx) > 0 {
		return htmx[0] == "true"
	}
	return false
}
