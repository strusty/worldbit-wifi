package database

import "time"

type Authentication struct {
	ID               string `gorm:"primary_key"`
	CreatedAt        time.Time
	PhoneNumber      string
	ConfirmationCode string
	ExpiryDate       time.Time
}
