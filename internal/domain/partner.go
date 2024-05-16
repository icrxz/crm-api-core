package domain

import (
	"context"
	"time"
)

type PartnerRepository interface {
	Create(ctx context.Context, provider Partner) (string, error)
	GetByID(ctx context.Context, providerID string) (*Partner, error)
	Search(ctx context.Context, filters map[string]string) ([]Partner, error)
	Update(ctx context.Context, providerToUpdate User) error
	Delete(ctx context.Context, providerID string) error
}

type Partner struct {
	PartnerID       string
	FirstName       string
	LastName        string
	CompanyName     string
	LegalName       string
	PartnerType     EntityType
	Document        string
	DocumentType    DocumentType
	ShippingAddress Address
	BillingAddress  Address
	PersonalPhone   string
	BusinessPhone   string
	PersonalEmail   string
	BusinessEmail   string
	Region          int
	Cases           []Case
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
}
