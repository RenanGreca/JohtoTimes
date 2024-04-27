package post

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"johtotimes.com/src/internal"
	T "johtotimes.com/src/templates"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("Hello posts")
	// fmt.Println(r.URL.Path)
	slug := req.PathValue("slug")
	fileName := internal.PostsPath + "/" + slug + ".md"
	fmt.Println(fileName)

	if _, err := os.Stat(fileName); err == nil {
		singlePage(fileName).Render(req.Context(), w)
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Error trying to open file " + fileName)
	}

}

func singlePage(fileName string) templ.Component {

	post, contents := parseHeaders(fileName)

	return T.Base(post.Title, unsafe(contents))
}

// ! Todo: Move to another package
func unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
