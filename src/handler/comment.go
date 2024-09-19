package handler

import (
	"image/png"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-tts/tts/pkg/speech"
	"github.com/google/uuid"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/constants"
	"johtotimes.com/src/database"
	"johtotimes.com/src/file"
	"johtotimes.com/src/model"
	"johtotimes.com/src/templates"
)

func CommentHandler(w http.ResponseWriter, req *http.Request) {
	postID, err := strconv.ParseInt(req.PathValue("postID"), 10, 64)
	assert.NoError(err, "CommentHandler: Error parsing postID")

	switch req.Method {
	case "POST":
		comment, errMsg := createPostComment(req, postID)
		if len(errMsg) == 0 {
			templates.CreateCommentButton(postID).Render(req.Context(), w)
		} else {
			commentForm(comment, errMsg...).Render(req.Context(), w)
		}
		getPostComments(postID).Render(req.Context(), w)
		return
	case "GET":
		getPostComments(postID).Render(req.Context(), w)
		return
	default:
		errorPage(405).Render(req.Context(), w)
		return
	}
}

func NewCommentHandler(w http.ResponseWriter, req *http.Request) {
	postID, err := strconv.ParseInt(req.PathValue("postID"), 10, 64)
	assert.NoError(err, "CommentHandler: Error parsing postID")

	commentForm(model.Comment{PostID: postID}).Render(req.Context(), w)
}

func getPostComments(postID int64) templ.Component {
	assert.NotNil(postID, "CommentHandler: postID cannot be nil")
	assert.LogDebug("Retrieving comments for post %d", postID)
	db := database.Connect()
	defer db.Close()
	comments := db.Comments.GetCommentsFromPost(postID)
	assert.LogDebug("Found %d comments", len(comments))
	return templates.CommentListTemplate(comments)
}

func createPostComment(req *http.Request, postID int64) (model.Comment, []string) {
	assert.NotNil(postID, "CommentHandler: postID cannot be nil")
	assert.LogDebug("Creating comment for post %d", postID)
	db := database.Connect()
	defer db.Close()

	var errMsg []string

	captchaID := req.FormValue("captchaID")
	captcha, err := db.Captchas.Retrieve(captchaID)
	if err != nil {
		errMsg = append(errMsg, "Invalid captcha")
	} else {
		captchaInput := strings.ToUpper(req.FormValue("captcha"))
		if captchaInput == "" || captchaInput != captcha.Value {
			errMsg = append(errMsg, "Invalid captcha")
		}
	}

	name := req.FormValue("name")
	if name == "" {
		errMsg = append(errMsg, "Name cannot be empty")
	}
	content := req.FormValue("content")
	if content == "" {
		errMsg = append(errMsg, "Content cannot be empty")
	}
	date := time.Now()

	comment := model.Comment{
		PostID:     postID,
		Name:       name,
		Email:      "",
		Content:    content,
		Date:       date,
		IsDeleted:  false,
		IsSpam:     false,
		IsApproved: true,
	}

	if len(errMsg) == 0 {
		db.Comments.Create(&comment)
		assert.LogDebug("Created comment with ID %d", comment.ID)
	}

	db.Captchas.Delete(captcha.UUID)
	file.Delete(constants.AudioPath + "/" + captchaID + ".mp3")
	return comment, errMsg
}

func commentForm(comment model.Comment, errMsg ...string) templ.Component {
	captchaID := uuid.New().String()
	return templates.CreateCommentTemplate(comment, captchaID, errMsg...)
}

func CaptchaHandler(w http.ResponseWriter, req *http.Request) {
	captchaID := req.PathValue("captchaID")
	captcha := model.NewCaptcha(captchaID)

	db := database.Connect()
	defer db.Close()
	assert.LogDebug("Creating captcha with ID %s\n", captchaID)
	db.Captchas.Create(&captcha)

	png.Encode(w, captcha.Image)
}

func NewCaptchaHandler(w http.ResponseWriter, req *http.Request) {
	captchaID := req.PathValue("captchaID")
	db := database.Connect()
	defer db.Close()

	db.Captchas.Delete(captchaID)
	file.Delete(constants.AudioPath + "/" + captchaID + ".mp3")

	captchaID = uuid.New().String()
	templates.CaptchaTemplate(captchaID).Render(req.Context(), w)
}

func AudioCaptchaHandler(w http.ResponseWriter, req *http.Request) {
	captchaID := req.PathValue("captchaID")
	db := database.Connect()
	defer db.Close()

	assert.LogDebug("Retrieving captcha with ID %s", captchaID)
	captcha, err := db.Captchas.Retrieve(captchaID)
	if err != nil {
		return
	}

	file := file.Create(constants.AudioPath + "/" + captchaID + ".mp3")
	assert.NoError(err, "AudioCaptchaHandler: Error creating file")
	defer file.Close()

	reader := strings.NewReader(captcha.Value)
	err = speech.WriteToAudioStream(reader, file, speech.LangEn)
	assert.NoError(err, "AudioCaptchaHandler: Error writing audio to file")

	templates.AudioCaptchaTemplate(captchaID, true).Render(req.Context(), w)
}
