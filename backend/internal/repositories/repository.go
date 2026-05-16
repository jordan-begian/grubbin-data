// Package repositories provides database access interfaces and implementations.
package repositories

import (
	"context"
	"fmt"

	"grubbin-data/backend/internal/models"
)

// Repository defines the interface for user data access.
type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	querier Querier
}

// NewRepository creates a new PostgresRepository.
func NewRepository(querier Querier) Repository {
	return &PostgresRepository{querier: querier}
}

// CreateUser inserts a new user into the database.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.querier.Exec(ctx,
		"INSERT INTO users (id, username, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.ID,
		user.Username,
		user.Password,
		user.Created,
		user.Updated,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by their ID.
// Password is not loaded; the field remains nil.
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := r.querier.QueryRow(ctx,
		"SELECT id, username, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Created,
		&user.Updated,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}
