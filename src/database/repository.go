package database

import (
	"johtotimes.com/src/mailbag"
	"johtotimes.com/src/post"
)

type Repository interface {
	post.PostRepository | mailbag.MailbagRepository
}
