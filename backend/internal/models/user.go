package models

import "time"

type User struct {
	ID       string     `json:"id" db:"id"`
	Username string     `json:"username" db:"username"`
	Password *string    `json:"-" db:"password_hash"`
	Profile  *Profile   `json:"profile,omitempty" db:"-"`
	Created  time.Time  `json:"created" db:"created_at"`
	Updated  *time.Time `json:"updated,omitempty" db:"updated_at"`
}

type Profile struct {
	ID        string   `json:"id" db:"id"`
	FirstName string   `json:"first_name" db:"first_name"`
	LastName  string   `json:"last_name" db:"last_name"`
	Vehicle   *Vehicle `json:"vehicle,omitempty" db:"-"`
}

type Vehicle struct {
	ID         string  `json:"id" db:"id"`
	Name       string  `json:"name" db:"name"`
	AverageMPG float64 `json:"average_mpg" db:"average_mpg"`
}
