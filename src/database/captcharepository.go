package database

import (
	"database/sql"
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL,
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

	res, err := r.db.Exec(query,
		captcha.UUID,
		captcha.Value,
		time.Now(),
	)

	id, err := res.LastInsertId()
	assert.NoError(err, "CaptchaRepository: Error getting last insert ID")
	captcha.ID = id
}

func (r *CaptchaRepository) Retrieve(uuid string) model.Captcha {
	query := `
	SELECT id, value, date
	FROM captcha
	WHERE uuid = ?`
	row := r.db.QueryRow(query, uuid)

	var captcha model.Captcha
	err := row.Scan(
		&captcha.ID,
		&captcha.Value,
		&captcha.CreatedAt,
	)
	assert.NoError(err, "CaptchaRepository: Error scanning row")
	captcha.UUID = uuid

	return captcha
}

func (r *CaptchaRepository) Delete(id int64) {
	query := `DELETE FROM captcha WHERE id = ?`
	_, err := r.db.Exec(query, id)
	assert.NoError(err, "CaptchaRepository: Error running query: %s", query)
}
