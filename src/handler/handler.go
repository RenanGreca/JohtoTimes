package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

// Extracts the page number from the request URL
func getPageNumber(req *http.Request) int {
	pagestr := req.PathValue("page")
	if len(pagestr) > 0 {
		page, err := strconv.Atoi(pagestr)
		if err != nil {
			log.Fatalf("Error converting %q to int\n", req.PathValue("page"))
		}
		return page
	}
	return 0
}

// Produces the templ component from raw HTML
func renderHTML(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
