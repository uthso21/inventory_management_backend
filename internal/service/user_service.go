package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req *entities.RegisterRequest) error
	Login(ctx context.Context, req *entities.LoginRequest) (string, error)
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
	return &userService{userRepo: userRepo}
}

// ─────────────────────────────────────
// Auth
// ─────────────────────────────────────

// CreateUser hashes the password with bcrypt before persisting.
// The plain-text password is never stored or logged.
func (s *userService) CreateUser(ctx context.Context, req *entities.RegisterRequest) error {
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return ErrInvalidInput
	}
	if req.Role == "" {
		req.Role = "staff"
	}

	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return repository.ErrUserExists
	}

	// Hash the password — plain text is never passed further
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Immediately overwrite password in request so it can be GC'd
	req.Password = ""

	user := &entities.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		WarehouseID:  req.WarehouseID,
	}

	return s.userRepo.Create(ctx, user)
}

// Login verifies credentials and returns a signed JWT on success.
func (s *userService) Login(ctx context.Context, req *entities.LoginRequest) (string, error) {
	if req.Email == "" || req.Password == "" {
		return "", ErrInvalidInput
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Return generic error so email enumeration is not possible
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return generateJWT(user)
}

// generateJWT creates a signed JWT with user_id, role, warehouse_id, and exp claims.
func generateJWT(user *entities.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme-in-production"
	}

	expiryStr := os.Getenv("JWT_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiry = 24 * time.Hour // default
	}

	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"role":         user.Role,
		"warehouse_id": user.WarehouseID,
		"exp":          time.Now().Add(expiry).Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ─────────────────────────────────────
// CRUD
// ─────────────────────────────────────

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
	_, err := s.userRepo.GetByID(ctx, user.ID)
	if err != nil {
		return err
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

// Keep strconv imported (used for numeric conversions in extended services)
var _ = strconv.Itoa
