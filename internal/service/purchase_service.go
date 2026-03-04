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

func (s *purchaseService) CreatePurchase(ctx context.Context, req *entities.CreatePurchaseRequest, userID int) (*entities.Purchase, error) {

	if err := s.ValidatePurchaseItems(req.Items); err != nil {
		return nil, err
	}

	warehouseExists, err := s.warehouseRepo.ExistsByID(ctx, req.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check warehouse: %w", err)
	}
	if !warehouseExists {
		return nil, ErrWarehouseNotFound
	}

	for _, item := range req.Items {
		productExists, err := s.productRepo.ExistsByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if !productExists {
			return nil, ErrProductNotFound
		}
	}

	tx, err := database.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	purchase := &entities.Purchase{
		WarehouseID: req.WarehouseID,
		CreatedBy:   userID,
	}

	_, err = s.purchaseRepo.CreateWithTx(ctx, tx, purchase)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range req.Items {

		purchaseItem := &entities.PurchaseItem{
			PurchaseID: purchase.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
		}

		err = s.purchaseRepo.CreatePurchaseItemWithTx(ctx, tx, purchaseItem)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		err = s.productRepo.IncrementStockWithTx(ctx, tx, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

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
			tx.Rollback()
			return nil, err
		}

		purchase.Items = append(purchase.Items, *purchaseItem)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return purchase, nil
}

func (s *purchaseService) GetPurchase(ctx context.Context, id int) (*entities.Purchase, error) {
	return s.purchaseRepo.GetByID(ctx, id)
}

func (s *purchaseService) ListPurchases(ctx context.Context) ([]*entities.Purchase, error) {
	return s.purchaseRepo.List(ctx)
}
