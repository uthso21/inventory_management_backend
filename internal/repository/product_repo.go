package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/uthso21/inventory_management_backend/internal/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductExists   = errors.New("product already exists")
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id int) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entity.Product, error)
}

// productRepository is the concrete implementation of ProductRepository
type productRepository struct {
	mu       sync.RWMutex
	products map[int]*entity.Product
	nextID   int
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository() ProductRepository {
	return &productRepository{
		products: make(map[int]*entity.Product),
		nextID:   1,
	}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate SKU
	for _, p := range r.products {
		if p.SKU == product.SKU {
			return ErrProductExists
		}
	}

	product.ID = r.nextID
	r.nextID++
	r.products[product.ID] = product
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	return p, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.SKU == sku {
			return p, nil
		}
	}
	return nil, ErrProductNotFound
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[product.ID]; !ok {
		return ErrProductNotFound
	}

	r.products[product.ID] = product
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return ErrProductNotFound
	}

	delete(r.products, id)
	return nil
}

func (r *productRepository) List(ctx context.Context) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*entity.Product, 0, len(r.products))
	for _, p := range r.products {
		products = append(products, p)
	}
	return products, nil
}
