package service

import (
	"context"
	"fmt"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

// InventoryMovementService exposes read operations on inventory_movements.
// Write operations (insert) are handled internally by PurchaseService (and future StockOut service)
// as part of their respective database transactions.
type InventoryMovementService interface {
	// ListAll returns every movement record, most recent first.
	ListAll(ctx context.Context) ([]*entities.InventoryMovement, error)
	// ListByProduct returns movements for a specific product.
	ListByProduct(ctx context.Context, productID int) ([]*entities.InventoryMovement, error)
	// ListByWarehouse returns movements for a specific warehouse.
	ListByWarehouse(ctx context.Context, warehouseID int) ([]*entities.InventoryMovement, error)
}

type inventoryMovementService struct {
	movementRepo repository.InventoryMovementRepository
}

func NewInventoryMovementService(movementRepo repository.InventoryMovementRepository) InventoryMovementService {
	return &inventoryMovementService{movementRepo: movementRepo}
}

func (s *inventoryMovementService) ListAll(ctx context.Context) ([]*entities.InventoryMovement, error) {
	movements, err := s.movementRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("inventory movement service: list all: %w", err)
	}
	return movements, nil
}

func (s *inventoryMovementService) ListByProduct(ctx context.Context, productID int) ([]*entities.InventoryMovement, error) {
	if productID <= 0 {
		return nil, fmt.Errorf("invalid product_id")
	}
	movements, err := s.movementRepo.ListByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("inventory movement service: list by product: %w", err)
	}
	return movements, nil
}

func (s *inventoryMovementService) ListByWarehouse(ctx context.Context, warehouseID int) ([]*entities.InventoryMovement, error) {
	if warehouseID <= 0 {
		return nil, fmt.Errorf("invalid warehouse_id")
	}
	movements, err := s.movementRepo.ListByWarehouseID(ctx, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("inventory movement service: list by warehouse: %w", err)
	}
	return movements, nil
}
