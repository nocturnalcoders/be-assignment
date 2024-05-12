package models

type PaymentAccount struct {
	ID           string   `bson:"_id,omitempty"`
	UserID       string   `bson:"userId"`
	Type         string   `bson:"type"` // Example: credit, debit, loan
	Balance      float64  `bson:"balance"`
	Currency     string   `bson:"currency"`
	Transactions []string `bson:"transactions,omitempty"`
}
