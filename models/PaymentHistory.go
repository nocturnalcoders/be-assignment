package models

import (
	"time"
)

type PaymentHistory struct {
	ID          string    `bson:"_id,omitempty"`
	AccountID   string    `bson:"accountId"`
	Amount      float64   `bson:"amount"`
	Currency    string    `bson:"currency"`
	Description string    `bson:"description"`
	Timestamp   time.Time `bson:"timestamp"`
}
