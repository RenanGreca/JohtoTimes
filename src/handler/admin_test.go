package handler

import (
	"crypto/rand"
	"log"
	"testing"

	"golang.org/x/crypto/argon2"
)

func TestPassword(t *testing.T) {
	password := []byte("password")
	salt := []byte("salt")
	key := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	log.Printf("Key: %x", key)
	log.Printf("Length: %d", len(string(key)))
}

func TestGenerateSalt(t *testing.T) {
	salt := make([]byte, 16)
	rand.Read(salt)
	log.Printf("Salt: %x", salt)
}
