package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uthso21/inventory_management_backend/internal/database"
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
type userRepository struct{}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, warehouse_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.WarehouseID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, warehouse_id, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &entities.User{}
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.Role, &user.WarehouseID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, warehouse_id, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &entities.User{}
	err := database.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.Role, &user.WarehouseID,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, role = $3, warehouse_id = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err := database.DB.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.Role, user.WarehouseID, user.ID,
	).Scan(&user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*entities.User, error) {
	query := `
		SELECT id, username, email, role, warehouse_id, created_at, updated_at
		FROM users ORDER BY id DESC
	`
	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		u := &entities.User{}
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email,
			&u.Role, &u.WarehouseID,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}
