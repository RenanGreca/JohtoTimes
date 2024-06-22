package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/category"
	"johtotimes.com/src/comment"
	"johtotimes.com/src/post"
)

type Database struct {
	Connection *sql.DB
	Posts      *post.PostRepository
	Categories *category.CategoryRepository
	Comments   *comment.CommentRepository
	Captchas   *comment.CaptchaRepository
}

// type Issue struct {
// 	Volume int
// 	Issue  int
// 	Title  string
// 	News   templ.Component
// 	// Post        *Post
// 	Mailbag     templ.Component
// 	Description string
// }
//
// type Category struct {
// 	Name   string
// 	Plural string
// 	Slug   string
// 	// Posts  []*Post
// }

const dbFile = "sqlite.db"

// var DB = NewDB()

// For now this function always creates a DB from scratch
func NewDB() {
	log.Println("Creating new database")
	os.Remove(dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	defer db.Close()
	assert.NoError(err, "Database: Error opening database")

	categoryRepository := category.NewCategoryRepository(db)
	categoryRepository.Migrate()

	postRepository := post.NewPostRepository(db)
	postRepository.Populate(db)

	commentRepository := comment.NewCommentRepository(db)
	commentRepository.Migrate()

	captchaRepository := comment.NewCaptchaRepository(db)
	captchaRepository.Migrate()

	// database := Database{
	// 	Connection: db,
	// 	Posts:      postRepository,
	// }
	// return &database
}

func Connect() *Database {
	log.Println("Opening connection to existing database")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	postRepository := post.NewPostRepository(db)
	categoryRepository := category.NewCategoryRepository(db)
	commentRepository := comment.NewCommentRepository(db)
	captchaRepository := comment.NewCaptchaRepository(db)

	database := Database{
		Connection: db,
		Posts:      postRepository,
		Categories: categoryRepository,
		Comments:   commentRepository,
		Captchas:   captchaRepository,
	}
	return &database
}

func (db *Database) Close() {
	log.Println("Closing connection")
	db.Connection.Close()
}
