package service

import (
	"context"
	"errors"

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
}

func NewPurchaseService(
	purchaseRepo repository.PurchaseRepository,
	warehouseRepo repository.WarehouseRepository,
	productRepo repository.ProductRepository,
) PurchaseService {
	return &purchaseService{
		purchaseRepo:  purchaseRepo,
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
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
	warehouseExists, err := s.warehouseRepo.ExistsByID(ctx, purchase.WarehouseID)
	if err != nil {
		return 0, err
	}
	if !warehouseExists {
		return 0, ErrWarehouseNotFound
	}

	purchaseID, err := s.purchaseRepo.Create(ctx, purchase)
	if err != nil {
		return 0, err
	}
	return purchaseID, nil
}
