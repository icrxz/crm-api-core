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
	CustomerID      string
	OwnerID         string
	FirstName       string
	LastName        string
	CompanyName     string
	LegalName       string
	Document        string
	DocumentType    string
	Type            EntityType
	ShippingAddress Address
	BillingAddress  Address
	PersonalPhone   string
	BusinessPhone   string
	PersonalEmail   string
	BusinessEmail   string
	Cases           []Case
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
}

type DocumentType string

const (
	CPF  DocumentType = "CPF"
	CNPJ DocumentType = "CNPJ"
	RG   DocumentType = "RG"
)

type EntityType string

const (
	NATURAL EntityType = "Natural"
	LEGAL   EntityType = "Legal"
)
