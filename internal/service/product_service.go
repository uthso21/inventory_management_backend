package usecases

import (
	"errors"

	"github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *entity.Product) error {

	existing, _ := s.repo.GetBySKU(product.SKU)

	if existing != nil {
		return errors.New("sku already exists")
	}

	return s.repo.CreateProduct(product)
}

func (s *ProductService) GetProducts() ([]entity.Product, error) {
	return s.repo.GetProducts()
}

func (s *ProductService) UpdateProduct(product *entity.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
