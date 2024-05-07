package domain

import (
	"context"
	"time"
)

type ProviderRepository interface {
	Create(ctx context.Context, provider Provider) (string, error)
	GetByID(ctx context.Context, providerID string) (*Provider, error)
	Search(ctx context.Context, filters map[string]string) ([]Provider, error)
	Update(ctx context.Context, providerToUpdate User) error
	Delete(ctx context.Context, providerID string) error
}

type Provider struct {
	FirstName     string
	LastName      string
	ContactNumber string
	Email         string
	Address       Address
	Region        int
	Orders        []Order
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedBy     string
	UpdatedAt     time.Time
}
