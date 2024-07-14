package database

import (
	"database/sql"
	"log"
	"time"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/model"
)

type CaptchaRepository struct {
	db *sql.DB
}

func NewCaptchaRepository(db *sql.DB) *CaptchaRepository {
	return &CaptchaRepository{
		db: db,
	}
}

func (r *CaptchaRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS captcha(
		uuid TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		date DATETIME NOT NULL
	);`

	_, err := r.db.Exec(query)
	assert.NoError(err, "CaptchaRepository: Error running query: %s", query)
}

func (r *CaptchaRepository) Create(captcha *model.Captcha) {
	query := `
	INSERT INTO captcha(uuid, value, date)
	values(?,?,?)`

	_, err := r.db.Exec(query,
		captcha.UUID,
		captcha.Value,
		time.Now(),
	)
	assert.NoError(err, "CaptchaRepository: Error running query: %s", query)

	log.Printf("Created captcha with ID %s\n", captcha.UUID)
}

func (r *CaptchaRepository) Retrieve(uuid string) (model.Captcha, error) {
	query := `
	SELECT uuid, value, date
	FROM captcha
	WHERE uuid = ?`
	row := r.db.QueryRow(query, uuid)

	var captcha model.Captcha
	err := row.Scan(
		&captcha.UUID,
		&captcha.Value,
		&captcha.CreatedAt,
	)

	log.Printf("Retrieved captcha with ID %s\n", captcha.UUID)
	return captcha, err
}

func (r *CaptchaRepository) Delete(uuid string) {
	query := `DELETE FROM captcha WHERE uuid = ?`
	_, err := r.db.Exec(query, uuid)
	assert.NoError(err, "CaptchaRepository: Error running query: %s", query)
}

func (r *CaptchaRepository) DeleteOld() {
	query := `DELETE FROM captcha WHERE date < ?`
	_, err := r.db.Exec(query, time.Now().AddDate(0, 0, -1))
	assert.NoError(err, "CaptchaRepository: Error running query: %s", query)
	// TODO: Also delete related files
}
