package post

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"johtotimes.com/src/internal"
)

// Received the path to a markdown file and returns a Post element
func parseHeaders(fileName string) (Post, string) {
	md := internal.ReadFile(fileName)

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(md), &buf, parser.WithContext(context)); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	metadata := meta.Get(context)

	return Post{
		// Contents: content,
		Slug:  extractSlug(fileName),
		Title: metadata["Title"].(string),
		// Category:    metadata["Category"].(string),
		Img:         metadata["Header"].(string),
		Description: metadata["Description"].(string),
		Date:        extractDate(fileName),
		// Tags:        extractTags(metadata),
	}, buf.String()
}

func extractSlug(fileName string) string {
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, ".")
	slug := split2[0]
	return slug
}

func extractDate(fileName string) time.Time {
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, "-")
	year, err := strconv.Atoi(split2[0])
	if err != nil {
		log.Fatal(err)
	}
	month, err := strconv.Atoi(split2[1])
	if err != nil {
		log.Fatal(err)
	}
	day, err := strconv.Atoi(split2[2])
	if err != nil {
		log.Fatal(err)
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func extractTags(metadata map[string]interface{}) []string {
	if metadata["Tags"] == nil {
		return []string{}
	}
	return strings.Split(metadata["Tags"].(string), ",")
}
