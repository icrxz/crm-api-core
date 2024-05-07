package domain

import "time"

type TransactionRepository interface{}

type Transaction struct {
	TransactionID string
	Type          TransactionType
	Value         float64
	OrderID       string
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedBy     string
	UpdatedAt     time.Time
}

type TransactionType int

const (
	INCOMING TransactionType = iota
	OUTGOING
)
