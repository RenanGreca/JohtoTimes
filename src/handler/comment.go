package handler

import (
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/afocus/captcha"
	"github.com/google/uuid"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/comment"
	"johtotimes.com/src/database"
	"johtotimes.com/src/internal"
	"johtotimes.com/src/templates"
)

func CommentHandler(w http.ResponseWriter, req *http.Request) {
	postID, err := strconv.ParseInt(req.PathValue("postID"), 10, 64)
	assert.NoError(err, "CommentHandler: Error parsing postID")

	switch req.Method {
	case "POST":
		createPostComment(req, postID)
		commentForm(postID).Render(req.Context(), w)
		getPostComments(postID).Render(req.Context(), w)
		return
	case "GET":
		commentForm(postID).Render(req.Context(), w)
		getPostComments(postID).Render(req.Context(), w)
		return
	default:
		errorPage(405).Render(req.Context(), w)
		return
	}
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

	// TODO: Check if captcha still exists in DB
	captchaID := req.FormValue("captchaID")
	captcha := db.Captchas.Retrieve(captchaID)
	captchaInput := req.FormValue("captcha")
	// TODO: Show a message if captcha is invalid
	assert.Equal(captcha.Value, captchaInput, "CommentHandler: Captcha value does not match")

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

	db.Captchas.Delete(captcha.ID)
}

func commentForm(postID int64) templ.Component {
	captchaID := uuid.New()
	return templates.CreateCommentTemplate(postID, captchaID.String())
}

func CaptchaHandler(w http.ResponseWriter, req *http.Request) {
	captchaID := req.PathValue("captchaID")
	cap := captcha.New()
	cap.SetSize(256, 64)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	cap.SetFont(internal.AssetPath + "/fonts/unown.ttf")
	img, str := cap.Create(6, captcha.UPPER)
	captcha := comment.Captcha{
		UUID:  captchaID,
		Value: str,
	}

	db := database.Connect()
	defer db.Close()
	db.Captchas.Create(&captcha)

	png.Encode(w, img)
}
