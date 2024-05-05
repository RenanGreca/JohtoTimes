package internal

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func ParseMarkdown(md string) (map[string]interface{}, bytes.Buffer) {
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

	return metadata, buf
}

func ExtractSlug(fileName string) string {
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, ".")
	slug := split2[0]
	return slug
}

func ExtractDate(fileName string) time.Time {
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

func ExtractTags(metadata map[string]interface{}) []string {
	if metadata["Tags"] == nil {
		return []string{}
	}
	return strings.Split(metadata["Tags"].(string), ",")
}
