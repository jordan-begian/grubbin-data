package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"grubbin-data/backend/internal/config"
	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/models"
)

// mockRepository implements the repositories.Repository interface for testing.
type mockRepository struct {
	createVehicleFunc         func(ctx context.Context, vehicle *models.Vehicle) error
	createUserWithProfileFunc func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error
	getUserByIDFunc           func(ctx context.Context, id string) (*models.User, error)
	getUserByUsernameFunc     func(ctx context.Context, username string) (*models.User, error)
}

func (m *mockRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	return m.createVehicleFunc(ctx, vehicle)
}

func (m *mockRepository) CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
	return m.createUserWithProfileFunc(ctx, user, profile, vehicle)
}

func (m *mockRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return m.getUserByIDFunc(ctx, id)
}

func (m *mockRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.getUserByUsernameFunc(ctx, username)
}

func (m *mockRepository) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	return nil
}

func (m *mockRepository) GetDeliveryByID(ctx context.Context, userID, deliveryID string) (*models.Delivery, error) {
	return nil, nil
}

func (m *mockRepository) GetDeliveriesByUser(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error) {
	return nil, nil
}

func (m *mockRepository) UpdateDelivery(ctx context.Context, userID, deliveryID string, fields core.DeliveryUpdateFields) error {
	return nil
}

func (m *mockRepository) UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
	return nil
}

func (m *mockRepository) DeleteDelivery(ctx context.Context, userID, deliveryID string) error {
	return nil
}

func (m *mockRepository) DeleteDeliveries(ctx context.Context, userID string, ids []string) error {
	return nil
}

func newTestConfig() *config.Config {
	return &config.Config{
		BcryptCost:     10, // Lower cost for faster tests
		PasswordPepper: "test-pepper",
	}
}

func TestAuthService_RegisterUser_Success(t *testing.T) {
	var capturedUser *models.User
	var capturedProfile *models.Profile
	var capturedVehicle *models.Vehicle

	repo := &mockRepository{
		createUserWithProfileFunc: func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
			// Verify the user has a hashed password, not plaintext
			if user.Password == nil || *user.Password == "" {
				t.Error("Expected password to be hashed, got nil or empty")
			}
			if user.Username != "testuser" {
				t.Errorf("Expected username 'testuser', got %q", user.Username)
			}
			if user.ID == "" {
				t.Error("Expected user ID to be generated, got empty")
			}
			if user.Created.IsZero() {
				t.Error("Expected Created timestamp to be set")
			}
			if profile.FirstName != "John" {
				t.Errorf("Expected first name 'John', got %q", profile.FirstName)
			}
			if profile.LastName != "Doe" {
				t.Errorf("Expected last name 'Doe', got %q", profile.LastName)
			}
			if profile.ID != user.ID {
				t.Errorf("Expected profile ID to match user ID, got %q vs %q", profile.ID, user.ID)
			}
			if vehicle != nil {
				t.Error("Expected no vehicle when none provided")
			}
			capturedUser = user
			capturedProfile = profile
			capturedVehicle = vehicle
			return nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	user, err := service.RegisterUser(context.Background(), "testuser", "Xk9#mP2$vL7@qR4!", "John", "Doe", nil, nil)

	_ = capturedVehicle // suppress unused warning when vehicle is nil

	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", user.Username)
	}
	if user.Password == nil || *user.Password == "" {
		t.Error("Expected hashed password on returned user")
	}
	if user.ID == "" {
		t.Error("Expected user ID to be generated")
	}
	if user.Created.IsZero() {
		t.Error("Expected Created timestamp to be set on returned user")
	}
	if user.Profile == nil {
		t.Fatal("Expected profile to be attached to user")
	}
	if user.Profile.FirstName != "John" {
		t.Errorf("Expected profile first name 'John', got %q", user.Profile.FirstName)
	}
	if user.Profile.LastName != "Doe" {
		t.Errorf("Expected profile last name 'Doe', got %q", user.Profile.LastName)
	}
	if user.Profile.Vehicle != nil {
		t.Error("Expected no vehicle on profile when none provided")
	}
	if capturedUser == nil || capturedProfile == nil {
		t.Error("Expected CreateUserWithProfile to be called")
	}
}

func TestAuthService_RegisterUser_WithVehicle(t *testing.T) {
	var capturedUser *models.User
	var capturedProfile *models.Profile
	var capturedVehicle *models.Vehicle

	vehicleName := "Honda Civic"
	vehicleMPG := 25.5

	repo := &mockRepository{
		createUserWithProfileFunc: func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
			if profile.FirstName != "John" {
				t.Errorf("Expected first name 'John', got %q", profile.FirstName)
			}
			if profile.LastName != "Doe" {
				t.Errorf("Expected last name 'Doe', got %q", profile.LastName)
			}
			if vehicle == nil {
				t.Fatal("Expected vehicle to be created")
			}
			if vehicle.Name != vehicleName {
				t.Errorf("Expected vehicle name %q, got %q", vehicleName, vehicle.Name)
			}
			if vehicle.AverageMPG != vehicleMPG {
				t.Errorf("Expected vehicle MPG %f, got %f", vehicleMPG, vehicle.AverageMPG)
			}
			if vehicle.ID != user.ID {
				t.Errorf("Expected vehicle ID to match user ID, got %q vs %q", vehicle.ID, user.ID)
			}
			capturedUser = user
			capturedProfile = profile
			capturedVehicle = vehicle
			return nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	user, err := service.RegisterUser(context.Background(), "testuser", "Xk9#mP2$vL7@qR4!", "John", "Doe", &vehicleName, &vehicleMPG)

	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Profile == nil {
		t.Fatal("Expected profile to be attached to user")
	}
	if user.Profile.Vehicle == nil {
		t.Fatal("Expected vehicle to be attached to profile")
	}
	if user.Profile.Vehicle.Name != vehicleName {
		t.Errorf("Expected vehicle name %q, got %q", vehicleName, user.Profile.Vehicle.Name)
	}
	if capturedUser == nil || capturedProfile == nil || capturedVehicle == nil {
		t.Error("Expected CreateUserWithProfile to be called with all entities")
	}
}

func TestAuthService_RegisterUser_InvalidPassword(t *testing.T) {
	repo := &mockRepository{
		createUserWithProfileFunc: func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
			t.Error("CreateUserWithProfile should not be called with invalid password")
			return nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err := service.RegisterUser(context.Background(), "testuser", "short", "John", "Doe", nil, nil)

	if err == nil {
		t.Fatal("Expected error for invalid password, got nil")
	}
	if !errors.Is(err, errors.New("password validation failed")) {
		// Error should mention password validation
		if err.Error() == "" || len(err.Error()) < 20 {
			t.Errorf("Expected password validation error, got: %v", err)
		}
	}
}

func TestAuthService_RegisterUser_RepoError(t *testing.T) {
	repo := &mockRepository{
		createUserWithProfileFunc: func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
			return errors.New("database connection failed")
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err := service.RegisterUser(context.Background(), "testuser", "Xk9#mP2$vL7@qR4!", "John", "Doe", nil, nil)

	if err == nil {
		t.Fatal("Expected error for repo failure, got nil")
	}
	if err.Error() != "failed to create user: database connection failed" {
		t.Errorf("Expected wrapped repo error, got: %v", err)
	}
}

func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	password := "MySecureP@ssw0rd123"
	hash, err := core.HashPassword(password, "test-pepper", 10)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	repo := &mockRepository{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			if username != "testuser" {
				t.Errorf("Expected username 'testuser', got %q", username)
			}
			return &models.User{
				ID:       "user-123",
				Username: "testuser",
				Password: &hash,
			}, nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	user, err := service.AuthenticateUser(context.Background(), "testuser", password)

	if err != nil {
		t.Fatalf("AuthenticateUser failed: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", user.Username)
	}
	if user.Password != nil {
		t.Error("Password should be cleared from user struct after successful auth")
	}
}

func TestAuthService_AuthenticateUser_UserNotFound(t *testing.T) {
	repo := &mockRepository{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err := service.AuthenticateUser(context.Background(), "unknownuser", "Xk9#mP2$vL7@qR4!")

	if err == nil {
		t.Fatal("Expected error for user not found, got nil")
	}
	if err.Error() != "authentication failed: user not found" {
		t.Errorf("Expected wrapped auth error, got: %v", err)
	}
}

func TestAuthService_AuthenticateUser_WrongPassword(t *testing.T) {
	password := "MySecureP@ssw0rd123"
	hash, err := core.HashPassword(password, "test-pepper", 10)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	repo := &mockRepository{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return &models.User{
				ID:       "user-123",
				Username: "testuser",
				Password: &hash,
			}, nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err = service.AuthenticateUser(context.Background(), "testuser", "WrongPassword456!")

	if err == nil {
		t.Fatal("Expected error for wrong password, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got: %v", err)
	}
}

func TestAuthService_AuthenticateUser_NilPassword(t *testing.T) {
	repo := &mockRepository{
		getUserByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return &models.User{
				ID:       "user-123",
				Username: "testuser",
				Password: nil, // No password hash
			}, nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err := service.AuthenticateUser(context.Background(), "testuser", "Xk9#mP2$vL7@qR4!")

	if err == nil {
		t.Fatal("Expected error for nil password, got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials' error, got: %v", err)
	}
}

func TestAuthService_RegisterUser_MissingFirstName(t *testing.T) {
	repo := &mockRepository{
		createUserWithProfileFunc: func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
			t.Error("CreateUserWithProfile should not be called with missing first name")
			return nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err := service.RegisterUser(context.Background(), "testuser", "Xk9#mP2$vL7@qR4!", "", "Doe", nil, nil)

	if err == nil {
		t.Fatal("Expected error for missing first name, got nil")
	}
	if err.Error() != "first name and last name are required" {
		t.Errorf("Expected 'first name and last name are required' error, got: %v", err)
	}
}

func TestAuthService_RegisterUser_MissingLastName(t *testing.T) {
	repo := &mockRepository{
		createUserWithProfileFunc: func(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
			t.Error("CreateUserWithProfile should not be called with missing last name")
			return nil
		},
	}

	service := NewAuthService(newTestConfig(), repo)
	_, err := service.RegisterUser(context.Background(), "testuser", "Xk9#mP2$vL7@qR4!", "John", "", nil, nil)

	if err == nil {
		t.Fatal("Expected error for missing last name, got nil")
	}
	if err.Error() != "first name and last name are required" {
		t.Errorf("Expected 'first name and last name are required' error, got: %v", err)
	}
}

func TestAuthService_ValidatePassword(t *testing.T) {
	service := NewAuthService(newTestConfig(), nil)

	errors := service.ValidatePassword("short")
	if len(errors) == 0 {
		t.Error("Expected validation errors for weak password")
	}

	errors = service.ValidatePassword("Xk9#mP2$vL7@qR4!")
	if len(errors) > 0 {
		t.Errorf("Expected no validation errors for strong password, got: %v", errors)
	}
}
