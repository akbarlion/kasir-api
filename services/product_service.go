package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

// Product By ID
func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// Update (By ID tentunya)
func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

// Delete (juga By ID)
func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
