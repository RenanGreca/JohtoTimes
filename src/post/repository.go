package post

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/category"
	"johtotimes.com/src/internal"
)

const selectPosts = `
	SELECT p.id, p.title, p.slug, p.img, p.description, p.date,
	p.type, p.filename, p.issue, p.volume, p.permalink,
	c.id, c.name, c.slug, c.type
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
		date DATETIME NOT NULL,
		type CHARACTER(1) NOT NULL,
		category_id INTEGER,
		FOREIGN KEY (category_id)
		REFERENCES category(id)
	);
	`

	_, err := r.db.Exec(query)
	assert.NoError(err, "PostRepository: Error running query: %s", query)
}

func (r *PostRepository) Populate(db *sql.DB) {

	r.Migrate()

	for t, path := range internal.PostTypePath {
		posts := getFromDirectory(path)
		cr := category.NewCategoryRepository(r.db)
		for _, p := range posts {
			// Create category
			var cat *category.Category
			if len(p.Metadata.Category) > 0 {
				cat = cr.Create(p.Metadata.Category, 'C')
			} else {
				cat = cr.Create(fmt.Sprintf("Uncategorized %s", string(t)), 'C')
			}

			// TODO: Create tags
			// for _, t := range p.Metadata.Tags {
			// 	tag := cr.Create(t, 'T')
			// 	tags = append(tags, tag)
			// }

			// Create permalink
			var permalink string

			var tags []*category.Category
			switch t {
			case 'P':
				permalink = "/posts/" + cat.Slug + "/" + p.Slug
			case 'M':
				permalink = "/mailbag/" + p.Slug
			case 'N':
				permalink = "/news/" + p.Slug
			case 'I':
				permalink = "/issues/" + p.Slug
			}
			assert.NotNil(permalink, "PostRepository: Error creating permalink")
			post := Post{
				Title:       p.Metadata.Title,
				String:      p.Contents,
				Slug:        p.Slug,
				Category:    cat,
				Tags:        tags,
				Img:         p.Metadata.Header,
				Description: p.Metadata.Description,
				Issue:       p.Metadata.Issue,
				Volume:      p.Metadata.Volume,
				Permalink:   permalink,
				Type:        t,
				Date:        p.Date,
			}
			created := r.Create(post)
			log.Printf(
				"Created post of type %s with slug %s and ID %d\n",
				string(created.Type), created.Slug, created.ID,
			)
		}
	}

}

func (r *PostRepository) Create(post Post) *Post {
	query := `
	INSERT INTO post(title, slug, img, description, date, category_id, type, filename, issue, volume, permalink)
	values(?,?,?,?,?,?,?,?,?,?,?)
	`

	var category_id int64
	if post.Category != nil {
		category_id = post.Category.ID
	}
	res, err := r.db.Exec(query,
		post.Title,
		post.Slug,
		post.Img,
		post.Description,
		post.Date,
		category_id,
		post.Type,
		post.FileName,
		post.Issue,
		post.Volume,
		post.Permalink,
	)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	id, err := res.LastInsertId()
	assert.NoError(err, "PostRepository: Error getting last insert ID")
	post.ID = id

	return &post
}

// Returns posts of a given type ('P', 'N', or 'M')
func (r *PostRepository) GetPage(postType byte, offset int, limit int) []Post {
	query := selectPosts + `
	WHERE p.type = ?
	ORDER BY p.date
	LIMIT ?, ?`
	rows, err := r.db.Query(query, postType, 0, 10)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	return parseRows(rows)
}

func printQuery(query string, args ...interface{}) {
	query = strings.ReplaceAll(query, "?", "%q")
	log.Printf(query, args...)
}

// Returns posts matching category slug.
func (r *PostRepository) GetByCategorySlug(category string, offset int, limit int) []Post {
	query := selectPosts + `
	WHERE c.slug = ?
	ORDER BY p.date
	LIMIT ?, ?`
	rows, err := r.db.Query(query, category, offset, limit)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	return parseRows(rows)
}

// Returns post matching the given slug. Should always find just 1 row.
func (r *PostRepository) GetBySlug(slug string, postType byte) (*Post, error) {
	query := selectPosts + `
	WHERE p.slug = ?
	AND p.type = ?`
	rows, err := r.db.Query(query, slug, postType)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	posts := parseRows(rows)
	if len(posts) > 1 {
		log.Fatalf(`Error: Query by slug %q found %d results`, slug, len(posts))
	}
	if len(posts) == 0 {
		return nil, fmt.Errorf("PostRepository: Error: Query by slug %q found 0 results", slug)
	}
	return &posts[0], nil
}

func (r *PostRepository) GetByDateAndType(date time.Time, postType byte) (*Post, error) {
	query := selectPosts + `
	WHERE p.date = ?
	AND p.type = ?
	ORDER BY p.date
	`
	rows, err := r.db.Query(query, date, postType)
	assert.NoError(err, "PostRepository: Error running query: %s", query)

	posts := parseRows(rows)
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

func parseRows(rows *sql.Rows) []Post {
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var category category.Category
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Slug,
			&post.Img,
			&post.Description,
			&post.Date,
			&post.Type,
			&post.FileName,
			&post.Issue,
			&post.Volume,
			&post.Permalink,
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Type,
		)
		assert.NoError(err, "PostRepository: Error scanning row")
		post.Category = &category
		posts = append(posts, post)
	}
	return posts
}
