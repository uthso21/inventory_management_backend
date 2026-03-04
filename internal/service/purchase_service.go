package service

import (
	"context"
	"errors"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

var (
	ErrWarehouseNotFound  = errors.New("warehouse not found")
	ErrInvalidQuantity    = errors.New("quantity must be greater than zero")
	ErrEmptyPurchaseItems = errors.New("purchase items are required")
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error)
	ValidatePurchaseItems(items []entities.PurchaseItem) error
}

type purchaseService struct {
	purchaseRepo  repository.PurchaseRepository
	warehouseRepo repository.WarehouseRepository
	productRepo   repository.ProductRepository
	movementRepo  repository.MovementRepository
}

func NewPurchaseService(
	purchaseRepo repository.PurchaseRepository,
	warehouseRepo repository.WarehouseRepository,
	productRepo repository.ProductRepository,
	movementRepo repository.MovementRepository,
) PurchaseService {
	return &purchaseService{
		purchaseRepo:  purchaseRepo,
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
		movementRepo:  movementRepo,
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

func (s *purchaseService) CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error) {

	if purchase.Quantity <= 0 {
		return 0, ErrInvalidQuantity
	}

	tx, err := database.BeginTx(ctx)
	if err != nil {
		return 0, err
	}

	// Check warehouse exists
	warehouseExists, err := s.warehouseRepo.ExistsByID(ctx, purchase.WarehouseID)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if !warehouseExists {
		_ = tx.Rollback()
		return 0, ErrWarehouseNotFound
	}

	// Create Purchase
	purchaseID, err := s.purchaseRepo.CreateWithTx(ctx, tx, purchase)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	// Create Movement Log (FIXED HERE)
	movement := &entities.Movement{
		ProductID:   purchase.ProductID,
		WarehouseID: purchase.WarehouseID,
		UserID:      0, // later replace with authenticated user ID
		Type:        "PURCHASE",
		Quantity:    purchase.Quantity,
	}

	if err := s.movementRepo.CreateWithTx(ctx, tx, movement); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return purchaseID, nil
}
