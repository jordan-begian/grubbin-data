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
