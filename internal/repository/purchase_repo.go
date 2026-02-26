package repository

import (
	"context"
	"sync"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type PurchaseRepository interface {
	Create(ctx context.Context, purchase *entities.Purchase) error
}

type purchaseRepository struct {
	mu        sync.Mutex
	purchases []entities.Purchase
}

func NewPurchaseRepository() PurchaseRepository {
	return &purchaseRepository{
		purchases: []entities.Purchase{},
	}
}

func (r *purchaseRepository) Create(ctx context.Context, purchase *entities.Purchase) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	purchase.ID = len(r.purchases) + 1
	r.purchases = append(r.purchases, *purchase)

	return nil
}
