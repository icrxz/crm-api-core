package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PartnerRepository interface {
	Create(ctx context.Context, partner Partner) (string, error)
	GetByID(ctx context.Context, partnerID string) (*Partner, error)
	Search(ctx context.Context, filters PartnerFilters) ([]Partner, error)
	Update(ctx context.Context, partnerToUpdate Partner) error
	Delete(ctx context.Context, partnerID string) error
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
	BusinessContact Contact
	PersonalContact Contact
	Region          int
	Cases           []Case
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
}

type PartnerFilters struct {
	PartnerID   []string
	Region      []string
	Document    []string
	PartnerType []string
}

func NewPartner(firstName, lastName, companyName, legalName, document, documentType, author string, personalContact, businessContact Contact, shippingAddress, billingAddress Address, region int) (Partner, error) {
	now := time.Now().UTC()

	partnerID, err := uuid.NewUUID()
	if err != nil {
		return Partner{}, err
	}

	var partnerType EntityType
	if firstName != "" {
		partnerType = NATURAL
	} else {
		partnerType = LEGAL
	}

	return Partner{
		PartnerID:       partnerID.String(),
		FirstName:       firstName,
		LastName:        lastName,
		CompanyName:     companyName,
		LegalName:       legalName,
		Document:        document,
		DocumentType:    DocumentType(documentType),
		PartnerType:     partnerType,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		PersonalContact: personalContact,
		BusinessContact: businessContact,
		Region:          region,
		CreatedAt:       now,
		CreatedBy:       author,
		UpdatedAt:       now,
		UpdatedBy:       author,
	}, nil
}
