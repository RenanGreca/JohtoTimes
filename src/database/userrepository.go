package database

import (
	"database/sql"
	"log"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password character(32) NOT NULL,
		salt character(16) NOT NULL
	);
	`

	_, err := r.db.Exec(query)
	assert.NoError(err, "UserRepository: Error running query: %s", query)
}

func (r *UserRepository) Populate() {
	user := model.CreateDefaultUser()

	created := r.Create(user)
	assert.LogDebug(
		"Created user with ID %d\n",
		created.ID,
	)

	u, err := r.GetByEmail(created.Email)
	assert.NoError(err, "UserRepository: Error getting user by email")
	log.Printf("UserRepository: Populated user with ID %d\n", u.ID)
	log.Printf("UserRepository: Populated user with name %s\n", u.Name)
	log.Printf("UserRepository: Populated user with email %s\n", u.Email)
	log.Printf("UserRepository: Populated user with password %s\n", string(u.Password))
	log.Printf("UserRepository: Populated user with salt %s\n", string(u.Salt))
}

func (r *UserRepository) Create(user model.User) *model.User {
	// Check if user already exists, return it if so.
	u, err := r.GetByEmail(user.Email)
	if err == nil && u != nil {
		return u
	}

	query := `
	INSERT INTO user(
		name, email, password, salt
	)
	values(?,?,?,?)
	`

	res, err := r.db.Exec(query,
		user.Name,
		user.Email,
		user.Password,
		user.Salt,
	)
	assert.NoError(err, "UserRepository: Error running query: %s", query)

	id, err := res.LastInsertId()
	assert.NoError(err, "UserRepository: Error getting last insert ID")
	user.ID = id

	return &user
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
	SELECT id, name, email, password, salt
	FROM user
	WHERE email = ?`

	row := r.db.QueryRow(query, email)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Salt)

	return &user, err
}
