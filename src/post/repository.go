package post

import (
	"database/sql"
	"fmt"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL,
		img TEXT NOT NULL,
		description TEXT NOT NULL,
		date DATETIME NOT NULL
	);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		fmt.Println(query)
	}
	return err
}

func (r *PostRepository) Create(post Post) (*Post, error) {
	query := `
	INSERT INTO posts(title, slug, img, description, date)
	values(?,?,?,?,?)
	`
	res, err := r.db.Exec(query, post.Title, post.Slug, post.Img, post.Description, post.Date)
	if err != nil {
		fmt.Println(query)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	post.ID = id

	return &post, nil
}

func (r *PostRepository) GetPage(offset int, limit int) ([]Post, error) {
	query := `
	SELECT *
	FROM posts
	ORDER BY date`
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		fmt.Println(query)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Img, &post.Description, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetBySlug(slug string) (*Post, error) {
	row := r.db.QueryRow("SELECT * FROM posts WHERE slug = ? ", slug)

	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Img, &post.Description, &post.Date)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
