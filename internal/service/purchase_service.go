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
	CreatePurchase(ctx context.Context, purchase *entities.Purchase) error
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

func (s *purchaseService) CreatePurchase(ctx context.Context, purchase *entities.Purchase) error {
	warehouseExists, err := s.warehouseRepo.ExistsByID(ctx, purchase.WarehouseID)
	if err != nil {
		return err
	}
	if !warehouseExists {
		return ErrWarehouseNotFound
	}

	return s.purchaseRepo.Create(ctx, purchase)
}
