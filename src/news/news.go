package news

import (
	"log"
	"os"
	"time"

	"johtotimes.com/src/internal"
)

type News struct {
	ID   int64
	Date time.Time
}

func getFromDirectory(newsDir string) []News {
	entries, err := os.ReadDir(newsDir)
	if err != nil {
		log.Fatalln(err)
	}

	var news []News
	for _, e := range entries {
		n, _ := parseHeaders(newsDir + "/" + e.Name())
		news = append(news, n)
	}

	return news
}

func parseHeaders(fileName string) (News, string) {
	md := internal.ReadFile(fileName)

	_, buf := internal.ParseMarkdown(md)

	return News{
		Date: internal.ExtractDate(fileName),
	}, buf.String()
}
