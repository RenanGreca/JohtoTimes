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
		date DATETIME NOT NULL,
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

	posts := getFromDirectory(internal.PostsPath)
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
			Date:        p.Date,
		}
		created, err := r.Create(post)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created post with slug %s and ID %d\n", created.Slug, created.ID)
	}

}

func (r *PostRepository) Create(post Post) (*Post, error) {
	query := `
	INSERT INTO post(title, slug, img, description, date, category_id)
	values(?,?,?,?,?,?)
	`
	res, err := r.db.Exec(query,
		post.Title,
		post.Slug,
		post.Img,
		post.Description,
		post.Date,
		post.Category.ID,
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

func (r *PostRepository) GetPage(offset int, limit int) ([]Post, error) {
	query := `
	SELECT *
	FROM post
	ORDER BY date
	LIMIT ?, ?`
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		log.Println(query)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var categoryID int64
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Img, &post.Description, &post.Date, &categoryID)

		cr := category.NewCategoryRepository(r.db)
		cat, err := cr.GetByID(categoryID)
		post.Category = cat
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) GetBySlug(slug string) (*Post, error) {
	row := r.db.QueryRow("SELECT * FROM post WHERE slug = ? ", slug)

	var post Post
	var categoryID int64
	err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Img, &post.Description, &post.Date, &categoryID)

	cr := category.NewCategoryRepository(r.db)
	cat, err := cr.GetByID(categoryID)
	post.Category = cat
	if err != nil {
		return nil, err
	}
	return &post, nil
}
