package model

import (
	"os"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/file"
	"johtotimes.com/src/markdown"
)

// Categories and Tags are defined as the same entity
// The difference is that posts and categories are 1-to-n
// posts and tags are n-to-n
type Category struct {
	ID          int64
	Singular    string
	Plural      string
	Slug        string
	Description string
	Type        byte // 'C' for category, 'T' for tag
}

func GetCategoryFromFile(dir string, slug string) Category {
	fileName := dir + "/" + slug + ".md"
	md := NewCategoryFromMarkdown(fileName)
	return md
}

func GetCategoriesFromDirectory(dir string) []Category {
	entries, err := os.ReadDir(dir)
	assert.NoError(err, "CategoryRepository: Error reading directory: %s", dir)

	var categories []Category
	for _, e := range entries {
		fileName := dir + "/" + e.Name()
		category := NewCategoryFromMarkdown(fileName)
		categories = append(categories, category)
	}
	return categories
}

// Received the path to a markdown file and returns a Post element
func NewCategoryFromMarkdown(fileName string) Category {
	md := file.ReadFile(fileName)

	metadata, buf := markdown.ParseMarkdown(md)

	category := Category{
		Slug:        markdown.ExtractSlug(fileName),
		Description: buf.String(),
		Type:        'C',
	}
	category.extractMetadata(metadata)

	return category
}

func (category *Category) extractMetadata(metadata map[string]interface{}) {
	if metadata["singular"] != nil {
		category.Singular = metadata["singular"].(string)
	}
	if metadata["plural"] != nil {
		category.Plural = metadata["plural"].(string)
	}
}

func extractTags(tags []interface{}) []Category {
	var result []Category
	for _, t := range tags {
		result = append(result, Category{
			Slug: t.(string),
		})
	}
	return result
}
