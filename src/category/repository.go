// Package category provides a repository for categories.
package category

import (
	"database/sql"
	"log"

	"github.com/gosimple/slug"
	"johtotimes.com/src/assert"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Migrate() {
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
	assert.NoError(err, "CategoryRepository: Error running query: %s", query)
}

func (r *CategoryRepository) Create(name string, t byte) *Category {
	slug := slug.Make(name)
	if len(slug) == 0 {
		return nil
	}
	assert.NotZero(len(slug), "CategoryRepository: Slug cannot be empty: %s\n", name)

	// Check if category already exists, return it if so.
	cat, err := r.GetBySlug(slug, t)
	if err == nil && cat != nil {
		return cat
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
	assert.NoError(err, "CategoryRepository: Error running query: %s", query)
	id, err := res.LastInsertId()
	assert.NoError(err, "CategoryRepository: Error getting last insert ID")

	category := Category{
		ID:   id,
		Name: name,
		Slug: slug,
		Type: t,
	}
	log.Println("Created category:", category.Name)
	return &category
}

// GetByID searches for a category by its ID.
// If the category is not found, (nil, ErrNotFound) is returned.
func (r *CategoryRepository) GetByID(id int64) (*Category, error) {
	row := r.db.QueryRow("SELECT * FROM category WHERE id = ?", id)

	var category Category
	err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Type)
	return &category, err
}

// GetBySlug searches for a category of a specified type t by its slug.
// If the category is not found, (nil, ErrNotFound) is returned.
func (r *CategoryRepository) GetBySlug(slug string, t byte) (*Category, error) {
	row := r.db.QueryRow("SELECT * FROM category WHERE slug = ? AND type = ?", slug, t)

	var category Category
	err := row.Scan(&category.ID, &category.Name, &category.Slug, &category.Type)
	return &category, err
}
