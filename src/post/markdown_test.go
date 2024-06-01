package post

import (
	"testing"

	"johtotimes.com/src/internal"
)

const fileName = "../../web/posts/2024-02-24-pokemon-legends-celebi-a-concept.md"

var post Markdown

func TestReadFile(t *testing.T) {
	_ = internal.ReadFile(fileName)
}

func TestBeforeAll(t *testing.T) {
	post = ParseHeaders(fileName)
}

func TestTitle(t *testing.T) {
	expected := "Pokémon Legends Celebi: A Concept"
	if post.Metadata.Title != expected {
		t.Fatalf(`Expected title %q, received %q`, expected, post.Metadata.Title)
	}
}

// func TestCategory(t *testing.T) {
// 	expected := "Features"
// 	if post.Category != expected {
// 		t.Fatalf(`Expected category %q, received %q`, expected, post.Category)
// 	}
// }
//
// func TestTags(t *testing.T) {
// 	expectedCount := 0
// 	if len(post.Tags) != expectedCount {
// 		t.Fatalf(`Expected tags count %d, received %d`, expectedCount, len(post.Tags))
// 	}
// }

func TestSlug(t *testing.T) {
	slug := extractSlug(fileName)
	expected := "2024-02-24-pokemon-legends-celebi-a-concept"
	if slug != expected {
		t.Fatalf(`Expected slug %q, received %q`, expected, slug)
	}
}
