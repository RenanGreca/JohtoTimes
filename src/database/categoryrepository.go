package database

import (
	"database/sql"
	"log"

	"johtotimes.com/src/assert"
	"johtotimes.com/src/constants"
	"johtotimes.com/src/model"
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
		singular TEXT NOT NULL,
		plural TEXT NOT NULL,
		slug TEXT NOT NULL,
		description TEXT,
		type CHARACTER(1) NOT NULL
	);
	`
	_, err := r.db.Exec(query)
	assert.NoError(err, "CategoryRepository: Error running query: %s", query)
}

func (r *CategoryRepository) Populate() {
	categories := model.GetCategoriesFromDirectory(constants.CategoriesPath)
	for _, category := range categories {
		r.Create(&category)
	}
}

func (r *CategoryRepository) Create(category *model.Category) {
	// Check if category already exists, return it if so.
	cat, err := r.GetBySlug(category.Slug, category.Type)
	if err == nil && cat != nil {
		category = cat
		return
	}

	query := `
	INSERT INTO category(singular, plural, slug, type)
	values(?,?,?,?)
	`
	res, err := r.db.Exec(query,
		category.Singular,
		category.Plural,
		category.Slug,
		category.Type,
	)
	assert.NoError(err, "CategoryRepository: Error running query: %s", query)
	id, err := res.LastInsertId()
	assert.NoError(err, "CategoryRepository: Error getting last insert ID")
	category.ID = id
	log.Printf(
		"Created category of type %s with slug %s and ID %d\n",
		string(category.Type), category.Slug, category.ID,
	)
}

// GetByID searches for a category by its ID.
// If the category is not found, (nil, ErrNotFound) is returned.
func (r *CategoryRepository) GetByID(id int64) (*model.Category, error) {
	row := r.db.QueryRow("SELECT * FROM category WHERE id = ?", id)

	var category model.Category
	err := row.Scan(&category.ID, &category.Singular, &category.Plural, &category.Slug, &category.Type)
	return &category, err
}

// GetBySlug searches for a category of a specified type t by its slug.
// If the category is not found, (nil, ErrNotFound) is returned.
func (r *CategoryRepository) GetBySlug(slug string, t byte) (*model.Category, error) {
	row := r.db.QueryRow("SELECT * FROM category WHERE slug = ? AND type = ?", slug, t)

	var category model.Category
	err := row.Scan(&category.ID, &category.Singular, &category.Plural, &category.Slug, &category.Type)
	return &category, err
}
