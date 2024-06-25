package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"johtotimes.com/src/assert"
)

type Database struct {
	Connection *sql.DB
	Posts      *PostRepository
	Categories *CategoryRepository
	Comments   *CommentRepository
	Captchas   *CaptchaRepository
}

type Repository interface {
	Migrate()
	Populate(db *sql.DB)
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

	categoryRepository := NewCategoryRepository(db)
	categoryRepository.Migrate()
	categoryRepository.Populate()

	postRepository := NewPostRepository(db)
	postRepository.Migrate()
	postRepository.Populate()

	commentRepository := NewCommentRepository(db)
	commentRepository.Migrate()

	captchaRepository := NewCaptchaRepository(db)
	captchaRepository.Migrate()
}

func Connect() *Database {
	log.Println("Opening connection to existing database")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	postRepository := NewPostRepository(db)
	categoryRepository := NewCategoryRepository(db)
	commentRepository := NewCommentRepository(db)
	captchaRepository := NewCaptchaRepository(db)

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

func printQuery(query string, args ...interface{}) {
	query = strings.ReplaceAll(query, "?", "%q")
	log.Printf(query, args...)
}
