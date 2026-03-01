package service

import (
	"context"
	"fmt"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error)
}

type purchaseService struct {
	purchaseRepo  repository.PurchaseRepository
	warehouseRepo repository.WarehouseRepository
}

func NewPurchaseService(
	purchaseRepo repository.PurchaseRepository,
	warehouseRepo repository.WarehouseRepository,
) PurchaseService {
	return &purchaseService{
		purchaseRepo:  purchaseRepo,
		warehouseRepo: warehouseRepo,
	}
}

func (s *purchaseService) CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error) {
	exists, err := s.warehouseRepo.ExistsByID(ctx, purchase.WarehouseID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, fmt.Errorf("warehouse not found")
	}

	purchaseID, err := s.purchaseRepo.Create(ctx, purchase)
	if err != nil {
		return 0, err
	}
	return purchaseID, nil
}
