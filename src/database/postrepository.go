package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gosimple/slug"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/constants"
	"johtotimes.com/src/model"
)

const selectPosts = `
	SELECT p.id, p.title, p.slug, p.img, p.description, 
	p.type, p.filename, p.issue, p.volume, p.permalink,
	p.created_at, p.modified_at, p.hash,
	c.id, c.singular, c.plural, c.slug, c.type
	FROM post AS p
	JOIN category AS c ON c.id = p.category_id`

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS post(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL,
		img TEXT NOT NULL,
		description TEXT NOT NULL,
		issue INTEGER NOT NULL,
		volume INTEGER NOT NULL,
		permalink TEXT NOT NULL,
		filename TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		modified_at DATETIME,
		hash TEXT NOT NULL,
		type CHARACTER(1) NOT NULL,
		category_id INTEGER,
		FOREIGN KEY (category_id)
		REFERENCES category(id)
	);
	`

	_, err := r.db.Exec(query)
	assert.NoError(err, "PostRepository: Error running query: %s", query)
}

func (r *PostRepository) Populate() {

	for t, path := range constants.PostTypePath {
		posts := model.GetPostsFromDirectory(path)
		cr := NewCategoryRepository(r.db)
		for _, post := range posts {
			// Find or create category
			if len(post.Category.Slug) == 0 {
				post.Category = model.Category{
					Singular: fmt.Sprintf("Uncategorized %s", string(t)),
					Plural:   fmt.Sprintf("Uncategorized %s", string(t)),
					Slug:     slug.Make(fmt.Sprintf("Uncategorized %s", string(t))),
					Type:     'C',
				}
			}
			cr.Create(&post.Category)

			// TODO: Create tags
			var tags []model.Category
			// for _, t := range p.Metadata.Tags {
			// 	tag := cr.Create(t, 'T')
			// 	tags = append(tags, tag)
			// }

			post.Tags = tags
			post.Type = t
			post.SetPermalink()

			created := r.Create(post)
			log.Printf(
				"Created post of type %s with slug %s and ID %d\n",
				string(created.Type), created.Slug, created.ID,
			)
		}
	}

}

func (r *PostRepository) Create(post model.Post) *model.Post {
	query := `
	INSERT INTO post(
		title, slug, img, description, category_id,
		type, filename, issue, volume, permalink,
		created_at, modified_at, hash
	)
	values(?,?,?,?,?,?,?,?,?,?,?,?,?)
	`

	assert.NotZero(int(post.Category.ID), "PostRepository: Category ID is zero")
	res, err := r.db.Exec(query,
		post.Title,
		post.Slug,
		post.Img,
		post.Description,
		post.Category.ID,
		post.Type,
		post.FileName,
		post.Issue,
		post.Volume,
		post.Permalink,
		post.CreatedAt,
		post.ModifiedAt,
		post.Hash,
	)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	id, err := res.LastInsertId()
	assert.NoError(err, "PostRepository: Error getting last insert ID")
	post.ID = id

	return &post
}

func (r *PostRepository) Search(search string, offset int, limit int) []model.Post {
	query := selectPosts + `
	WHERE p.title LIKE ? OR p.description LIKE ?
	ORDER BY p.created_at, p.title
	LIMIT ?, ?`
	rows, err := r.db.Query(query, "%"+search+"%", "%"+search+"%", offset, limit)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	return parsePostRows(rows)
}

// Returns posts of a given type ('P', 'N', or 'M')
func (r *PostRepository) GetPage(postType byte, offset int, limit int) []model.Post {
	query := selectPosts + `
	WHERE p.type = ?
	ORDER BY p.created_at, p.title
	LIMIT ?, ?`
	rows, err := r.db.Query(query, postType, offset, limit)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	return parsePostRows(rows)
}

// Returns posts matching category slug.
func (r *PostRepository) GetByCategorySlug(category string, offset int, limit int) []model.Post {
	query := selectPosts + `
	WHERE c.slug = ?
	ORDER BY p.created_at, p.title
	LIMIT ?, ?`
	rows, err := r.db.Query(query, category, offset, limit)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	return parsePostRows(rows)
}

// Returns post matching the given slug. Should always find just 1 row.
func (r *PostRepository) GetBySlug(slug string, postType byte) (*model.Post, error) {
	query := selectPosts + `
	WHERE p.slug = ?
	AND p.type = ?`
	rows, err := r.db.Query(query, slug, postType)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	posts := parsePostRows(rows)
	if len(posts) > 1 {
		log.Fatalf(`Error: Query by slug %q found %d results`, slug, len(posts))
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("PostRepository: Error: Query by slug %q found 0 results", slug)
	}
	return &posts[0], nil
}

func (r *PostRepository) GetByDateAndType(date time.Time, postType byte) (*model.Post, error) {
	query := selectPosts + `
	WHERE p.created_at = ?
	AND p.type = ?
	ORDER BY p.created_at`
	rows, err := r.db.Query(query, date, postType)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	posts := parsePostRows(rows)
	if len(posts) > 1 {
		log.Printf(
			"Warning: Query by date %q and type %q found %d results\n",
			date, string(postType), len(posts),
		)
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("PostRepository: Error: Query by date %q and type %q found 0 results", date, string(postType))
	}
	return &posts[0], nil
}

func parsePostRows(rows *sql.Rows) []model.Post {
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var category model.Category
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Img,
			&post.Description,
			&post.Type,
			&post.FileName,
			&post.Issue,
			&post.Volume,
			&post.Permalink,
			&post.CreatedAt,
			&post.ModifiedAt,
			&post.Hash,
			&category.ID,
			&category.Singular,
			&category.Plural,
			&category.Slug,
			&category.Type,
		)
		assert.NoError(err, "PostRepository: Error scanning row")
		post.Category = category
		posts = append(posts, post)
	}
	return posts
}
