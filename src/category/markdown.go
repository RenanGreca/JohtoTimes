package category

import (
	"bytes"
	"log"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"johtotimes.com/src/internal"
)

type Markdown struct {
	Slug     string
	Metadata Metadata
	Contents string
}

type Metadata struct {
	Singular string
	Plural   string
}

// Received the path to a markdown file and returns a Post element
func parseHeaders(fileName string) Markdown {
	md := internal.ReadFile(fileName)

	metadata, buf := parseMarkdown(md)

	return Markdown{
		Slug:     extractSlug(fileName),
		Metadata: extractMetadata(metadata),
		Contents: buf.String(),
	}
}

func parseMarkdown(md string) (map[string]interface{}, bytes.Buffer) {
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
	log.Println(metadata)

	return metadata, buf
}

func extractSlug(fileName string) string {
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, ".")
	slug := split2[0]
	return slug
}

func extractMetadata(metadata map[string]interface{}) Metadata {
	var result Metadata
	if metadata["singular"] != nil {
		result.Singular = metadata["singular"].(string)
	}
	if metadata["plural"] != nil {
		result.Plural = metadata["plural"].(string)
	}
	return result
}
