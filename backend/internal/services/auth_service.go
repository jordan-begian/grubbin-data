// Package services orchestrates business logic and side effects, composing
// pure functions from core with external dependencies like databases and time.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"grubbin-data/backend/internal/config"
	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/models"
	"grubbin-data/backend/internal/repositories"
	"grubbin-data/backend/internal/utilities"
)

// AuthService handles authentication business logic including user registration and login.
type AuthService struct {
	config *config.Config
	repo   repositories.Repository
}

// NewAuthService creates a new authentication service.
func NewAuthService(config *config.Config, repo repositories.Repository) *AuthService {
	return &AuthService{
		config: config,
		repo:   repo,
	}
}

// RegisterUser creates a new user with a validated and hashed password,
// along with their profile and optionally their vehicle.
// All are created atomically within a transaction.
// Returns an error if the password doesn't meet complexity requirements
// or if either the username or profile names are empty.
func (s *AuthService) RegisterUser(
	ctx context.Context,
	username string,
	password string,
	firstName string,
	lastName string,
	vehicleName *string,
	vehicleMPG *float64,
) (*models.User, error) {
	// Validate required fields
	if username == "" {
		return nil, errors.New("username is required")
	}
	if firstName == "" || lastName == "" {
		return nil, errors.New("first name and last name are required")
	}

	// Validate password complexity
	validationErrors := core.ValidatePassword(password)
	if len(validationErrors) > 0 {
		return nil, fmt.Errorf("password validation failed: %v", validationErrors)
	}

	// Hash the password with the server-side pepper
	hash, err := core.HashPassword(password, s.config.PasswordPepper, s.config.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate a sortable ULID for the user
	id, err := utilities.GenerateULID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user id: %w", err)
	}

	now := time.Now().UTC()

	// Create the user model
	user := &models.User{
		ID:       id,
		Username: username,
		Password: &hash,
		Created:  now,
	}

	// Create the profile model (1:1 with user, shares the same ID)
	profile := &models.Profile{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}

	// Create the vehicle model if provided (1:1 with profile, shares the same ID)
	var vehicle *models.Vehicle
	if vehicleName != nil && *vehicleName != "" {
		mpg := 0.0
		if vehicleMPG != nil {
			mpg = *vehicleMPG
		}
		vehicle = &models.Vehicle{
			ID:         id,
			Name:       *vehicleName,
			AverageMPG: mpg,
		}
		profile.Vehicle = vehicle
	}

	// Persist all atomically within a transaction (side effect)
	if err := s.repo.CreateUserWithProfile(ctx, user, profile, vehicle); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Attach profile to user for the response
	user.Profile = profile

	return user, nil
}

// AuthenticateUser verifies a user's credentials and returns the user if valid.
// Returns an error if the user doesn't exist or the password is incorrect.
func (s *AuthService) AuthenticateUser(ctx context.Context, username string, password string) (*models.User, error) {
	// Fetch user from database (side effect)
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if user == nil || user.Password == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password using pure core function
	valid, err := core.VerifyPassword(password, s.config.PasswordPepper, *user.Password)
	if err != nil {
		return nil, fmt.Errorf("password verification failed: %w", err)
	}

	if !valid {
		return nil, errors.New("invalid credentials")
	}

	// Clear password before returning (security best practice)
	user.Password = nil

	return user, nil
}

// ValidatePassword is a convenience method that checks if a password meets complexity requirements.
func (*AuthService) ValidatePassword(password string) []string {
	return core.ValidatePassword(password)
}
