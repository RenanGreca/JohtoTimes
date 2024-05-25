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
	FileName    string
	Category    *category.Category
	Tags        []*category.Category
	Img         string
	Description string
	Type        byte // 'P' for post, 'N' for news, 'M' for mailbag
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
		post := ParseHeaders(fileName)
		// posts[i] = post
		posts = append(posts, post)
	}
	return posts
}

// Received the path to a markdown file and returns a Post element
func ParseHeaders(fileName string) Markdown {
	md := internal.ReadFile(fileName)

	metadata, buf := ParseMarkdown(md)

	return Markdown{
		FileName: fileName,
		Slug:     ExtractSlug(fileName),
		Date:     ExtractDate(fileName),
		Metadata: ExtractMetadata(metadata),
		Contents: buf.String(),
	}
}
