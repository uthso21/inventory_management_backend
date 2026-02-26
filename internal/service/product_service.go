package usecases

import (
	"context"

	"github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *entity.Product) error {
	existing, _ := s.repo.GetBySKU(ctx, product.SKU)
	if existing != nil {
		return repository.ErrProductExists
	}
	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetProducts(ctx context.Context) ([]*entity.Product, error) {
	return s.repo.List(ctx)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
