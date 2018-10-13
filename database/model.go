package database

import "time"

type Authentication struct {
	ID               string `gorm:"primary_key"`
	CreatedAt        time.Time
	PhoneNumber      string
	ConfirmationCode string
	ExpiryDate       time.Time
}

type PricingPlan struct {
	ID        string `gorm:"primary_key"`
	AmountUSD float64
	Duration  int64
	MaxUsers  int64
	UpLimit   int64
	DownLimit int64
	PurgeDays int64
}

type Admin struct {
	ID       string `gorm:"primary_key"`
	Login    string `gorm:"unique"`
	Password string
}
