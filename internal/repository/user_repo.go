package repository

import (
	"context"
	"errors"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id int) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entities.User, error)
}

// userRepository is the concrete implementation of UserRepository
type userRepository struct {
	// Add database connection here (e.g., *sql.DB, *gorm.DB)
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	// TODO: Implement database creation logic
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	// TODO: Implement database retrieval logic
	return nil, ErrUserNotFound
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	// TODO: Implement database retrieval logic
	return nil, ErrUserNotFound
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	// TODO: Implement database update logic
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	// TODO: Implement database deletion logic
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*entities.User, error) {
	// TODO: Implement database list logic
	return nil, nil
}
