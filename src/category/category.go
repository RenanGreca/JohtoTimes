package category

// Categories and Tags are defined as the same entity
// The difference is that posts and categories are 1-to-n
// posts and tags are n-to-n
type Category struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	Type        byte // 'C' for category, 'T' for tag
}

func GetFromFile(dir string, slug string) Markdown {
	fileName := dir + "/" + slug + ".md"
	md := parseHeaders(fileName)
	return md
}
