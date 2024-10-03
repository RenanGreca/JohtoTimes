package model

import (
	"crypto/rand"
	"log"
	"os"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID       int64
	Name     string
	Email    string
	Password []byte
	Salt     []byte
}

func CreateDefaultUser() User {
	name := os.Getenv("JOHTOTIMES_ADMIN_NAME")
	email := os.Getenv("JOHTOTIMES_ADMIN_EMAIL")
	password := []byte(os.Getenv("JOHTOTIMES_ADMIN_PASSWORD"))
	salt := generateSalt()

	log.Printf("Creating default user with name %s, email %s, password %s", name, email, password)

	passwordhash := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)

	return User{
		Name:     name,
		Email:    email,
		Password: passwordhash,
		Salt:     salt,
	}
}

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}
