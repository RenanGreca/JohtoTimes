package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/database"
	"johtotimes.com/src/internal"
	"johtotimes.com/src/post"
	"johtotimes.com/src/templates"
)

// PostHandler handles GET requests to /posts/{slug}
func PostHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling Post request to " + req.URL.Path)
	slug := req.PathValue("slug")
	db := database.Connect()
	defer db.Close()
	post, err := db.Posts.GetBySlug(slug, 'P')
	if err != nil {
		log.Printf("PostHandler: Error getting post: %s", err)
		errorPage(404).Render(req.Context(), w)
		return
	}
	render(postPage(post), isHTMX(req), post.Title, w)
}

// IssueHandler handles GET requests to /issues/{slug}.
// This involves building an issue based on Issue, Post, News and Mailbag.
func IssueHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling Issue request to " + req.URL.Path)
	slug := req.PathValue("slug")
	db := database.Connect()
	defer db.Close()
	issue, err := db.Posts.GetBySlug(slug, 'I')
	if err != nil {
		log.Printf("IssueHandler: Error getting issue: %s", err)
		errorPage(404).Render(req.Context(), w)
		return
	}
	post, err := db.Posts.GetByDateAndType(issue.Date, 'P')
	if err != nil {
		log.Printf("IssueHandler: Error getting post: %s", err)
		errorPage(404).Render(req.Context(), w)
		return
	}
	news, err := db.Posts.GetByDateAndType(issue.Date, 'N')
	assert.NoError(err, "IssueHandler: Error getting news")
	mailbag, err := db.Posts.GetByDateAndType(issue.Date, 'M')
	assert.NoError(err, "IssueHandler: Error getting mailbag")

	render(issuePage(issue, post, news, mailbag), isHTMX(req), issue.Title, w)
}

func postPage(p *post.Post) templ.Component {
	filePath := internal.PostTypePath[p.Type]
	fileName := filePath + "/" + p.Slug + ".md"
	fmt.Printf("File path: %s %s\n", string(p.Type), fileName)

	headers := post.ParseHeaders(fileName)
	return templates.SingleTemplate(p, renderHTML(headers.Contents))
}

func postBody(p *post.Post) templ.Component {
	if p == nil {
		return errorPage(404)
	}
	filePath := internal.PostTypePath[p.Type]
	fileName := filePath + "/" + p.Slug + ".md"
	fmt.Printf("File path: %s %s\n", string(p.Type), fileName)

	// In the issue builder, if the file doesn't exist, we show a message
	if !internal.FileExists(fileName) {
		switch p.Type {
		case 'M':
			return templates.MailbagTemplate(p, renderHTML("No mailbag this week!"))
		case 'N':
			return templates.NewsTemplate(p, renderHTML("No news this week!"))
		default:
			return templates.PostTemplate(p, renderHTML("Error: post not found"))
		}
	}
	headers := post.ParseHeaders(fileName)
	return templates.PostTemplate(p, renderHTML(headers.Contents))
}

func issuePage(issue *post.Post, post *post.Post, news *post.Post, mailbag *post.Post) templ.Component {
	return templates.IssueTemplate(
		issue,
		postBody(issue),
		postBody(post),
		postBody(news),
		postBody(mailbag),
	)
}

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
	htmx := req.URL.Query().Get("htmx")
	return htmx == "true"
}
