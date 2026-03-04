package service

import "github.com/uthso21/inventory_management_backend/internal/repository"

type StockOutService struct {
    repo *repository.StockOutRepository
}

func NewStockOutService(repo *repository.StockOutRepository) *StockOutService {
    return &StockOutService{repo: repo}
}

func (s *StockOutService) StockOutProduct(productID, warehouseID, quantity int, reason string) error {
    return s.repo.StockOut(productID, warehouseID, quantity, reason)
}