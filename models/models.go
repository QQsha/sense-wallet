package models

import "time"

type Transaction struct {
	UserID     string    `json:"user_id" validate:"required"`
	Currency   string    `json:"currency" validate:"required,oneof=EUR USD"`
	Amount     float64   `json:"amount" validate:"required,gte=0"`
	TimePlaced string    `json:"time_placed" validate:"required"`
	Type       string    `json:"type" validate:"required,oneof=deposit withdrawal"`
	Time       time.Time `json:"-"`
}

type Balance struct {
	UserID       string  `json:"user_id"`
	AmountEuro   float64 `json:"euro_balance"`
	AmountUSD float64 `json:"usd_balance"`
}
