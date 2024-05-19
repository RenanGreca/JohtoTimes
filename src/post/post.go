package post

import (
	"log"
	"os"
	"time"

	"johtotimes.com/src/category"
	"johtotimes.com/src/internal"
)

type Post struct {
	ID          int64
	Title       string
	String      string
	Slug        string
	Category    *category.Category
	Tags        []*category.Category
	Img         string
	Description string
	Date        time.Time
}

type Markdown struct {
	FileName string
	Slug     string
	Date     time.Time
	Metadata Metadata
	Contents string
}

type Metadata struct {
	Title       string
	Header      string
	Category    string
	Description string
	Tags        []string
}

func getFromDirectory(postsDir string) []Markdown {
	entries, err := os.ReadDir(postsDir)
	if err != nil {
		log.Fatalln(err)
	}

	var posts []Markdown
	for _, e := range entries {
		fileName := postsDir + "/" + e.Name()
		post := parseHeaders(fileName)
		// posts[i] = post
		posts = append(posts, post)
	}
	return posts
}

// Received the path to a markdown file and returns a Post element
func parseHeaders(fileName string) Markdown {
	md := internal.ReadFile(fileName)

	metadata, buf := internal.ParseMarkdown(md)

	// db := database.Connect()
	// defer db.Close()
	// catSlug := slug.Make(metadata["Category"].(string))
	// cat, err := db.Categories.Create(catSlug, 'C')
	// if err != nil {
	// 	log.Println(err)
	// }
	return Markdown{
		FileName: fileName,
		Slug:     internal.ExtractSlug(fileName),
		Date:     internal.ExtractDate(fileName),
		Metadata: Metadata{
			Title:  metadata["Title"].(string),
			Header: metadata["Header"].(string),
			// Tags:        metadata["Tags"].([]string),
			Category:    metadata["Category"].(string),
			Description: metadata["Description"].(string),
		},
		Contents: buf.String(),
	}

	// return Post{
	// 	// Contents: content,
	// 	Slug:        internal.ExtractSlug(fileName),
	// 	Title:       metadata["Title"].(string),
	// 	Category:    cat,
	// 	Img:         metadata["Header"].(string),
	// 	Description: metadata["Description"].(string),
	// 	Date:        internal.ExtractDate(fileName),
	// 	// Tags:        extractTags(metadata),
	// }, buf.String()
}
