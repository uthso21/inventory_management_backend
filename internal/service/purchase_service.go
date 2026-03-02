package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

var (
	ErrWarehouseNotFound  = errors.New("warehouse not found")
	ErrProductNotFound    = errors.New("product not found")
	ErrInvalidQuantity    = errors.New("quantity must be greater than zero")
	ErrEmptyPurchaseItems = errors.New("purchase items are required")
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, req *entities.CreatePurchaseRequest, userID int) (*entities.Purchase, error)
	GetPurchase(ctx context.Context, id int) (*entities.Purchase, error)
	ListPurchases(ctx context.Context) ([]*entities.Purchase, error)
	ValidatePurchaseItems(items []entities.PurchaseItem) error
}

type purchaseService struct {
	purchaseRepo          repository.PurchaseRepository
	warehouseRepo         repository.WarehouseRepository
	productRepo           repository.ProductRepository
	inventoryMovementRepo repository.InventoryMovementRepository
}

func NewPurchaseService(
	purchaseRepo repository.PurchaseRepository,
	warehouseRepo repository.WarehouseRepository,
	productRepo repository.ProductRepository,
	inventoryMovementRepo repository.InventoryMovementRepository,
) PurchaseService {
	return &purchaseService{
		purchaseRepo:          purchaseRepo,
		warehouseRepo:         warehouseRepo,
		productRepo:           productRepo,
		inventoryMovementRepo: inventoryMovementRepo,
	}
}

func (s *purchaseService) ValidatePurchaseItems(items []entities.PurchaseItem) error {
	if len(items) == 0 {
		return ErrEmptyPurchaseItems
	}
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}
	}
	return nil
}

// CreatePurchase creates a purchase with items and updates stock in a single transaction
// This implements tasks #40, #41, #42, #43, #44, #45
func (s *purchaseService) CreatePurchase(ctx context.Context, req *entities.CreatePurchaseRequest, userID int) (*entities.Purchase, error) {
	// Task #44: Validate positive quantity input
	if err := s.ValidatePurchaseItems(req.Items); err != nil {
		return nil, err
	}

	// Task #43: Validate warehouse existence
	warehouseExists, err := s.warehouseRepo.ExistsByID(ctx, req.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check warehouse: %w", err)
	}
	if !warehouseExists {
		return nil, ErrWarehouseNotFound
	}

	// Task #43: Validate all products exist before starting transaction
	for _, item := range req.Items {
		productExists, err := s.productRepo.ExistsByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to check product: %w", err)
		}
		if !productExists {
			return nil, fmt.Errorf("%w: product_id=%d", ErrProductNotFound, item.ProductID)
		}
	}

	// Task #42: Begin database transaction
	tx, err := database.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create purchase record
	purchase := &entities.Purchase{
		WarehouseID: req.WarehouseID,
		CreatedBy:   userID,
	}

	_, err = s.purchaseRepo.CreateWithTx(ctx, tx, purchase)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to create purchase: %w", err)
	}

	// Process each item
	for _, item := range req.Items {
		// Create purchase item
		purchaseItem := &entities.PurchaseItem{
			PurchaseID: purchase.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
		}

		err = s.purchaseRepo.CreatePurchaseItemWithTx(ctx, tx, purchaseItem)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("failed to create purchase item: %w", err)
		}

		// Task #41: Increment stock
		err = s.productRepo.IncrementStockWithTx(ctx, tx, item.ProductID, item.Quantity)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("failed to increment stock: %w", err)
		}

		// Task #45: Insert inventory movement log
		movement := &entities.InventoryMovement{
			ProductID:     item.ProductID,
			WarehouseID:   req.WarehouseID,
			MovementType:  "purchase",
			Quantity:      item.Quantity,
			ReferenceType: "purchase",
			ReferenceID:   purchase.ID,
			CreatedBy:     userID,
		}

		err = s.inventoryMovementRepo.CreateWithTx(ctx, tx, movement)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("failed to create inventory movement: %w", err)
		}

		purchase.Items = append(purchase.Items, *purchaseItem)
	}

	// Task #42: Commit transaction
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return purchase, nil
}

func (s *purchaseService) GetPurchase(ctx context.Context, id int) (*entities.Purchase, error) {
	return s.purchaseRepo.GetByID(ctx, id)
}

func (s *purchaseService) ListPurchases(ctx context.Context) ([]*entities.Purchase, error) {
	return s.purchaseRepo.List(ctx)
}
