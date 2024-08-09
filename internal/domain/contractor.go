package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ContractorRepository interface {
	Create(ctx context.Context, contractor Contractor) (string, error)
	GetByID(ctx context.Context, contractorID string) (*Contractor, error)
	Search(ctx context.Context, filters ContractorFilters) (PagingResult[Contractor], error)
	Update(ctx context.Context, contractor Contractor) error
	Delete(ctx context.Context, contractorID string) error
}

type Contractor struct {
	ContractorID    string
	CompanyName     string
	LegalName       string
	Document        string
	DocumentType    DocumentType
	BusinessContact Contact
	Template        ContractorPlatformTemplate
	Cases           []Case
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
	Active          bool
}

type UpdateContractor struct {
	CompanyName     *string
	LegalName       *string
	Document        *string
	DocumentType    *DocumentType
	BusinessContact *Contact
	UpdatedBy       string
}

func (c *Contractor) MergeUpdate(newContractor UpdateContractor) {
	now := time.Now().UTC()
	c.UpdatedAt = now
	c.UpdatedBy = newContractor.UpdatedBy

	if newContractor.CompanyName != nil {
		c.CompanyName = *newContractor.CompanyName
	}

	if newContractor.LegalName != nil {
		c.LegalName = *newContractor.LegalName
	}

	if newContractor.Document != nil {
		c.Document = *newContractor.Document
	}

	if newContractor.DocumentType != nil {
		c.DocumentType = *newContractor.DocumentType
	}

	if newContractor.BusinessContact != nil {
		c.BusinessContact = *newContractor.BusinessContact
	}
}

type ContractorFilters struct {
	ContractorID []string
	CompanyName  []string
	Document     []string
	Active       *bool
	PagingFilter
}

func NewContractor(legalName, companyName, document, author string, businessContact Contact, platformTemplate ContractorPlatformTemplate) (Contractor, error) {
	now := time.Now().UTC()

	contractorID, err := uuid.NewRandom()
	if err != nil {
		return Contractor{}, err
	}

	return Contractor{
		ContractorID:    contractorID.String(),
		CompanyName:     companyName,
		LegalName:       legalName,
		Document:        document,
		DocumentType:    CNPJ,
		BusinessContact: businessContact,
		Template:        platformTemplate,
		CreatedBy:       author,
		CreatedAt:       now,
		UpdatedBy:       author,
		UpdatedAt:       now,
		Active:          true,
	}, nil
}
