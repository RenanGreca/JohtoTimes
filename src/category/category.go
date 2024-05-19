package category

// Categories and Tags are defined as the same entity
// The difference is that posts and categories are 1-to-n
// posts and tags are n-to-n
type Category struct {
	ID   int64
	Name string
	Slug string
	Type byte // 'C' for category, 'T' for tag
}
