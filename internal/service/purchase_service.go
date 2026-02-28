package service

import (
	"context"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type PurchaseService interface {
	CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error)
}

type purchaseService struct {
	purchaseRepo repository.PurchaseRepository
}

func NewPurchaseService(purchaseRepo repository.PurchaseRepository) PurchaseService {
	return &purchaseService{
		purchaseRepo: purchaseRepo,
	}
}

func (s *purchaseService) CreatePurchase(ctx context.Context, purchase *entities.Purchase) (int64, error) {
	purchaseID, err := s.purchaseRepo.Create(ctx, purchase)
	if err != nil {
		return 0, err
	}
	return purchaseID, nil
}
