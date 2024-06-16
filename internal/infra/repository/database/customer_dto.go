package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CustomerDTO struct {
	CustomerID      string    `db:"customer_id"`
	FirstName       string    `db:"first_name"`
	LastName        string    `db:"last_name"`
	CompanyName     string    `db:"company_name"`
	LegalName       string    `db:"legal_name"`
	CustomerType    string    `db:"customer_type"`
	Document        string    `db:"document"`
	DocumentType    string    `db:"document_type"`
	ShippingAddress string    `db:"shipping_address"`
	ShippingCity    string    `db:"shipping_city"`
	ShippingState   string    `db:"shipping_state"`
	ShippingZipCode string    `db:"shipping_zip_code"`
	ShippingCountry string    `db:"shipping_country"`
	BillingAddress  string    `db:"billing_address"`
	BillingCity     string    `db:"billing_city"`
	BillingState    string    `db:"billing_state"`
	BillingZipCode  string    `db:"billing_zip_code"`
	BillingCountry  string    `db:"billing_country"`
	PersonalPhone   string    `db:"personal_phone"`
	BusinessPhone   string    `db:"business_phone"`
	PersonalEmail   string    `db:"personal_email"`
	BusinessEmail   string    `db:"business_email"`
	OwnerID         *string   `db:"owner_id"`
	CreatedBy       string    `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedBy       string    `db:"updated_by"`
	UpdatedAt       time.Time `db:"updated_at"`
	Active          bool      `db:"active"`
}

func mapCustomerToCustomerDTO(customer domain.Customer) CustomerDTO {
	return CustomerDTO{
		CustomerID:      customer.CustomerID,
		FirstName:       customer.FirstName,
		LastName:        customer.LastName,
		CompanyName:     customer.CompanyName,
		LegalName:       customer.LegalName,
		CustomerType:    string(customer.Type),
		Document:        customer.Document,
		DocumentType:    string(customer.DocumentType),
		ShippingAddress: customer.ShippingAddress.Address,
		ShippingCity:    customer.ShippingAddress.City,
		ShippingState:   customer.ShippingAddress.State,
		ShippingZipCode: customer.ShippingAddress.ZipCode,
		ShippingCountry: customer.ShippingAddress.Country,
		BillingAddress:  customer.BillingAddress.Address,
		BillingCity:     customer.BillingAddress.City,
		BillingState:    customer.BillingAddress.State,
		BillingZipCode:  customer.BillingAddress.ZipCode,
		BillingCountry:  customer.BillingAddress.Country,
		OwnerID:         &customer.OwnerID,
		PersonalPhone:   customer.PersonalContact.PhoneNumber,
		BusinessPhone:   customer.BusinessContact.PhoneNumber,
		PersonalEmail:   customer.PersonalContact.Email,
		BusinessEmail:   customer.BusinessContact.Email,
		CreatedBy:       customer.CreatedBy,
		CreatedAt:       customer.CreatedAt,
		UpdatedBy:       customer.UpdatedBy,
		UpdatedAt:       customer.UpdatedAt,
		Active:          customer.Active,
	}
}

func mapCustomerDTOToCustomer(customerDTO CustomerDTO) domain.Customer {
	var ownerID string
	if customerDTO.OwnerID != nil {
		ownerID = *customerDTO.OwnerID
	}

	return domain.Customer{
		CustomerID:   customerDTO.CustomerID,
		FirstName:    customerDTO.FirstName,
		LastName:     customerDTO.LastName,
		CompanyName:  customerDTO.CompanyName,
		LegalName:    customerDTO.LegalName,
		Type:         domain.EntityType(customerDTO.CustomerType),
		Document:     customerDTO.Document,
		DocumentType: domain.DocumentType(customerDTO.DocumentType),
		ShippingAddress: domain.Address{
			Address: customerDTO.ShippingAddress,
			City:    customerDTO.ShippingCity,
			State:   customerDTO.ShippingState,
			ZipCode: customerDTO.ShippingZipCode,
			Country: customerDTO.ShippingCountry,
		},
		BillingAddress: domain.Address{
			Address: customerDTO.BillingAddress,
			City:    customerDTO.BillingCity,
			State:   customerDTO.BillingState,
			ZipCode: customerDTO.BillingZipCode,
			Country: customerDTO.BillingCountry,
		},
		OwnerID: ownerID,
		PersonalContact: domain.Contact{
			PhoneNumber: customerDTO.PersonalPhone,
			Email:       customerDTO.PersonalEmail,
		},
		BusinessContact: domain.Contact{
			PhoneNumber: customerDTO.BusinessPhone,
			Email:       customerDTO.BusinessEmail,
		},
		CreatedBy: customerDTO.CreatedBy,
		CreatedAt: customerDTO.CreatedAt,
		UpdatedBy: customerDTO.UpdatedBy,
		UpdatedAt: customerDTO.UpdatedAt,
		Active:    customerDTO.Active,
	}
}

func mapCustomerDTOsToCustomers(customerDTOs []CustomerDTO) []domain.Customer {
	customers := make([]domain.Customer, 0, len(customerDTOs))
	for _, customerDTO := range customerDTOs {
		customer := mapCustomerDTOToCustomer(customerDTO)
		customers = append(customers, customer)
	}

	return customers
}
