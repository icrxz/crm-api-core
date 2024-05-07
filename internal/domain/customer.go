package domain

import (
	"context"
	"time"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer Customer) (string, error)
	GetByID(ctx context.Context, customerID string) (*Customer, error)
	List(ctx context.Context) ([]Customer, error)
	Update(ctx context.Context, customer Customer) error
}

type Customer struct {
	CustomerID    string
	FirstName     string
	LastName      string
	Document      string
	Type          CustomerType
	Address       Address
	ContactNumber string
	Email         string
	Orders        []Order
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedBy     string
	UpdatedAt     time.Time
}

type CustomerType string

const (
	NATURAL CustomerType = "Natural"
	LEGAL   CustomerType = "Legal"
)
