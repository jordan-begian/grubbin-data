package core

import (
	"testing"
	"time"

	"grubbin-data/backend/internal/models"
)

func TestValidateCreateDelivery(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	tests := []struct {
		name    string
		req     models.CreateDeliveryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: models.CreateDeliveryRequest{
				Start:   now,
				End:     later,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: 500},
			},
			wantErr: false,
		},
		{
			name: "missing start time",
			req: models.CreateDeliveryRequest{
				End:     later,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: 500},
			},
			wantErr: true,
			errMsg:  "start time is required",
		},
		{
			name: "missing end time",
			req: models.CreateDeliveryRequest{
				Start:   now,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: 500},
			},
			wantErr: true,
			errMsg:  "end time is required",
		},
		{
			name: "end before start",
			req: models.CreateDeliveryRequest{
				Start:   later,
				End:     now,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: 500},
			},
			wantErr: true,
			errMsg:  "end time must be after start time",
		},
		{
			name: "missing pickup name",
			req: models.CreateDeliveryRequest{
				Start:   now,
				End:     later,
				Pickup:  models.PickupRequest{Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: 500},
			},
			wantErr: true,
			errMsg:  "pickup name is required",
		},
		{
			name: "negative tip",
			req: models.CreateDeliveryRequest{
				Start:   now,
				End:     later,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: -100, Base: 500},
			},
			wantErr: true,
			errMsg:  "tip cannot be negative",
		},
		{
			name: "negative base pay",
			req: models.CreateDeliveryRequest{
				Start:   now,
				End:     later,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: -500},
			},
			wantErr: true,
			errMsg:  "base pay cannot be negative",
		},
		{
			name: "negative bonus",
			req: models.CreateDeliveryRequest{
				Start:   now,
				End:     later,
				Pickup:  models.PickupRequest{Name: "Restaurant", Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffRequest{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsRequest{Tip: 300, Base: 500, Bonus: ptrInt(-50)},
			},
			wantErr: true,
			errMsg:  "bonus cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateCreateDelivery(tt.req)
			hasErrors := len(errors) > 0

			if hasErrors != tt.wantErr {
				t.Errorf("ValidateCreateDelivery() errors = %v, wantErr %v", errors, tt.wantErr)
			}

			if tt.wantErr && tt.errMsg != "" {
				found := false
				for _, e := range errors {
					if e == tt.errMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ValidateCreateDelivery() expected error %q, got %v", tt.errMsg, errors)
				}
			}
		})
	}
}

func TestValidateUpdateDelivery(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)

	tests := []struct {
		name    string
		req     models.UpdateDeliveryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update - single field",
			req: models.UpdateDeliveryRequest{
				ID:   "delivery-1",
				Note: strPtr("Updated note"),
			},
			wantErr: false,
		},
		{
			name: "valid update - multiple fields",
			req: models.UpdateDeliveryRequest{
				ID:    "delivery-1",
				Start: &now,
				End:   &later,
				Note:  strPtr("Updated note"),
			},
			wantErr: false,
		},
		{
			name: "valid update - clear note",
			req: models.UpdateDeliveryRequest{
				ID:   "delivery-1",
				Note: strPtr(""),
			},
			wantErr: false,
		},
		{
			name: "missing id",
			req: models.UpdateDeliveryRequest{
				Note: strPtr("Updated note"),
			},
			wantErr: true,
			errMsg:  "delivery id is required",
		},
		{
			name: "no fields to update",
			req: models.UpdateDeliveryRequest{
				ID: "delivery-1",
			},
			wantErr: true,
			errMsg:  "at least one field must be provided for update",
		},
		{
			name: "end before start",
			req: models.UpdateDeliveryRequest{
				ID:    "delivery-1",
				Start: &later,
				End:   &now,
			},
			wantErr: true,
			errMsg:  "end time must be after start time",
		},
		{
			name: "empty pickup name",
			req: models.UpdateDeliveryRequest{
				ID:     "delivery-1",
				Pickup: &models.PickupRequest{Name: ""},
			},
			wantErr: true,
			errMsg:  "pickup name cannot be empty",
		},
		{
			name: "negative tip in update",
			req: models.UpdateDeliveryRequest{
				ID:       "delivery-1",
				Earnings: &models.EarningsRequest{Tip: -100, Base: 500},
			},
			wantErr: true,
			errMsg:  "tip cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateUpdateDelivery(tt.req)
			hasErrors := len(errors) > 0

			if hasErrors != tt.wantErr {
				t.Errorf("ValidateUpdateDelivery() errors = %v, wantErr %v", errors, tt.wantErr)
			}

			if tt.wantErr && tt.errMsg != "" {
				found := false
				for _, e := range errors {
					if e == tt.errMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ValidateUpdateDelivery() expected error %q, got %v", tt.errMsg, errors)
				}
			}
		})
	}
}

func TestBuildDeliveryUpdateFields(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(30 * time.Minute)
	bonus := 100

	tests := []struct {
		name           string
		req            models.UpdateDeliveryRequest
		wantDelivery   bool
		wantPickup     bool
		wantDropoff    bool
		wantEarnings   bool
	}{
		{
			name: "only note update",
			req: models.UpdateDeliveryRequest{
				ID:   "delivery-1",
				Note: strPtr("Updated note"),
			},
			wantDelivery: true,
			wantPickup:   false,
			wantDropoff:  false,
			wantEarnings: false,
		},
		{
			name: "only pickup update",
			req: models.UpdateDeliveryRequest{
				ID:     "delivery-1",
				Pickup: &models.PickupRequest{Name: "New Restaurant", Lat: 41.0, Lon: -74.5},
			},
			wantDelivery: false,
			wantPickup:   true,
			wantDropoff:  false,
			wantEarnings: false,
		},
		{
			name: "only earnings update",
			req: models.UpdateDeliveryRequest{
				ID:       "delivery-1",
				Earnings: &models.EarningsRequest{Tip: 400, Base: 600, Bonus: &bonus},
			},
			wantDelivery: false,
			wantPickup:   false,
			wantDropoff:  false,
			wantEarnings: true,
		},
		{
			name: "all fields updated",
			req: models.UpdateDeliveryRequest{
				ID:      "delivery-1",
				Start:   &now,
				End:     &later,
				Note:    strPtr("Note"),
				Pickup:  &models.PickupRequest{Name: "Place", Lat: 40.0, Lon: -73.0},
				Dropoff: &models.DropoffRequest{Lat: 41.0, Lon: -74.0},
				Earnings: &models.EarningsRequest{Tip: 300, Base: 500},
			},
			wantDelivery: true,
			wantPickup:   true,
			wantDropoff:  true,
			wantEarnings: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := BuildDeliveryUpdateFields(tt.req)

			if (fields.Delivery != nil) != tt.wantDelivery {
				t.Errorf("Delivery fields: got %v, want %v", fields.Delivery != nil, tt.wantDelivery)
			}
			if (fields.Pickup != nil) != tt.wantPickup {
				t.Errorf("Pickup fields: got %v, want %v", fields.Pickup != nil, tt.wantPickup)
			}
			if (fields.Dropoff != nil) != tt.wantDropoff {
				t.Errorf("Dropoff fields: got %v, want %v", fields.Dropoff != nil, tt.wantDropoff)
			}
			if (fields.Earnings != nil) != tt.wantEarnings {
				t.Errorf("Earnings fields: got %v, want %v", fields.Earnings != nil, tt.wantEarnings)
			}
		})
	}
}

func TestComputeDeliveryStats(t *testing.T) {
	t.Run("empty list returns zero stats", func(t *testing.T) {
		stats := ComputeDeliveryStats([]models.DeliveryResponse{})
		if stats.TotalTime != 0 {
			t.Errorf("TotalTime = %d, want 0", stats.TotalTime)
		}
		if stats.TotalMiles != 0 {
			t.Errorf("TotalMiles = %f, want 0", stats.TotalMiles)
		}
	})

	t.Run("single delivery", func(t *testing.T) {
		start := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)

		deliveries := []models.DeliveryResponse{
			{
				Start: start,
				End:   end,
				Pickup: models.PickupResponse{Lat: 40.7128, Lon: -74.0060},
				Dropoff: models.DropoffResponse{Lat: 40.7580, Lon: -73.9855},
				Earnings: models.EarningsResponse{Tip: 300, Base: 500},
			},
		}

		stats := ComputeDeliveryStats(deliveries)

		if stats.TotalTime != 1800 {
			t.Errorf("TotalTime = %d, want 1800", stats.TotalTime)
		}
		if stats.AverageDeliveryTime != 1800 {
			t.Errorf("AverageDeliveryTime = %d, want 1800", stats.AverageDeliveryTime)
		}
		if stats.AverageTip != 300 {
			t.Errorf("AverageTip = %d, want 300", stats.AverageTip)
		}
		if stats.AverageBasePay != 500 {
			t.Errorf("AverageBasePay = %d, want 500", stats.AverageBasePay)
		}
	})

	t.Run("multiple deliveries", func(t *testing.T) {
		start1 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		end1 := time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)
		start2 := time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
		end2 := time.Date(2024, 1, 1, 13, 20, 0, 0, time.UTC)

		deliveries := []models.DeliveryResponse{
			{
				Start: start1, End: end1,
				Pickup: models.PickupResponse{Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffResponse{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsResponse{Tip: 300, Base: 500},
			},
			{
				Start: start2, End: end2,
				Pickup: models.PickupResponse{Lat: 40.7, Lon: -74.0},
				Dropoff: models.DropoffResponse{Lat: 40.8, Lon: -73.9},
				Earnings: models.EarningsResponse{Tip: 500, Base: 700},
			},
		}

		stats := ComputeDeliveryStats(deliveries)

		expectedTotalTime := 1800 + 1200 // 30min + 20min
		if stats.TotalTime != expectedTotalTime {
			t.Errorf("TotalTime = %d, want %d", stats.TotalTime, expectedTotalTime)
		}
		if stats.AverageTip != 400 {
			t.Errorf("AverageTip = %d, want 400", stats.AverageTip)
		}
		if stats.AverageBasePay != 600 {
			t.Errorf("AverageBasePay = %d, want 600", stats.AverageBasePay)
		}
	})
}

func TestToDeliveryResponse(t *testing.T) {
	now := time.Now().UTC()
	bonus := 100

	delivery := models.Delivery{
		ID:      "delivery-1",
		UserID:  "user-1",
		Created: now,
		Start:   now,
		End:     now.Add(30 * time.Minute),
		Pickup: models.PickupLocation{
			ID:   "pickup-1",
			Name: "Restaurant",
			Lat:  40.7,
			Lon:  -74.0,
		},
		Dropoff: models.DropoffLocation{
			ID:  "dropoff-1",
			Lat: 40.8,
			Lon: -73.9,
		},
		Earnings: models.DeliveryEarnings{
			ID:    "earnings-1",
			Tip:   300,
			Base:  500,
			Bonus: &bonus,
		},
		Note: strPtr("Leave at door"),
	}

	resp := ToDeliveryResponse(delivery)

	if resp.ID != delivery.ID {
		t.Errorf("ID = %s, want %s", resp.ID, delivery.ID)
	}
	if resp.UserID != delivery.UserID {
		t.Errorf("UserID = %s, want %s", resp.UserID, delivery.UserID)
	}
	if resp.Pickup.Name != delivery.Pickup.Name {
		t.Errorf("Pickup.Name = %s, want %s", resp.Pickup.Name, delivery.Pickup.Name)
	}
	if resp.Earnings.Tip != delivery.Earnings.Tip {
		t.Errorf("Earnings.Tip = %d, want %d", resp.Earnings.Tip, delivery.Earnings.Tip)
	}
	if resp.Note == nil || *resp.Note != *delivery.Note {
		t.Errorf("Note = %v, want %v", resp.Note, delivery.Note)
	}
}

func TestFormatValidationError(t *testing.T) {
	t.Run("nil for empty errors", func(t *testing.T) {
		err := FormatValidationError([]string{})
		if err != nil {
			t.Errorf("FormatValidationError([]) = %v, want nil", err)
		}
	})

	t.Run("error for non-empty errors", func(t *testing.T) {
		err := FormatValidationError([]string{"field required"})
		if err == nil {
			t.Error("FormatValidationError() = nil, want error")
		}
	})
}

func TestHaversineDistance(t *testing.T) {
	// NYC to LA is approximately 2451 miles
	distance := haversineDistance(40.7128, -74.0060, 34.0522, -118.2437)

	// Allow 5% tolerance
	expected := 2451.0
	tolerance := expected * 0.05

	if distance < expected-tolerance || distance > expected+tolerance {
		t.Errorf("haversineDistance(NYC, LA) = %f, want ~%f (±5%%)", distance, expected)
	}
}

func TestRoundTo2(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{1.234, 1.23},
		{1.235, 1.24},
		{1.236, 1.24},
		{0.0, 0.0},
		{100.999, 101.0},
	}

	for _, tt := range tests {
		got := roundTo2(tt.input)
		if got != tt.want {
			t.Errorf("roundTo2(%f) = %f, want %f", tt.input, got, tt.want)
		}
	}
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
