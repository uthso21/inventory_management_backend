package usecases

import (
	"context"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type WarehouseService interface {
	CreateWarehouse(ctx context.Context, warehouse *entities.Warehouse) error
	ListWarehouses(ctx context.Context) ([]*entities.Warehouse, error)
	UpdateWarehouse(ctx context.Context, warehouse *entities.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
}

type warehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{repo: repo}
}

// CREATE
func (s *warehouseService) CreateWarehouse(ctx context.Context, warehouse *entities.Warehouse) error {
	return s.repo.Create(ctx, warehouse)
}

// LIST
func (s *warehouseService) ListWarehouses(ctx context.Context) ([]*entities.Warehouse, error) {
	return s.repo.List(ctx)
}

// UPDATE
func (s *warehouseService) UpdateWarehouse(ctx context.Context, warehouse *entities.Warehouse) error {
	return s.repo.Update(ctx, warehouse)
}

// DELETE
func (s *warehouseService) DeleteWarehouse(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
