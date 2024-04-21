package internal

import (
	"log"
	"os"
)

func ReadFile(fileName string) string {
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
