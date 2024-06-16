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
	Active          bool
}

type EditPartner struct {
	FirstName       *string
	LastName        *string
	CompanyName     *string
	LegalName       *string
	PartnerType     *EntityType
	Document        *string
	DocumentType    *DocumentType
	ShippingAddress *Address
	BillingAddress  *Address
	BusinessContact *Contact
	PersonalContact *Contact
	Region          *int
	Active          *bool
	UpdatedBy       string
}

type PartnerFilters struct {
	PartnerID   []string
	Region      []string
	Document    []string
	PartnerType []string
	Active      *bool
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
		Active:          true,
	}, nil
}

func (p *Partner) MergeUpdate(updatePartner EditPartner) {
	p.UpdatedBy = updatePartner.UpdatedBy
	p.UpdatedAt = time.Now().UTC()

	if updatePartner.FirstName != nil {
		p.FirstName = *updatePartner.FirstName
	}

	if updatePartner.LastName != nil {
		p.LastName = *updatePartner.LastName
	}

	if updatePartner.CompanyName != nil {
		p.CompanyName = *updatePartner.CompanyName
	}

	if updatePartner.LegalName != nil {
		p.LegalName = *updatePartner.LegalName
	}

	if updatePartner.Document != nil {
		p.Document = *updatePartner.Document
	}

	if updatePartner.DocumentType != nil {
		p.DocumentType = *updatePartner.DocumentType
	}

	if updatePartner.ShippingAddress != nil {
		p.ShippingAddress = *updatePartner.ShippingAddress
	}

	if updatePartner.BillingAddress != nil {
		p.BillingAddress = *updatePartner.BillingAddress
	}

	if updatePartner.BusinessContact != nil {
		p.BusinessContact = *updatePartner.BusinessContact
	}

	if updatePartner.PersonalContact != nil {
		p.PersonalContact = *updatePartner.PersonalContact
	}

	if updatePartner.Region != nil {
		p.Region = *updatePartner.Region
	}

	if updatePartner.Active != nil {
		p.Active = *updatePartner.Active
	}
}
