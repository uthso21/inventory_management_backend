package service

import (
	"context"
	"errors"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error)
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
