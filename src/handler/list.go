package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/database"
	"johtotimes.com/src/templates"
)

func IssuesHandler(w http.ResponseWriter, req *http.Request) {
	assert.LogDebug("IssuesHandler: Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	assert.LogDebug("Page number is %d\n", page)
	body := listPage("Issues", 'I', page)
	render(body, isHTMX(req), "Issues", w)
}

func PostsHandler(w http.ResponseWriter, req *http.Request) {
	assert.LogDebug("PostsHandler: Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	body := listPage("", 'P', page)
	render(body, isHTMX(req), "Posts", w)
}

func MailbagHandler(w http.ResponseWriter, req *http.Request) {
	assert.LogDebug("MailbagHandler: Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	body := listPage("Mailbag", 'M', page)
	render(body, isHTMX(req), "Mailbag", w)
}

func NewsHandler(w http.ResponseWriter, req *http.Request) {
	assert.LogDebug("NewsHandler: Handling request to " + req.URL.Path)
	page := getPageNumber(req)
	body := listPage("News", 'N', page)
	render(body, isHTMX(req), "News", w)
}

func listPage(title string, postType byte, page int) templ.Component {
	assert.LogDebug("Rendering content of type " + string(postType))
	db := database.Connect()
	defer db.Close()
	assert.NotNil(page, "Page number cannot be nil")
	posts := db.Posts.GetPage(postType, page, 10)

	assert.LogDebug("Found %d posts", len(posts))
	description := renderHTML("")

	return templates.ListTemplate(title, title, description, posts, page)
}
