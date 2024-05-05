package news

import (
	"database/sql"
	"log"

	"johtotimes.com/src/internal"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

func (r *NewsRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATETIME NOT NULL
	)
	`

	_, err := r.db.Exec(query)
	if err != nil {
		log.Println("Error running query")
		log.Println(query)
	}
	return err
}

func (r *NewsRepository) Populate(db *sql.DB) {
	if err := r.Migrate(); err != nil {
		log.Fatal(err)
	}

	news := getFromDirectory(internal.NewsPath)
	for _, n := range news {
		created, err := r.Create(n)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created news with ID %d\n", created.ID)
	}
}

func (r *NewsRepository) Create(news News) (*News, error) {
	query := `
	INSERT INTO news(date)
	values(?)
	`

	res, err := r.db.Exec(query, news.Date)
	if err != nil {
		log.Println(query)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	news.ID = id

	return &news, nil
}
