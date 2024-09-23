package database

import (
	"crypto/subtle"
	"database/sql"
	"os"
	"testing"

	"golang.org/x/crypto/argon2"
	"johtotimes.com/src/assert"
)

var userRepository *UserRepository

func TestBeforeAll(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(err, "Error opening database")
	userRepository = NewUserRepository(db)
	userRepository.Migrate()
	userRepository.Populate()
}

func TestPasswordMatch(t *testing.T) {
	u, err := userRepository.GetByEmail(os.Getenv("JOHTOTIMES_ADMIN_EMAIL"))
	assert.NoError(err, "UserRepository: Error getting user by email")

	password := []byte(os.Getenv("JOHTOTIMES_ADMIN_PASSWORD"))
	passwordhash := argon2.IDKey(password, u.Salt, 1, 64*1024, 4, 32)
	passwordmatch := subtle.ConstantTimeCompare(passwordhash, u.Password[:])
	if passwordmatch != 1 {
		t.Fatalf("Passwords do not match")
	}
}
