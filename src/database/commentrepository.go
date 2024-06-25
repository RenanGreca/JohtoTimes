package database

import (
	"database/sql"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/model"
)

const selectComments = `
	SELECT c.id, c.post_id, c.name, c.email, c.date,
	c.content, c.is_deleted, c.is_spam, c.is_approved
	FROM comment AS c
	JOIN post AS p ON p.id = c.post_id`

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS comment(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		date DATETIME NOT NULL,
		content TEXT NOT NULL,
		is_deleted BOOLEAN NOT NULL,
		is_spam BOOLEAN NOT NULL,
		is_approved BOOLEAN NOT NULL,
		FOREIGN KEY (post_id)
		REFERENCES post(id)
	);`

	_, err := r.db.Exec(query)
	assert.NoError(err, "CommentRepository: Error running query: %s", query)
}

func (r *CommentRepository) Create(comment *model.Comment) {
	query := `
	INSERT INTO comment(post_id, name, email, date, content, is_deleted, is_spam, is_approved)
	values(?,?,?,?,?,?,?,?)`

	res, err := r.db.Exec(query,
		comment.PostID,
		comment.Name,
		comment.Email,
		comment.Date,
		comment.Content,
		comment.IsDeleted,
		comment.IsSpam,
		comment.IsApproved,
	)
	assert.NoError(err, "CommentRepository: Error running query: %s", query)

	id, err := res.LastInsertId()
	assert.NoError(err, "CommentRepository: Error getting last insert ID")
	comment.ID = id
}

func (r *CommentRepository) GetCommentsFromPost(postID int64) []model.Comment {
	query := selectComments + `
	WHERE c.post_id = ?
	ORDER BY c.date DESC`
	rows, err := r.db.Query(query, postID)
	assert.NoError(err, "CommentRepository: Error running query: %s", query)

	return parseCommentRows(rows)
}

func parseCommentRows(rows *sql.Rows) []model.Comment {
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Name,
			&comment.Email,
			&comment.Date,
			&comment.Content,
			&comment.IsDeleted,
			&comment.IsSpam,
			&comment.IsApproved,
		)
		assert.NoError(err, "CommentRepository: Error scanning row")
		comments = append(comments, comment)
	}
	return comments
}
