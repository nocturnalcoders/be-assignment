package models

import (
	"time"
)

type Transaction struct {
	ID          string    `bson:"_id,omitempty"`
	Amount      float64   `bson:"amount"`
	Description string    `bson:"description"`
	Currency    string    `bson:"currency"`
	Timestamp   time.Time `bson:"timestamp"`
	AccountID   string    `bson:"accountId"`
	HistoryID   string    `bson:"historyId,omitempty"`
}
