package service

import (
	"context"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

type MovementService interface {
	GetMovements(
		ctx context.Context,
		mType string,
		productID string,
		date string,
	) ([]*entities.Movement, error)
}

type movementService struct {
	movementRepo repository.MovementRepository
}

func NewMovementService(repo repository.MovementRepository) MovementService {
	return &movementService{
		movementRepo: repo,
	}
}

func (s *movementService) GetMovements(
	ctx context.Context,
	mType string,
	productID string,
	date string,
) ([]*entities.Movement, error) {

	return s.movementRepo.GetFiltered(ctx, mType, productID, date)
}
