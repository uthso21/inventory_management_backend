package service

import (
	"context"
	"errors"

	"github.com/uthso21/inventory_management_backend/internal/database"
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
	if purchase.Quantity <= 0 {
		return 0, errors.New("quantity must be greater than zero")
	}

	tx, err := database.BeginTx(ctx)
	if err != nil {
		return 0, err
	}

	purchaseID, err := s.purchaseRepo.CreateWithTx(ctx, tx, purchase)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return purchaseID, nil
}
