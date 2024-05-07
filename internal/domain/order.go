package domain

import (
	"context"
	"time"
)

type OrderRepository interface {
	Create(ctx context.Context, order Order) (string, error)
	GetByID(ctx context.Context, orderID string) (*Order, error)
	List(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, order Order) error
}

type Order struct {
	OrderID           string
	ContractorID      string
	CustomerID        string
	ResponsibleUserID string
	Transactions      []Transaction
	Status            OrderStatus
	DueDate           time.Time
	CreatedBy         string
	CreatedAt         time.Time
	UpdatedBy         string
	UpdatedAt         time.Time
}

type OrderStatus string

const (
	PENDING   OrderStatus = "Pending"
	PROGRESS  OrderStatus = "Progress"
	HOLD      OrderStatus = "Hold"
	COMPLETED OrderStatus = "Completed"
	REJECTED  OrderStatus = "Rejected"
)
