package post

import (
	"database/sql"
	"log"
	"time"

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

func (r *PostRepository) Migrate() error {
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
		category_id INTEGER NOT NULL,
		FOREIGN KEY (category_id)
		REFERENCES category(id)
	);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		log.Println("Error running query")
		log.Println(query)
	}
	return err
}

func (r *PostRepository) Populate(db *sql.DB) {

	if err := r.Migrate(); err != nil {
		log.Fatal(err)
	}

	for t, path := range internal.PostTypePath {
		posts := getFromDirectory(path)
		cr := category.NewCategoryRepository(r.db)
		for _, p := range posts {
			cat, err := cr.Create(p.Metadata.Category, 'C')
			if err != nil {
				log.Println("Error creating category: " + p.Metadata.Category)
				log.Println(err)
			}
			var tags []*category.Category
			for _, t := range p.Metadata.Tags {
				tag, _ := cr.Create(t, 'T')
				if err != nil {
					log.Println("Error creating tag: " + t)
					log.Println(err)
				}
				tags = append(tags, tag)
			}
			var permalink string
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
			created, err := r.Create(post)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Created post of type %s with slug %s and ID %d\n", string(created.Type), created.Slug, created.ID)
		}
	}

}

func (r *PostRepository) Create(post Post) (*Post, error) {
	query := `
	INSERT INTO post(title, slug, img, description, date, category_id, type, filename, issue, volume, permalink)
	values(?,?,?,?,?,?,?,?,?,?, ?)
	`
	res, err := r.db.Exec(query,
		post.Title,
		post.Slug,
		post.Img,
		post.Description,
		post.Date,
		post.Category.ID,
		post.Type,
		post.FileName,
		post.Issue,
		post.Volume,
		post.Permalink,
	)
	if err != nil {
		log.Println("Error running query.")
		log.Println(query)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	post.ID = id

	return &post, nil
}

// Returns posts of a given type ('P', 'N', or 'M')
func (r *PostRepository) GetPage(postType byte, offset int, limit int) ([]Post, error) {
	query := selectPosts + `
	WHERE p.type = ?
	ORDER BY p.date
	LIMIT ?, ?`
	rows, err := r.db.Query(query, postType, offset, limit)
	if err != nil {
		log.Println(query)
		return nil, err
	}

	return parseRows(rows)
}

// Returns posts matching category slug.
func (r *PostRepository) GetByCategorySlug(category string, offset int, limit int) ([]Post, error) {
	query := selectPosts + `
	WHERE c.slug = ?
	ORDER BY p.date
	LIMIT ?, ?`
	rows, err := r.db.Query(query, category, offset, limit)
	if err != nil {
		log.Println(query)
		return nil, err
	}

	return parseRows(rows)
}

// Returns post matching the given slug. Should always find just 1 row.
func (r *PostRepository) GetBySlug(slug string, postType byte) (*Post, error) {
	query := selectPosts + `
	WHERE p.slug = ?
	AND p.type = ?`
	rows, err := r.db.Query(query, slug, postType)
	if err != nil {
		log.Println(query)
		return nil, err
	}
	posts, err := parseRows(rows)
	if len(posts) > 1 {
		log.Fatalf(`Error: Query by slug %q found %d results`, slug, len(posts))
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
	if err != nil {
		log.Println(query)
		return nil, err
	}

	posts, err := parseRows(rows)
	if len(posts) > 1 {
		log.Fatalf(`Warning: Query by date %q and type %q found %d results`, date, string(postType), len(posts))
	}

	return &posts[0], nil
}

func parseRows(rows *sql.Rows) ([]Post, error) {
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

		post.Category = &category
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
