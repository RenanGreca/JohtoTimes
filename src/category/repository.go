package category

import (
	"database/sql"
	"log"

	"github.com/gosimple/slug"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS category (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		slug TEXT NOT NULL,
		description TEXT,
		type CHARACTER(1) NOT NULL
	);
	`
	_, err := r.db.Exec(query)
	if err != nil {
		log.Println("Error running query")
		log.Println(query)
	}
	return err
}

func (r *CategoryRepository) Create(name string, t byte) (*Category, error) {
	slug := slug.Make(name)
	// Check if category already exists
	cat, err := r.GetBySlug(slug, t)
	if err == nil {
		return cat, nil
	}

	query := `
	INSERT INTO category(name, slug, type)
	values(?,?,?)
	`
	res, err := r.db.Exec(query,
		name,
		slug,
		t,
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
	category := Category{
		ID:   id,
		Name: name,
		Slug: slug,
		Type: t,
	}

	return &category, nil
}

func (r *CategoryRepository) GetByID(id int64) (*Category, error) {
	row := r.db.QueryRow("SELECT * FROM category WHERE id = ?", id)

	var category Category
	err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Type)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) GetBySlug(slug string, t byte) (*Category, error) {
	row := r.db.QueryRow("SELECT * FROM category WHERE slug = ? AND type = ?", slug, t)

	var category Category
	err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Type)
	if err != nil {
		return nil, err
	}
	return &category, nil
}
