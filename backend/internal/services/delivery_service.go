// Package services orchestrates business logic and side effects, composing
// pure functions from core with external dependencies like databases and time.
package services

import (
	"context"
	"fmt"
	"time"

	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/models"
	"grubbin-data/backend/internal/repositories"
	"grubbin-data/backend/internal/utilities"
)

// DeliveryService handles delivery business logic.
type DeliveryService struct {
	repo repositories.Repository
}

// NewDeliveryService creates a new delivery service.
func NewDeliveryService(repo repositories.Repository) *DeliveryService {
	return &DeliveryService{repo: repo}
}

// CreateDelivery validates and creates a new delivery with related records.
func (s *DeliveryService) CreateDelivery(ctx context.Context, userID string, req models.CreateDeliveryRequest) (*models.DeliveryResponse, error) {
	// Pure validation
	if errors := core.ValidateCreateDelivery(req); len(errors) > 0 {
		return nil, core.FormatValidationError(errors)
	}

	// Generate ULIDs for all entities
	deliveryID, err := utilities.GenerateULID()
	if err != nil {
		return nil, fmt.Errorf("generate delivery id: %w", err)
	}
	pickupID, err := utilities.GenerateULID()
	if err != nil {
		return nil, fmt.Errorf("generate pickup id: %w", err)
	}
	dropoffID, err := utilities.GenerateULID()
	if err != nil {
		return nil, fmt.Errorf("generate dropoff id: %w", err)
	}
	earningsID, err := utilities.GenerateULID()
	if err != nil {
		return nil, fmt.Errorf("generate earnings id: %w", err)
	}

	now := time.Now().UTC()

	// Build the delivery model
	delivery := &models.Delivery{
		ID:                deliveryID,
		UserID:            userID,
		PickupLocationID:  pickupID,
		DropoffLocationID: dropoffID,
		EarningsID:        earningsID,
		Created:           now,
		Start:             req.Start,
		End:               req.End,
		Note:              req.Note,
		Pickup: models.PickupLocation{
			ID:   pickupID,
			Name: req.Pickup.Name,
			Lat:  req.Pickup.Lat,
			Lon:  req.Pickup.Lon,
		},
		Dropoff: models.DropoffLocation{
			ID:  dropoffID,
			Lat: req.Dropoff.Lat,
			Lon: req.Dropoff.Lon,
		},
		Earnings: models.DeliveryEarnings{
			ID:    earningsID,
			Tip:   req.Earnings.Tip,
			Base:  req.Earnings.Base,
			Bonus: req.Earnings.Bonus,
		},
	}

	// Persist (side effect)
	if err := s.repo.CreateDelivery(ctx, delivery); err != nil {
		return nil, fmt.Errorf("create delivery: %w", err)
	}

	// Build response
	resp := core.ToDeliveryResponse(*delivery)
	return &resp, nil
}

// GetDelivery retrieves a single delivery by ID.
func (s *DeliveryService) GetDelivery(ctx context.Context, userID, deliveryID string) (*models.DeliveryResponse, error) {
	delivery, err := s.repo.GetDeliveryByID(ctx, userID, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("get delivery: %w", err)
	}

	resp := core.ToDeliveryResponse(*delivery)
	return &resp, nil
}

// GetDeliveries retrieves all deliveries for a user, optionally filtered by date range.
func (s *DeliveryService) GetDeliveries(ctx context.Context, userID string, startDate, endDate *time.Time) (*models.DeliveryListResponse, error) {
	deliveries, err := s.repo.GetDeliveriesByUser(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("get deliveries: %w", err)
	}

	// Convert to responses
	responses := make([]models.DeliveryResponse, 0, len(deliveries))
	for _, d := range deliveries {
		responses = append(responses, core.ToDeliveryResponse(d))
	}

	// Compute stats (pure function)
	stats := core.ComputeDeliveryStats(responses)

	return &models.DeliveryListResponse{
		Deliveries: responses,
		Stats:      stats,
	}, nil
}

// UpdateDeliveries updates one or more deliveries.
func (s *DeliveryService) UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
	// Validate each update request
	for _, req := range updates {
		if errors := core.ValidateUpdateDelivery(req); len(errors) > 0 {
			return core.FormatValidationError(errors)
		}
	}

	// Persist (side effect)
	if err := s.repo.UpdateDeliveries(ctx, userID, updates); err != nil {
		return fmt.Errorf("update deliveries: %w", err)
	}

	return nil
}

// DeleteDeliveries deletes one or more deliveries.
func (s *DeliveryService) DeleteDeliveries(ctx context.Context, userID string, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("no delivery ids provided")
	}

	// Persist (side effect)
	if err := s.repo.DeleteDeliveries(ctx, userID, ids); err != nil {
		return fmt.Errorf("delete deliveries: %w", err)
	}

	return nil
}
