package service

import (
	"context"
	"errors"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type SaleService interface {
	CreateSale(ctx context.Context, sale *entities.Sale) (int64, error)
}

type saleService struct {
	saleRepo     repository.SaleRepository
	productRepo  repository.ProductRepository
	movementRepo repository.MovementRepository
}

func NewSaleService(
	saleRepo repository.SaleRepository,
	productRepo repository.ProductRepository,
	movementRepo repository.MovementRepository,
) SaleService {
	return &saleService{saleRepo, productRepo, movementRepo}
}

func (s *saleService) CreateSale(ctx context.Context, sale *entities.Sale) (int64, error) {

	if sale.Quantity <= 0 {
		return 0, errors.New("invalid quantity")
	}

	tx, err := database.BeginTx(ctx)
	if err != nil {
		return 0, err
	}

	// Check product stock
	product, err := s.productRepo.GetByID(ctx, sale.ProductID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if product.Stock < sale.Quantity {
		tx.Rollback()
		return 0, errors.New("insufficient stock")
	}

	// Reduce stock
	product.Stock -= sale.Quantity
	if err := s.productRepo.Update(ctx, product); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create sale
	id, err := s.saleRepo.CreateWithTx(ctx, tx, sale)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create movement log
	movement := &entities.Movement{
		ProductID:   sale.ProductID,
		WarehouseID: sale.WarehouseID,
		UserID:      0,
		Type:        "SALE",
		Quantity:    sale.Quantity,
	}

	if err := s.movementRepo.CreateWithTx(ctx, tx, movement); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}
