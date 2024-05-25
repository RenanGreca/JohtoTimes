package database

import (
	"johtotimes.com/src/post"
)

type Repository interface {
	post.PostRepository
}
