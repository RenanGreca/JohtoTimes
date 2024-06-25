package file

import (
	"log"
	"os"
)

// ReadFile reads a file and returns its contents as a string.
func ReadFile(fileName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

// FileExists checks if a file exists in the filesystem.
func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}
	return true
}
