package issue

import (
	"time"

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

	// These are stored in the DB as FKs to the respective tables
	Post    *post.Post
	Mailbag *post.Post
	News    *post.Post
}
