package domain

import (
	"context"
	"time"
)

type ContractorRepository interface {
	Create(ctx context.Context, contractor Contractor) (string, error)
	GetByID(ctx context.Context, customerID string) (*Customer, error)
	Search(ctx context.Context) ([]Customer, error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, customerID string) error
}

type Contractor struct {
	ContractorID  string
	CompanyName   string
	LegalName     string
	Document      string
	DocumentType  DocumentType
	BusinessPhone string
	BusinessEmail string
	Template      ContractorPlatformTemplate
	Cases         []Case
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedBy     string
	UpdatedAt     time.Time
}

type ContractorPlatformTemplate struct {
	URL       string
	LoginName string
	LoginPass string
	Fields    map[string]map[string]string
}

func NewContractor(legalName, companyName, document, author, businessPhone, businessEmail string, platformTemplate ContractorPlatformTemplate) (Contractor, error) {
	now := time.Now().UTC()
	contractorID, err := uuid.NewRandom()
	if err != nil {
		return Contractor{}, err
	}

	return Contractor{
		ContractorID:  contractorID.String(),
		CompanyName:   companyName,
		LegalName:     legalName,
		Document:      document,
		DocumentType:  CNPJ,
		BusinessPhone: businessPhone,
		BusinessEmail: businessEmail,
		Template:      platformTemplate,
		CreatedBy:     author,
		CreatedAt:     now,
		UpdatedBy:     author,
		UpdatedAt:     now,
	}, nil
}
