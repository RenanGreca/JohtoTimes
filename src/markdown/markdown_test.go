package markdown

import (
	"log"
	"testing"

	"johtotimes.com/src/file"
)

const fileName = "../../web/posts/2024-02-24-pokemon-legends-celebi-a-concept.md"

type TestCase struct {
	FileName string
	Title    string
	Category string
	Tags     []string
	Slug     string
	Post     Markdown
}

var testCases = []TestCase{
	{
		FileName: "../../web/posts/2024-02-24-pokemon-legends-celebi-a-concept.md",
		Title:    "Pokémon Legends Celebi: A Concept",
		Category: "Feature",
		Tags:     []string{"celebi", "johto"},
		Slug:     "2024-02-24-pokemon-legends-celebi-a-concept",
	},
	{
		FileName: "../../web/posts/2024-02-29-interview-with-the-pokemasters.md",
		Title:    "Interview with The Pokémasters",
		Category: "Interview",
		Tags:     []string{"interview", "pokemasters"},
		Slug:     "2024-02-29-interview-with-the-pokemasters",
	},
	{
		FileName: "../../web/posts/2023-11-30-interview-with-johto.md",
		Title:    "Interview with Johto",
		Category: "Interview",
		Tags:     []string{"interview", "johto"},
		Slug:     "2023-11-30-interview-with-johto",
	},
}

func TestReadFile(t *testing.T) {
	for _, tc := range testCases {
		_ = file.ReadFile(tc.FileName)
	}
}

func TestBeforeAll(t *testing.T) {
	for i, tc := range testCases {
		fileName := tc.FileName
		post := ParseHeaders(fileName)
		log.Printf("======================")
		log.Printf("TEST %d", i)
		log.Printf("Parsed %s", fileName)
		log.Printf("Title: %s", post.Metadata.Title)
		log.Printf("Category: %s", post.Metadata.Category)
		log.Printf("Tags: %v", post.Metadata.Tags)
		log.Printf("Slug: %s", post.Slug)
		log.Printf("======================")
	}
}

func TestTitle(t *testing.T) {
	for _, tc := range testCases {
		post := ParseHeaders(tc.FileName)
		expected := tc.Title
		if post.Metadata.Title != expected {
			t.Fatalf(`Expected title %q, received %q`, expected, post.Metadata.Title)
		}
	}
}

func TestCategory(t *testing.T) {
	for _, tc := range testCases {
		post := ParseHeaders(tc.FileName)
		expected := tc.Category
		if post.Metadata.Category != expected {
			t.Fatalf(`Expected category %q, received %q`, expected, post.Metadata.Category)
		}
	}
}

func TestTags(t *testing.T) {
	for _, tc := range testCases {
		post := ParseHeaders(tc.FileName)
		expected := tc.Tags
		if len(post.Metadata.Tags) != len(expected) {
			t.Fatalf(`Expected tags count %d, received %d`, len(expected), len(post.Metadata.Tags))
		}
	}
}

func TestSlug(t *testing.T) {
	for _, tc := range testCases {
		slug := extractSlug(tc.FileName)
		expected := tc.Slug
		if slug != expected {
			t.Fatalf(`Expected slug %q, received %q`, expected, slug)
		}
	}
}
