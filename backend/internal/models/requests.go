package models

import "time"

type CreateDeliveryRequest struct {
	Start    time.Time       `json:"start"`
	End      time.Time       `json:"end"`
	Pickup   PickupRequest   `json:"pickup"`
	Dropoff  DropoffRequest  `json:"dropoff"`
	Earnings EarningsRequest `json:"earnings"`
	Note     *string         `json:"note,omitempty"`
}

type PickupRequest struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type DropoffRequest struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type EarningsRequest struct {
	Tip   int  `json:"tip"`             // cents
	Base  int  `json:"base"`            // cents
	Bonus *int `json:"bonus,omitempty"` // cents, optional
}

type RegisterUserRequest struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	VehicleName  *string  `json:"vehicle_name,omitempty"`
	VehicleMPG   *float64 `json:"vehicle_mpg,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateDeliveryRequest represents a partial update for a delivery.
// Only non-nil fields are updated; nil fields are ignored.
// Setting a nullable field to an explicit null pointer clears it.
type UpdateDeliveryRequest struct {
	ID       string            `json:"id"`
	Start    *time.Time        `json:"start,omitempty"`
	End      *time.Time        `json:"end,omitempty"`
	Pickup   *PickupRequest    `json:"pickup,omitempty"`
	Dropoff  *DropoffRequest   `json:"dropoff,omitempty"`
	Earnings *EarningsRequest  `json:"earnings,omitempty"`
	Note     *string           `json:"note,omitempty"`
}
