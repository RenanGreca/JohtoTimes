package mailbag

import (
	"database/sql"
	"log"

	"johtotimes.com/src/internal"
)

type MailbagRepository struct {
	db *sql.DB
}

func NewMailbagRepository(db *sql.DB) *MailbagRepository {
	return &MailbagRepository{
		db: db,
	}
}

func (r *MailbagRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS mailbag(
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

func (r *MailbagRepository) Populate(db *sql.DB) {
	if err := r.Migrate(); err != nil {
		log.Fatal(err)
	}

	mailbags := getFromDirectory(internal.MailbagPath)
	for _, m := range mailbags {
		created, err := r.Create(m)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created mailbag with ID %d\n", created.ID)
	}
}

func (r *MailbagRepository) Create(mailbag Mailbag) (*Mailbag, error) {
	query := `
	INSERT INTO mailbag(date)
	values(?)
	`

	res, err := r.db.Exec(query, mailbag.Date)
	if err != nil {
		log.Println(query)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	mailbag.ID = id

	return &mailbag, nil
}
