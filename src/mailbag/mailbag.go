package mailbag

import (
	"log"
	"os"
	"time"

	"johtotimes.com/src/internal"
)

type Mailbag struct {
	ID   int64
	Date time.Time
}

func getFromDirectory(mailbagDir string) []Mailbag {
	entries, err := os.ReadDir(mailbagDir)
	if err != nil {
		log.Fatalln(err)
	}

	var mailbags []Mailbag
	for _, e := range entries {
		mailbag, _ := parseHeaders(mailbagDir + "/" + e.Name())
		mailbags = append(mailbags, mailbag)
	}

	return mailbags
}

func parseHeaders(fileName string) (Mailbag, string) {
	md := internal.ReadFile(fileName)

	_, buf := internal.ParseMarkdown(md)

	return Mailbag{
		Date: internal.ExtractDate(fileName),
	}, buf.String()
}
