package models

import "time"

type DeliveryResponse struct {
	ID       string           `json:"id"`
	UserID   string           `json:"user"`
	Created  time.Time        `json:"created"`
	Updated  *time.Time       `json:"updated,omitempty"`
	Start    time.Time        `json:"start"`
	End      time.Time        `json:"end"`
	Pickup   PickupResponse   `json:"pickup"`
	Dropoff  DropoffResponse  `json:"dropoff"`
	Earnings EarningsResponse `json:"earnings"`
	Note     *string          `json:"note,omitempty"`
}

type PickupResponse struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type DropoffResponse struct {
	ID  string  `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type EarningsResponse struct {
	ID    string `json:"id"`
	Tip   int    `json:"tip"`             // cents
	Base  int    `json:"base"`            // cents
	Bonus *int   `json:"bonus,omitempty"` // cents, optional
}

type DeliveryListResponse struct {
	Deliveries []DeliveryResponse `json:"deliveries"`
	Stats      DeliveryStats      `json:"stats"`
}

// DeliveryStats values are computed when building the response and are not stored in the database.
type DeliveryStats struct {
	TotalTime               int      `json:"total_time"` // seconds
	TotalMiles              float64  `json:"total_miles"`
	FuelUsed                *float64 `json:"fuel_used,omitempty"`      // optional
	UsedFuelCost            *int     `json:"used_fuel_cost,omitempty"` // optional, cents
	AverageDeliveryTime     int      `json:"average_delivery_time"`    // seconds
	AverageDeliveryDistance float64  `json:"average_delivery_distance"`
	AverageTip              int      `json:"average_tip"`      // cents
	AverageBasePay          int      `json:"average_base_pay"` // cents
}

type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Created   time.Time  `json:"created"`
	Updated   *time.Time `json:"updated,omitempty"`
	Profile   *Profile   `json:"profile,omitempty"`
}
