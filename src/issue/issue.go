package issue

import (
	"time"

	"johtotimes.com/src/mailbag"
	"johtotimes.com/src/news"
	"johtotimes.com/src/post"
)

type Issue struct {
	ID          int64
	Title       string
	Description string
	Slug        string
	Volume      int
	Number      int
	Date        time.Time

	Post    *post.Post
	Mailbag *mailbag.Mailbag
	News    *news.News
}
