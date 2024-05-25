package post

import (
	"database/sql"
	"log"

	"johtotimes.com/src/category"
	"johtotimes.com/src/internal"
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
	CREATE TABLE IF NOT EXISTS post(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL,
		img TEXT NOT NULL,
		description TEXT NOT NULL,
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
	types := map[byte]string{
		'P': internal.PostsPath,
		'N': internal.NewsPath,
		'M': internal.MailbagPath,
	}

	for t, path := range types {
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
			post := Post{
				Title:       p.Metadata.Title,
				String:      p.Contents,
				Slug:        p.Slug,
				Category:    cat,
				Tags:        tags,
				Img:         p.Metadata.Header,
				Description: p.Metadata.Description,
				Type:        t,
				Date:        p.Date,
			}
			created, err := r.Create(post)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Created post with slug %s and ID %d\n", created.Slug, created.ID)
		}
	}

}

func (r *PostRepository) Create(post Post) (*Post, error) {
	query := `
	INSERT INTO post(title, slug, img, description, date, category_id, type, filename)
	values(?,?,?,?,?,?,?,?)
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
	query := `
	SELECT p.id, p.title, p.slug, p.img, p.description, p.date, p.type, p.filename,
	c.id, c.name, c.slug, c.type
	FROM post AS p
	JOIN category AS c ON c.id = p.category_id
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
	query := `
	SELECT p.id, p.title, p.slug, p.img, p.description, p.date, p.type, p.filename,
	c.id, c.name, c.slug, c.type
	FROM post AS p
	JOIN category AS c ON c.id = p.category_id
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
func (r *PostRepository) GetBySlug(slug string) (*Post, error) {
	query := `
	SELECT p.id, p.title, p.slug, p.img, p.description, p.date, p.type, p.filename,
	c.id, c.name, c.slug, c.type
	FROM post AS p
	JOIN category AS c ON c.id = p.category_id
	WHERE p.slug = ?`
	rows, err := r.db.Query(query, slug)
	if err != nil {
		log.Println(query)
		return nil, err
	}
	posts, err := parseRows(rows)
	if len(posts) > 1 {
		log.Printf(`Warning: Query by slug %q found %d results`, slug, len(posts))
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
