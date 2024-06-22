package comment

import "time"

type Comment struct {
	ID      int64
	PostID  int64
	Name    string
	Email   string
	Date    time.Time
	Content string
	// Replies    []*Comment
	// Parent     *Comment
	IsDeleted  bool
	IsSpam     bool
	IsApproved bool
}
