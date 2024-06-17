package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer Customer) (string, error)
	GetByID(ctx context.Context, customerID string) (*Customer, error)
	Search(ctx context.Context, filters CustomerFilters) ([]Customer, error)
	Update(ctx context.Context, customer Customer) error
	Delete(ctx context.Context, customerID string) error
}

type Customer struct {
	CustomerID      string
	OwnerID         string
	FirstName       string
	LastName        string
	CompanyName     string
	LegalName       string
	Document        string
	DocumentType    DocumentType
	Type            EntityType
	ShippingAddress Address
	BillingAddress  Address
	BusinessContact Contact
	PersonalContact Contact
	Cases           []Case
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
	Active          bool
	Region          int
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

type CustomerFilters struct {
	CustomerID   []string
	OwnerID      []string
	CustomerType []string
	Document     []string
	Active       bool
}

func NewCustomer(firstName, lastName, companyName, legalName, document, documentType, author string, personalContact, businessContact Contact, shippingAddress, billingAddress Address) (Customer, error) {
	now := time.Now().UTC()

	customerID, err := uuid.NewUUID()
	if err != nil {
		return Customer{}, err
	}

	var customerType EntityType
	if firstName != "" {
		customerType = NATURAL
	} else {
		customerType = LEGAL
	}

	region := regions[shippingAddress.State]

	return Customer{
		CustomerID:      customerID.String(),
		Type:            customerType,
		FirstName:       firstName,
		LastName:        lastName,
		CompanyName:     companyName,
		LegalName:       legalName,
		Document:        document,
		DocumentType:    DocumentType(documentType),
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		PersonalContact: personalContact,
		BusinessContact: businessContact,
		CreatedAt:       now,
		CreatedBy:       author,
		UpdatedAt:       now,
		UpdatedBy:       author,
		Active:          true,
		Region:          region,
	}, nil
}

type UpdateCustomer struct {
	FirstName       *string
	LastName        *string
	CompanyName     *string
	LegalName       *string
	Document        *string
	DocumentType    *string
	ShippingAddress *Address
	BillingAddress  *Address
	BusinessContact *Contact
	PersonalContact *Contact
	UpdatedBy       string
}

func (c *Customer) MergeUpdate(updateCustomer UpdateCustomer) {
	c.UpdatedBy = updateCustomer.UpdatedBy
	c.UpdatedAt = time.Now().UTC()

	if updateCustomer.FirstName != nil {
		c.FirstName = *updateCustomer.FirstName
	}

	if updateCustomer.LastName != nil {
		c.LastName = *updateCustomer.LastName
	}

	if updateCustomer.CompanyName != nil {
		c.CompanyName = *updateCustomer.CompanyName
	}

	if updateCustomer.LegalName != nil {
		c.LegalName = *updateCustomer.LegalName
	}

	if updateCustomer.Document != nil {
		c.Document = *updateCustomer.Document
	}

	if updateCustomer.DocumentType != nil {
		c.DocumentType = DocumentType(*updateCustomer.DocumentType)
	}

	if updateCustomer.ShippingAddress != nil {
		c.ShippingAddress = *updateCustomer.ShippingAddress
		c.Region = regions[updateCustomer.ShippingAddress.State]
	}

	if updateCustomer.BillingAddress != nil {
		c.BillingAddress = *updateCustomer.BillingAddress
	}

	if updateCustomer.BusinessContact != nil {
		c.BusinessContact = *updateCustomer.BusinessContact
	}

	if updateCustomer.PersonalContact != nil {
		c.PersonalContact = *updateCustomer.PersonalContact
	}
}
