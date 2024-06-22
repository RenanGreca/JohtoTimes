package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/comment"
	"johtotimes.com/src/database"
	"johtotimes.com/src/templates"
)

func CommentHandler(w http.ResponseWriter, req *http.Request) {
	postID, err := strconv.ParseInt(req.PathValue("postID"), 10, 64)
	assert.NoError(err, "CommentHandler: Error parsing postID")

	var component templ.Component
	switch req.Method {
	case "POST":
		createPostComment(req, postID)
		component = getPostComments(postID)
	case "GET":
		component = getPostComments(postID)
	default:
		component = errorPage(405)
	}

	component.Render(req.Context(), w)
}

func getPostComments(postID int64) templ.Component {
	assert.NotNil(postID, "CommentHandler: postID cannot be nil")
	log.Printf("Retrieving comments for post %d", postID)
	db := database.Connect()
	defer db.Close()
	comments := db.Comments.GetCommentsFromPost(postID)
	log.Printf("Found %d comments", len(comments))
	return templates.CommentListTemplate(comments)
}

func createPostComment(req *http.Request, postID int64) {
	assert.NotNil(postID, "CommentHandler: postID cannot be nil")
	log.Printf("Creating comment for post %d", postID)
	db := database.Connect()
	defer db.Close()

	name := req.FormValue("name")
	email := req.FormValue("email")
	content := req.FormValue("content")
	date := time.Now()

	comment := comment.Comment{
		PostID:     postID,
		Name:       name,
		Email:      email,
		Content:    content,
		Date:       date,
		IsDeleted:  false,
		IsSpam:     false,
		IsApproved: true,
	}
	db.Comments.Create(&comment)
	log.Printf("Created comment with ID %d", comment.ID)
}
