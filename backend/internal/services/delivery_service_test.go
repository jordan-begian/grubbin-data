package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/models"
)

// mockDeliveryRepo implements repositories.Repository for delivery tests.
type mockDeliveryRepo struct {
	createDeliveryFunc    func(ctx context.Context, delivery *models.Delivery) error
	getDeliveryByIDFunc   func(ctx context.Context, userID, deliveryID string) (*models.Delivery, error)
	getDeliveriesByUserFunc func(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error)
	updateDeliveryFunc    func(ctx context.Context, userID, deliveryID string, fields core.DeliveryUpdateFields) error
	updateDeliveriesFunc  func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error
	deleteDeliveryFunc    func(ctx context.Context, userID, deliveryID string) error
	deleteDeliveriesFunc  func(ctx context.Context, userID string, ids []string) error
}

func (m *mockDeliveryRepo) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	return nil
}

func (m *mockDeliveryRepo) CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
	return nil
}

func (m *mockDeliveryRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

func (m *mockDeliveryRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return nil, nil
}

func (m *mockDeliveryRepo) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	return m.createDeliveryFunc(ctx, delivery)
}

func (m *mockDeliveryRepo) GetDeliveryByID(ctx context.Context, userID, deliveryID string) (*models.Delivery, error) {
	return m.getDeliveryByIDFunc(ctx, userID, deliveryID)
}

func (m *mockDeliveryRepo) GetDeliveriesByUser(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error) {
	return m.getDeliveriesByUserFunc(ctx, userID, startDate, endDate)
}

func (m *mockDeliveryRepo) UpdateDelivery(ctx context.Context, userID, deliveryID string, fields core.DeliveryUpdateFields) error {
	return m.updateDeliveryFunc(ctx, userID, deliveryID, fields)
}

func (m *mockDeliveryRepo) UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
	return m.updateDeliveriesFunc(ctx, userID, updates)
}

func (m *mockDeliveryRepo) DeleteDelivery(ctx context.Context, userID, deliveryID string) error {
	return m.deleteDeliveryFunc(ctx, userID, deliveryID)
}

func (m *mockDeliveryRepo) DeleteDeliveries(ctx context.Context, userID string, ids []string) error {
	return m.deleteDeliveriesFunc(ctx, userID, ids)
}

func TestDeliveryService_CreateDelivery_Success(t *testing.T) {
	var capturedDelivery *models.Delivery

	repo := &mockDeliveryRepo{
		createDeliveryFunc: func(ctx context.Context, delivery *models.Delivery) error {
			capturedDelivery = delivery
			return nil
		},
	}

	service := NewDeliveryService(repo)
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	req := models.CreateDeliveryRequest{
		Start:   now,
		End:     later,
		Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
		Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
		Earnings: models.EarningsRequest{Tip: 300, Base: 500},
		Note:    strPtr("Leave at door"),
	}

	resp, err := service.CreateDelivery(context.Background(), "user-1", req)

	if err != nil {
		t.Fatalf("CreateDelivery failed: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.UserID != "user-1" {
		t.Errorf("UserID = %s, want user-1", resp.UserID)
	}
	if resp.Pickup.Name != "Restaurant" {
		t.Errorf("Pickup.Name = %s, want Restaurant", resp.Pickup.Name)
	}
	if resp.Earnings.Tip != 300 {
		t.Errorf("Earnings.Tip = %d, want 300", resp.Earnings.Tip)
	}
	if capturedDelivery == nil {
		t.Error("Expected CreateDelivery to be called")
	}
}

func TestDeliveryService_CreateDelivery_ValidationErrors(t *testing.T) {
	repo := &mockDeliveryRepo{
		createDeliveryFunc: func(ctx context.Context, delivery *models.Delivery) error {
			t.Error("CreateDelivery should not be called with invalid request")
			return nil
		},
	}

	service := NewDeliveryService(repo)

	// Missing start time
	req := models.CreateDeliveryRequest{
		End:     time.Now().UTC(),
		Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
		Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
		Earnings: models.EarningsRequest{Tip: 300, Base: 500},
	}

	_, err := service.CreateDelivery(context.Background(), "user-1", req)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
	if err.Error() == "" {
		t.Error("Expected error message, got empty")
	}
}

func TestDeliveryService_CreateDelivery_RepoError(t *testing.T) {
	repo := &mockDeliveryRepo{
		createDeliveryFunc: func(ctx context.Context, delivery *models.Delivery) error {
			return errors.New("database connection failed")
		},
	}

	service := NewDeliveryService(repo)
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	req := models.CreateDeliveryRequest{
		Start:   now,
		End:     later,
		Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
		Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
		Earnings: models.EarningsRequest{Tip: 300, Base: 500},
	}

	_, err := service.CreateDelivery(context.Background(), "user-1", req)

	if err == nil {
		t.Fatal("Expected repo error, got nil")
	}
	if err.Error() != "create delivery: database connection failed" {
		t.Errorf("Expected wrapped repo error, got: %v", err)
	}
}

func TestDeliveryService_GetDelivery_Success(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	repo := &mockDeliveryRepo{
		getDeliveryByIDFunc: func(ctx context.Context, userID, deliveryID string) (*models.Delivery, error) {
			if userID != "user-1" {
				t.Errorf("UserID = %s, want user-1", userID)
			}
			if deliveryID != "delivery-1" {
				t.Errorf("DeliveryID = %s, want delivery-1", deliveryID)
			}
			return &models.Delivery{
				ID:       "delivery-1",
				UserID:   "user-1",
				Start:    now,
				End:      later,
				Pickup:   models.PickupLocation{ID: "pickup-1", Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff:  models.DropoffLocation{ID: "dropoff-1", Lat: 40.8, Lon: -73.9},
				Earnings: models.DeliveryEarnings{ID: "earnings-1", Tip: 300, Base: 500},
			}, nil
		},
	}

	service := NewDeliveryService(repo)
	resp, err := service.GetDelivery(context.Background(), "user-1", "delivery-1")

	if err != nil {
		t.Fatalf("GetDelivery failed: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
	if resp.ID != "delivery-1" {
		t.Errorf("ID = %s, want delivery-1", resp.ID)
	}
	if resp.Pickup.Name != "Restaurant" {
		t.Errorf("Pickup.Name = %s, want Restaurant", resp.Pickup.Name)
	}
}

func TestDeliveryService_GetDelivery_NotFound(t *testing.T) {
	repo := &mockDeliveryRepo{
		getDeliveryByIDFunc: func(ctx context.Context, userID, deliveryID string) (*models.Delivery, error) {
			return nil, errors.New("no rows in result set")
		},
	}

	service := NewDeliveryService(repo)
	_, err := service.GetDelivery(context.Background(), "user-1", "nonexistent")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDeliveryService_GetDeliveries_Success(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	repo := &mockDeliveryRepo{
		getDeliveriesByUserFunc: func(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error) {
			if userID != "user-1" {
				t.Errorf("UserID = %s, want user-1", userID)
			}
			return []models.Delivery{
				{
					ID:       "delivery-1",
					UserID:   "user-1",
					Start:    now,
					End:      later,
					Pickup:   models.PickupLocation{ID: "pickup-1", Name: "Restaurant", Lat: 40.7, Lon: -74.0},
					Dropoff:  models.DropoffLocation{ID: "dropoff-1", Lat: 40.8, Lon: -73.9},
					Earnings: models.DeliveryEarnings{ID: "earnings-1", Tip: 300, Base: 500},
				},
			}, nil
		},
	}

	service := NewDeliveryService(repo)
	result, err := service.GetDeliveries(context.Background(), "user-1", nil, nil)

	if err != nil {
		t.Fatalf("GetDeliveries failed: %v", err)
	}
	if result == nil {
		t.Fatal("Expected result, got nil")
	}
	if len(result.Deliveries) != 1 {
		t.Errorf("Deliveries count = %d, want 1", len(result.Deliveries))
	}
	if result.Stats.TotalTime == 0 {
		t.Error("Expected non-zero stats")
	}
}

func TestDeliveryService_GetDeliveries_Empty(t *testing.T) {
	repo := &mockDeliveryRepo{
		getDeliveriesByUserFunc: func(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error) {
			return []models.Delivery{}, nil
		},
	}

	service := NewDeliveryService(repo)
	result, err := service.GetDeliveries(context.Background(), "user-1", nil, nil)

	if err != nil {
		t.Fatalf("GetDeliveries failed: %v", err)
	}
	if len(result.Deliveries) != 0 {
		t.Errorf("Expected empty deliveries, got %d", len(result.Deliveries))
	}
	if result.Stats.TotalTime != 0 {
		t.Error("Expected zero stats for empty list")
	}
}

func TestDeliveryService_UpdateDeliveries_Success(t *testing.T) {
	var capturedUpdates []models.UpdateDeliveryRequest

	repo := &mockDeliveryRepo{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			capturedUpdates = updates
			return nil
		},
	}

	service := NewDeliveryService(repo)
	updates := []models.UpdateDeliveryRequest{
		{ID: "delivery-1", Note: strPtr("Updated note")},
		{ID: "delivery-2", Note: strPtr("Another note")},
	}

	err := service.UpdateDeliveries(context.Background(), "user-1", updates)

	if err != nil {
		t.Fatalf("UpdateDeliveries failed: %v", err)
	}
	if len(capturedUpdates) != 2 {
		t.Errorf("Captured updates count = %d, want 2", len(capturedUpdates))
	}
}

func TestDeliveryService_UpdateDeliveries_ValidationErrors(t *testing.T) {
	repo := &mockDeliveryRepo{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			t.Error("UpdateDeliveries should not be called with invalid request")
			return nil
		},
	}

	service := NewDeliveryService(repo)

	// Missing ID
	updates := []models.UpdateDeliveryRequest{
		{Note: strPtr("Updated note")},
	}

	err := service.UpdateDeliveries(context.Background(), "user-1", updates)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestDeliveryService_UpdateDeliveries_NoFields(t *testing.T) {
	repo := &mockDeliveryRepo{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			t.Error("UpdateDeliveries should not be called with no fields")
			return nil
		},
	}

	service := NewDeliveryService(repo)

	// ID provided but no fields to update
	updates := []models.UpdateDeliveryRequest{
		{ID: "delivery-1"},
	}

	err := service.UpdateDeliveries(context.Background(), "user-1", updates)

	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
}

func TestDeliveryService_UpdateDeliveries_RepoError(t *testing.T) {
	repo := &mockDeliveryRepo{
		updateDeliveriesFunc: func(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
			return errors.New("database connection failed")
		},
	}

	service := NewDeliveryService(repo)
	updates := []models.UpdateDeliveryRequest{
		{ID: "delivery-1", Note: strPtr("Updated note")},
	}

	err := service.UpdateDeliveries(context.Background(), "user-1", updates)

	if err == nil {
		t.Fatal("Expected repo error, got nil")
	}
	if err.Error() != "update deliveries: database connection failed" {
		t.Errorf("Expected wrapped repo error, got: %v", err)
	}
}

func TestDeliveryService_DeleteDeliveries_Success(t *testing.T) {
	var capturedIDs []string

	repo := &mockDeliveryRepo{
		deleteDeliveriesFunc: func(ctx context.Context, userID string, ids []string) error {
			capturedIDs = ids
			return nil
		},
	}

	service := NewDeliveryService(repo)

	err := service.DeleteDeliveries(context.Background(), "user-1", []string{"delivery-1", "delivery-2"})

	if err != nil {
		t.Fatalf("DeleteDeliveries failed: %v", err)
	}
	if len(capturedIDs) != 2 {
		t.Errorf("Captured IDs count = %d, want 2", len(capturedIDs))
	}
}

func TestDeliveryService_DeleteDeliveries_EmptyIDs(t *testing.T) {
	repo := &mockDeliveryRepo{
		deleteDeliveriesFunc: func(ctx context.Context, userID string, ids []string) error {
			t.Error("DeleteDeliveries should not be called with empty IDs")
			return nil
		},
	}

	service := NewDeliveryService(repo)

	err := service.DeleteDeliveries(context.Background(), "user-1", []string{})

	if err == nil {
		t.Fatal("Expected error for empty IDs, got nil")
	}
}

func TestDeliveryService_DeleteDeliveries_RepoError(t *testing.T) {
	repo := &mockDeliveryRepo{
		deleteDeliveriesFunc: func(ctx context.Context, userID string, ids []string) error {
			return errors.New("database connection failed")
		},
	}

	service := NewDeliveryService(repo)

	err := service.DeleteDeliveries(context.Background(), "user-1", []string{"delivery-1"})

	if err == nil {
		t.Fatal("Expected repo error, got nil")
	}
	if err.Error() != "delete deliveries: database connection failed" {
		t.Errorf("Expected wrapped repo error, got: %v", err)
	}
}

// Helper function
func strPtr(s string) *string {
	return &s
}
