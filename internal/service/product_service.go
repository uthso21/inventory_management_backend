package usecases

import (
	"context"

	"github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

// ProductService defines the interface for product business logic
type ProductService interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProduct(ctx context.Context, id int) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id int) error
	ListProducts(ctx context.Context) ([]*entity.Product, error)
}

// productService is the concrete implementation of ProductService
type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new instance of ProductService
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, product *entity.Product) error {
	if product.Name == "" || product.SKU == "" {
		return ErrInvalidInput
	}

	existing, _ := s.repo.GetBySKU(ctx, product.SKU)
	if existing != nil {
		return repository.ErrProductExists
	}

	return s.repo.Create(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id int) (*entity.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	return s.repo.GetByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *entity.Product) error {
	if product.ID <= 0 {
		return ErrInvalidInput
	}

	existing, err := s.repo.GetByID(ctx, product.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return repository.ErrProductNotFound
	}

	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]*entity.Product, error) {
	return s.repo.List(ctx)
}
