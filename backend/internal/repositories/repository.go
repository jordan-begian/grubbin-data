// Package repositories provides database access interfaces and implementations.
package repositories

import (
	"context"
	"fmt"

	"grubbin-data/backend/internal/models"
)

// Repository defines the interface for user data access.
type Repository interface {
	CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	db *DB
}

// NewRepository creates a new PostgresRepository.
func NewRepository(db *DB) Repository {
	return &PostgresRepository{db: db}
}

// CreateVehicle inserts a new vehicle into the database.
func (r *PostgresRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO vehicles (id, name, average_mpg) VALUES ($1, $2, $3)",
		vehicle.ID,
		vehicle.Name,
		vehicle.AverageMPG,
	)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}
	return nil
}

// CreateUserWithProfile creates a user, their profile, and optionally their vehicle atomically within a transaction.
// If any operation fails, the entire transaction is rolled back.
func (r *PostgresRepository) CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		_, err := tx.Exec(ctx,
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

		_, err = tx.Exec(ctx,
			"INSERT INTO profiles (id, first_name, last_name) VALUES ($1, $2, $3)",
			profile.ID,
			profile.FirstName,
			profile.LastName,
		)
		if err != nil {
			return fmt.Errorf("create profile: %w", err)
		}

		if vehicle != nil {
			_, err = tx.Exec(ctx,
				"INSERT INTO vehicles (id, name, average_mpg) VALUES ($1, $2, $3)",
				vehicle.ID,
				vehicle.Name,
				vehicle.AverageMPG,
			)
			if err != nil {
				return fmt.Errorf("create vehicle: %w", err)
			}
		}

		return nil
	})
}

// GetUserByID retrieves a user by their ID, including their profile and vehicle.
// Password is not loaded; the field remains nil.
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	var firstName, lastName, vehicleName *string
	var vehicleMPG *float64

	err := r.db.QueryRow(ctx,
		`SELECT
			u.id, u.username, u.created_at, u.updated_at,
			p.first_name, p.last_name,
			v.name, v.average_mpg
		FROM users u
		LEFT JOIN profiles p ON u.id = p.id
		LEFT JOIN vehicles v ON p.id = v.id
		WHERE u.id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Created,
		&user.Updated,
		&firstName,
		&lastName,
		&vehicleName,
		&vehicleMPG,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	// Attach profile if present
	if firstName != nil && lastName != nil {
		user.Profile = &models.Profile{
			ID:        user.ID,
			FirstName: *firstName,
			LastName:  *lastName,
		}
		// Attach vehicle if present
		if vehicleName != nil {
			user.Profile.Vehicle = &models.Vehicle{
				ID:         user.ID,
				Name:       *vehicleName,
				AverageMPG: 0,
			}
			if vehicleMPG != nil {
				user.Profile.Vehicle.AverageMPG = *vehicleMPG
			}
		}
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username, including their profile and vehicle.
// Password hash is loaded for authentication purposes.
func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	var passwordHash string
	var firstName, lastName, vehicleName *string
	var vehicleMPG *float64

	err := r.db.QueryRow(ctx,
		`SELECT
			u.id, u.username, u.password_hash, u.created_at, u.updated_at,
			p.first_name, p.last_name,
			v.name, v.average_mpg
		FROM users u
		LEFT JOIN profiles p ON u.id = p.id
		LEFT JOIN vehicles v ON p.id = v.id
		WHERE u.username = $1`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&passwordHash,
		&user.Created,
		&user.Updated,
		&firstName,
		&lastName,
		&vehicleName,
		&vehicleMPG,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	user.Password = &passwordHash

	// Attach profile if present
	if firstName != nil && lastName != nil {
		user.Profile = &models.Profile{
			ID:        user.ID,
			FirstName: *firstName,
			LastName:  *lastName,
		}
		// Attach vehicle if present
		if vehicleName != nil {
			user.Profile.Vehicle = &models.Vehicle{
				ID:         user.ID,
				Name:       *vehicleName,
				AverageMPG: 0,
			}
			if vehicleMPG != nil {
				user.Profile.Vehicle.AverageMPG = *vehicleMPG
			}
		}
	}

	return user, nil
}
