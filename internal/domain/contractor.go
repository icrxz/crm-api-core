package domain

import (
	"context"
	"time"
)

type ContractorRepository interface {
	Create(ctx context.Context, contractor Contractor) (string, error)
	GetByID(ctx context.Context, customerID string) (*Customer, error)
	List(ctx context.Context) ([]Customer, error)
	Update(ctx context.Context, customer Customer) error
}

type Contractor struct {
	ContractorID string
	LegalName    string
	CNPJ         string
	Template     ContractorPlatformTemplate
	Orders       []Order
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    string
	UpdatedAt    time.Time
}

type ContractorPlatformTemplate struct {
	URL       string
	LoginName string
	LoginPass string
	Fields    map[string]map[string]string
}

func NewContractor(legalName, CNPJ, author string, platformTemplate ContractorPlatformTemplate) Contractor {
	now := time.Now().UTC()

	return Contractor{
		LegalName: legalName,
		CNPJ:      CNPJ,
		Template:  platformTemplate,
		CreatedBy: author,
		CreatedAt: now,
		UpdatedBy: author,
		UpdatedAt: now,
	}
}
