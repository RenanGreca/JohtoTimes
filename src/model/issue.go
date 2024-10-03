package model

type Issue struct {
	ID          int64
	Title       string
	Slug        string
	FileName    string
	Category    Category
	Tags        []Category
	Img         string
	Description string
	Issue       int
	Volume      int
}
