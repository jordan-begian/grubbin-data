package models

import "time"

type Delivery struct {
	ID                string           `json:"id" db:"id"`
	UserID            string           `json:"user" db:"user_id"`
	PickupLocationID  string           `json:"-" db:"pickup_location_id"`
	DropoffLocationID string           `json:"-" db:"dropoff_location_id"`
	EarningsID        string           `json:"-" db:"earnings_id"`
	Created           time.Time        `json:"created" db:"created_at"`
	Updated           *time.Time       `json:"updated,omitempty" db:"updated_at"`
	Start             time.Time        `json:"start" db:"start_time"`
	End               time.Time        `json:"end" db:"end_time"`
	Pickup            PickupLocation   `json:"pickup" db:"-"`
	Dropoff           DropoffLocation  `json:"dropoff" db:"-"`
	Earnings          DeliveryEarnings `json:"earnings" db:"-"`
	Note              *string          `json:"note,omitempty" db:"note"`
}

type PickupLocation struct {
	ID   string  `json:"id" db:"id"`
	Name string  `json:"name" db:"name"`
	Lat  float64 `json:"lat" db:"lat"`
	Lon  float64 `json:"lon" db:"lon"`
}

type DropoffLocation struct {
	ID  string  `json:"id" db:"id"`
	Lat float64 `json:"lat" db:"lat"`
	Lon float64 `json:"lon" db:"lon"`
}

type DeliveryEarnings struct {
	ID    string `json:"id" db:"id"`
	Tip   int    `json:"tip" db:"tip_cents"`
	Base  int    `json:"base" db:"base_cents"`
	Bonus *int   `json:"bonus,omitempty" db:"bonus_cents"`
}
