package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/constants"
	"johtotimes.com/src/database"
	"johtotimes.com/src/file"
	"johtotimes.com/src/model"
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
	post, err := db.Posts.GetByDateAndType(issue.CreatedAt, 'P')
	if err != nil {
		log.Printf("IssueHandler: Error getting post: %s", err)
		errorPage(404).Render(req.Context(), w)
		return
	}
	news, err := db.Posts.GetByDateAndType(issue.CreatedAt, 'N')
	assert.NoError(err, "IssueHandler: Error getting news")
	mailbag, err := db.Posts.GetByDateAndType(issue.CreatedAt, 'M')
	assert.NoError(err, "IssueHandler: Error getting mailbag")

	render(issuePage(issue, post, news, mailbag), isHTMX(req), issue.Title, w)
}

func postPage(p *model.Post) templ.Component {
	filePath := constants.PostTypePath[p.Type]
	fileName := filePath + "/" + p.Slug + ".md"
	fmt.Printf("File path: %s %s\n", string(p.Type), fileName)

	return templates.SingleTemplate(p, renderHTML(p.Content()))
}

func postBody(p *model.Post) templ.Component {
	if p == nil {
		return errorPage(404)
	}
	filePath := constants.PostTypePath[p.Type]
	fileName := filePath + "/" + p.Slug + ".md"
	fmt.Printf("File path: %s %s\n", string(p.Type), fileName)

	// In the issue builder, if the file doesn't exist, we show a message
	if !file.FileExists(fileName) {
		switch p.Type {
		case 'M':
			return templates.MailbagTemplate(p, renderHTML("No mailbag this week!"))
		case 'N':
			return templates.NewsTemplate(p, renderHTML("No news this week!"))
		default:
			return templates.PostTemplate(p, renderHTML("Error: post not found"))
		}
	}
	return templates.PostTemplate(p, renderHTML(p.Content()))
}

func issuePage(issue *model.Post, post *model.Post, news *model.Post, mailbag *model.Post) templ.Component {
	return templates.IssueTemplate(
		issue,
		postBody(issue),
		postBody(post),
		postBody(news),
		postBody(mailbag),
		post.ID,
	)
}
