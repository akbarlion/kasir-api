package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	query := "SELECT id, name, description FROM category"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (repo *CategoryRepository) Create(data *models.Category) error {
	query := "INSERT INTO category (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, data.Name, data.Description).Scan(&data.ID)
	return err
}

// Category GetByID
func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM category WHERE id = $1"

	var c models.Category
	err := repo.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("Category tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update Category
func (repo *CategoryRepository) Update(category *models.Category) error {
	query := "UPDATE category SET name = $1, description = $2 WHERE id = $3"

	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Category tidak ditemukan")
	}

	return nil
}

// Delete Category
func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM category WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Kategori tidak ditemukan")
	}
	return nil
}
