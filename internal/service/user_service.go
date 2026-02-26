package usecases

import (
	"context"
	"errors"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, id int) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context) ([]*entities.User, error)
}

// userService is the concrete implementation of UserService
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *entities.User) error {
	// Add business logic validation here
	if user.Email == "" || user.Username == "" {
		return ErrInvalidInput
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return repository.ErrUserExists
	}

	return s.userRepo.Create(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id int) (*entities.User, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *entities.User) error {
	if user.ID <= 0 {
		return ErrInvalidInput
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return repository.ErrUserNotFound
	}

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context) ([]*entities.User, error) {
	return s.userRepo.List(ctx)
}
