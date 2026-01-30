package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

// GetByID
func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

// Update (By ID tentunya)
func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

// Delete (By ID jugaaa)
func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
