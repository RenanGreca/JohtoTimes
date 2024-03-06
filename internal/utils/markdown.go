package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"johtotimes.com/internal/types"
)

func ParseMarkdown(fileName string) types.Post {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	md := readFile(fileName)
	split := strings.Split(fileName, "/")
	last := split[len(split)-1]
	split2 := strings.Split(last, ".")
	slug := split2[0]
	log.Println("Post slug: " + slug)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert([]byte(md), &buf, parser.WithContext(context)); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	content := unsafe(buf.String())
	metadata := meta.Get(context)
	title := metadata["Title"]
	fmt.Println("Title: ", title)

	return types.Post{
		Contents: content,
		Metadata: metadata,
		Slug:     slug,
	}
}

func unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func readFile(fileName string) string {
	log.Println("Opening file: " + fileName)
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
