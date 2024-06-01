package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"johtotimes.com/src/database"
	"johtotimes.com/src/internal"
	"johtotimes.com/src/post"
	"johtotimes.com/src/templates"
)

func PostHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling Post request to " + req.URL.Path)
	slug := req.PathValue("slug")
	db := database.Connect()
	defer db.Close()
	post, err := db.Posts.GetBySlug(slug, 'P')
	if err != nil {
		log.Fatal(err)
	}

	postPage(post).Render(req.Context(), w)

}

func postPage(p *post.Post) templ.Component {
	filePath := internal.PostTypePath[p.Type]
	fileName := filePath + "/" + p.Slug + ".md"
	fmt.Printf("File path: %s %s\n", string(p.Type), fileName)

	headers := post.ParseHeaders(fileName)
	postBody := templates.SingleTemplate(p, renderHTML(headers.Contents))

	return templates.Base(p.Title, postBody)
}

func postBody(p *post.Post) templ.Component {
	filePath := internal.PostTypePath[p.Type]
	fileName := filePath + "/" + p.Slug + ".md"
	fmt.Printf("File path: %s %s\n", string(p.Type), fileName)

	if _, err := os.Stat(fileName); err != nil {
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
	postBody := templates.PostTemplate(p, renderHTML(headers.Contents))

	return postBody
}

func IssueHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Handling Issue request to " + req.URL.Path)
	slug := req.PathValue("slug")
	db := database.Connect()
	defer db.Close()
	issue, err := db.Posts.GetBySlug(slug, 'I')
	if err != nil {
		log.Fatal(err)
	}

	post, err := db.Posts.GetByDateAndType(issue.Date, 'P')
	if err != nil {
		log.Fatal(err)
	}

	news, err := db.Posts.GetByDateAndType(issue.Date, 'N')
	if err != nil {
		log.Fatal(err)
	}

	mailbag, err := db.Posts.GetByDateAndType(issue.Date, 'M')
	if err != nil {
		log.Fatal(err)
	}

	issuePage(issue, post, news, mailbag).Render(req.Context(), w)
}

func issuePage(i *post.Post, p *post.Post, n *post.Post, m *post.Post) templ.Component {
	issueBody := templates.IssueTemplate(i, postBody(i), postBody(p), postBody(n), postBody(m))

	return templates.Base(p.Title, issueBody)
}
