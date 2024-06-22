package handler

import (
	"fmt"

	"github.com/a-h/templ"
	"johtotimes.com/src/templates"
)

func errorPage(code int) templ.Component {
	errorBody := templates.Error(code)
	title := fmt.Sprintf("Error %d", code)
	return templates.Base(title, errorBody)
}
