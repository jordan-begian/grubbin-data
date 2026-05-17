// Package core provides pure business logic functions with no side effects.
// All functions in this package are deterministic and free of I/O,
// making them trivial to unit test.
package core

import (
	"fmt"
	"math"

	"grubbin-data/backend/internal/models"
)

// ValidateCreateDelivery validates a CreateDeliveryRequest and returns
// a list of validation error messages. Returns nil if valid.
func ValidateCreateDelivery(req models.CreateDeliveryRequest) []string {
	var errors []string

	if req.Start.IsZero() {
		errors = append(errors, "start time is required")
	}
	if req.End.IsZero() {
		errors = append(errors, "end time is required")
	}
	if !req.Start.IsZero() && !req.End.IsZero() && req.End.Before(req.Start) {
		errors = append(errors, "end time must be after start time")
	}
	if req.Pickup.Name == "" {
		errors = append(errors, "pickup name is required")
	}
	if req.Pickup.Lat == 0 && req.Pickup.Lon == 0 {
		errors = append(errors, "pickup coordinates are required")
	}
	if req.Dropoff.Lat == 0 && req.Dropoff.Lon == 0 {
		errors = append(errors, "dropoff coordinates are required")
	}
	if req.Earnings.Tip < 0 {
		errors = append(errors, "tip cannot be negative")
	}
	if req.Earnings.Base < 0 {
		errors = append(errors, "base pay cannot be negative")
	}
	if req.Earnings.Bonus != nil && *req.Earnings.Bonus < 0 {
		errors = append(errors, "bonus cannot be negative")
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// ValidateUpdateDelivery validates an UpdateDeliveryRequest and returns
// a list of validation error messages. Returns nil if valid.
// The ID field is required; at least one other field must be provided.
func ValidateUpdateDelivery(req models.UpdateDeliveryRequest) []string {
	var errors []string

	if req.ID == "" {
		errors = append(errors, "delivery id is required")
	}

	hasUpdate := req.Start != nil ||
		req.End != nil ||
		req.Pickup != nil ||
		req.Dropoff != nil ||
		req.Earnings != nil ||
		req.Note != nil

	if !hasUpdate {
		errors = append(errors, "at least one field must be provided for update")
	}

	if req.Start != nil && req.End != nil && req.End.Before(*req.Start) {
		errors = append(errors, "end time must be after start time")
	}

	if req.Pickup != nil && req.Pickup.Name == "" {
		errors = append(errors, "pickup name cannot be empty")
	}

	if req.Earnings != nil {
		if req.Earnings.Tip < 0 {
			errors = append(errors, "tip cannot be negative")
		}
		if req.Earnings.Base < 0 {
			errors = append(errors, "base pay cannot be negative")
		}
		if req.Earnings.Bonus != nil && *req.Earnings.Bonus < 0 {
			errors = append(errors, "bonus cannot be negative")
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// DeliveryUpdateFields holds the dynamic SQL update maps for a delivery
// and its related tables. Only non-nil fields are included.
type DeliveryUpdateFields struct {
	Delivery   map[string]any
	Pickup     map[string]any
	Dropoff    map[string]any
	Earnings   map[string]any
}

// BuildDeliveryUpdateFields extracts non-nil fields from an UpdateDeliveryRequest
// into separate maps for each table. Returns nil maps for tables with no updates.
func BuildDeliveryUpdateFields(req models.UpdateDeliveryRequest) DeliveryUpdateFields {
	fields := DeliveryUpdateFields{}

	// Delivery table fields
	if req.Start != nil || req.End != nil || req.Note != nil {
		fields.Delivery = make(map[string]any)
		if req.Start != nil {
			fields.Delivery["start_time"] = *req.Start
		}
		if req.End != nil {
			fields.Delivery["end_time"] = *req.End
		}
		if req.Note != nil {
			fields.Delivery["note"] = *req.Note
		}
	}

	// Pickup table fields
	if req.Pickup != nil {
		fields.Pickup = make(map[string]any)
		if req.Pickup.Name != "" {
			fields.Pickup["name"] = req.Pickup.Name
		}
		if req.Pickup.Lat != 0 || req.Pickup.Lon != 0 {
			fields.Pickup["lat"] = req.Pickup.Lat
			fields.Pickup["lon"] = req.Pickup.Lon
		}
	}

	// Dropoff table fields
	if req.Dropoff != nil {
		fields.Dropoff = make(map[string]any)
		if req.Dropoff.Lat != 0 || req.Dropoff.Lon != 0 {
			fields.Dropoff["lat"] = req.Dropoff.Lat
			fields.Dropoff["lon"] = req.Dropoff.Lon
		}
	}

	// Earnings table fields
	if req.Earnings != nil {
		fields.Earnings = make(map[string]any)
		fields.Earnings["tip_cents"] = req.Earnings.Tip
		fields.Earnings["base_cents"] = req.Earnings.Base
		if req.Earnings.Bonus != nil {
			fields.Earnings["bonus_cents"] = *req.Earnings.Bonus
		}
	}

	return fields
}

// ComputeDeliveryStats calculates aggregate statistics from a list of deliveries.
// Returns zero-value stats if the list is empty.
func ComputeDeliveryStats(deliveries []models.DeliveryResponse) models.DeliveryStats {
	if len(deliveries) == 0 {
		return models.DeliveryStats{}
	}

	var totalTime int
	var totalMiles float64
	var totalTip, totalBase int
	var fuelUsed float64
	var hasFuelData bool

	for _, d := range deliveries {
		// Total time in seconds
		duration := int(d.End.Sub(d.Start).Seconds())
		totalTime += duration

		// Distance calculation (Haversine formula)
		miles := haversineDistance(d.Pickup.Lat, d.Pickup.Lon, d.Dropoff.Lat, d.Dropoff.Lon)
		totalMiles += miles

		// Earnings
		totalTip += d.Earnings.Tip
		totalBase += d.Earnings.Base

		// Fuel data (if available)
		if miles > 0 {
			hasFuelData = true
			fuelUsed += miles / 25.0 // Default 25 MPG assumption
		}
	}

	count := float64(len(deliveries))
	stats := models.DeliveryStats{
		TotalTime:               totalTime,
		TotalMiles:              roundTo2(totalMiles),
		AverageDeliveryTime:     int(float64(totalTime) / count),
		AverageDeliveryDistance: roundTo2(totalMiles / count),
		AverageTip:              int(float64(totalTip) / count),
		AverageBasePay:          int(float64(totalBase) / count),
	}

	if hasFuelData {
		fuel := roundTo2(fuelUsed)
		stats.FuelUsed = &fuel
		cost := int(fuelUsed * 350) // $3.50/gallon assumption in cents
		stats.UsedFuelCost = &cost
	}

	return stats
}

// haversineDistance calculates the distance in miles between two lat/lon points.
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusMiles = 3958.8

	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degToRad(lat1))*math.Cos(degToRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMiles * c
}

// degToRad converts degrees to radians.
func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// roundTo2 rounds a float64 to 2 decimal places.
func roundTo2(val float64) float64 {
	return math.Round(val*100) / 100
}

// ToDeliveryResponse converts a models.Delivery to a models.DeliveryResponse.
// This is a pure mapping function with no side effects.
func ToDeliveryResponse(d models.Delivery) models.DeliveryResponse {
	return models.DeliveryResponse{
		ID:      d.ID,
		UserID:  d.UserID,
		Created: d.Created,
		Updated: d.Updated,
		Start:   d.Start,
		End:     d.End,
		Pickup: models.PickupResponse{
			ID:   d.Pickup.ID,
			Name: d.Pickup.Name,
			Lat:  d.Pickup.Lat,
			Lon:  d.Pickup.Lon,
		},
		Dropoff: models.DropoffResponse{
			ID:  d.Dropoff.ID,
			Lat: d.Dropoff.Lat,
			Lon: d.Dropoff.Lon,
		},
		Earnings: models.EarningsResponse{
			ID:    d.Earnings.ID,
			Tip:   d.Earnings.Tip,
			Base:  d.Earnings.Base,
			Bonus: d.Earnings.Bonus,
		},
		Note: d.Note,
	}
}

// FormatValidationError formats a list of validation errors into a single message.
func FormatValidationError(errors []string) error {
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("validation failed: %v", errors)
}
