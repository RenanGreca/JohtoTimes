package internal

import (
	"testing"
)

const fileName = "../../web/posts/2024-02-22-pokemon-legends-celebi-a-concept.md"

var post Post

func TestReadFile(t *testing.T) {
	_ = ReadFile(fileName)
}

func TestBeforeAll(t *testing.T) {
	post = ParseMarkdown(fileName)
}

func TestTitle(t *testing.T) {
	expected := "Pokémon Legends Celebi: A Concept"
	if post.Title != expected {
		t.Fatalf(`Expected title %q, received %q`, expected, post.Title)
	}
}

func TestCategory(t *testing.T) {
	expected := "Features"
	if post.Category != expected {
		t.Fatalf(`Expected category %q, received %q`, expected, post.Category)
	}
}

func TestTags(t *testing.T) {
	expectedCount := 0
	if len(post.Tags) != expectedCount {
		t.Fatalf(`Expected tags count %d, received %d`, expectedCount, len(post.Tags))
	}
}

func TestSlug(t *testing.T) {
	slug := extractSlug(fileName)
	expected := "2024-02-22-pokemon-legends-celebi-a-concept"
	if slug != expected {
		t.Fatalf(`Expected slug %q, received %q`, expected, slug)
	}
}
