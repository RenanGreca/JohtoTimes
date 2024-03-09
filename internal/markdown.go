package internal

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// Received the path to a markdown file and returns a Post element
func ParseMarkdown(fileName string) Post {
	md := ReadFile(fileName)

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
	content := unsafe(buf.String())
	metadata := meta.Get(context)

	return Post{
		Contents:    content,
		Slug:        extractSlug(fileName),
		Title:       metadata["Title"].(string),
		Category:    metadata["Category"].(string),
		Img:         metadata["Header"].(string),
		Description: metadata["Description"].(string),
		Tags:        extractTags(metadata),
	}
}

func extractSlug(fileName string) string {
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, ".")
	slug := split2[0]
	return slug
}

func extractTags(metadata map[string]interface{}) []string {
	if metadata["Tags"] == nil {
		return []string{}
	}
	return strings.Split(metadata["Tags"].(string), ",")
}

func unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
