package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/file"
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

const DEV_DB_FILE = "sqlite_test.db"
const PROD_DB_FILE = "sqlite.db"

var selectedDbFile string

func NewDB(dbFile string) {
	selectedDbFile = dbFile

	if file.FileExists(dbFile) && os.Getenv("ENV") == "prod" {
		// Skip if database already exists or program is in prod mode
		return
	}

	assert.LogDebug("Creating new database")
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
	assert.LogDebug("Opening connection to existing database")
	db, err := sql.Open("sqlite3", selectedDbFile)
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
	assert.LogDebug("Closing connection")
	db.Connection.Close()
}
