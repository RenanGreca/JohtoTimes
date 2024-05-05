package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"johtotimes.com/src/post"
)

type Database struct {
	Connection *sql.DB
	Posts      *post.PostRepository
	// Issues     []*Issue
	// Categories []Category
	// Tags       []Category
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

	if err != nil {
		log.Fatal(err)
	}
	postRepository := post.NewPostRepository(db)
	postRepository.Populate(db)

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

	database := Database{
		Connection: db,
		Posts:      postRepository,
	}
	return &database
}

func (db *Database) Close() {
	log.Println("Closing connection")
	db.Connection.Close()
}
