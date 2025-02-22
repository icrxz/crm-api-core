package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateCustomerDTO struct {
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	CompanyName     string     `json:"company_name"`
	LegalName       string     `json:"legal_name"`
	PartnerType     string     `json:"partner_type"`
	Document        string     `json:"document"`
	DocumentType    string     `json:"document_type"`
	ShippingAddress AddressDTO `json:"shipping"`
	BillingAddress  AddressDTO `json:"billing"`
	PersonalContact ContactDTO `json:"personal_contact"`
	BusinessContact ContactDTO `json:"business_contact"`
	CreatedBy       string     `json:"created_by"`
}

type CustomerDTO struct {
	CustomerID      string     `json:"customer_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	CompanyName     string     `json:"company_name"`
	LegalName       string     `json:"legal_name"`
	Document        string     `json:"document"`
	DocumentType    string     `json:"document_type"`
	ShippingAddress AddressDTO `json:"shipping"`
	BillingAddress  AddressDTO `json:"billing"`
	PersonalContact ContactDTO `json:"personal_contact"`
	BusinessContact ContactDTO `json:"business_contact"`
	Cases           []any      `json:"cases"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       string     `json:"updated_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Active          bool       `json:"active"`
	Region          int        `json:"region"`
}

type UpdateCustomerDTO struct {
	FirstName       *string           `json:"first_name"`
	LastName        *string           `json:"last_name"`
	CompanyName     *string           `json:"company_name"`
	LegalName       *string           `json:"legal_name"`
	PartnerType     *string           `json:"partner_type"`
	Document        *string           `json:"document"`
	DocumentType    *string           `json:"document_type"`
	ShippingAddress *UpdateAddressDTO `json:"shipping"`
	BillingAddress  *UpdateAddressDTO `json:"billing"`
	PersonalContact *UpdateContactDTO `json:"personal_contact"`
	BusinessContact *UpdateContactDTO `json:"business_contact"`
	UpdatedBy       string            `json:"updated_by"`
}

func mapCustomerToCustomerDTO(customer domain.Customer) CustomerDTO {
	return CustomerDTO{
		CustomerID:      customer.CustomerID,
		FirstName:       customer.FirstName,
		LastName:        customer.LastName,
		CompanyName:     customer.CompanyName,
		LegalName:       customer.LegalName,
		Document:        customer.Document,
		DocumentType:    string(customer.DocumentType),
		ShippingAddress: mapAddressToAddressDTO(customer.ShippingAddress),
		BillingAddress:  mapAddressToAddressDTO(customer.BillingAddress),
		PersonalContact: mapContactToContactDTO(customer.PersonalContact),
		BusinessContact: mapContactToContactDTO(customer.BusinessContact),
		CreatedBy:       customer.CreatedBy,
		CreatedAt:       customer.CreatedAt,
		UpdatedBy:       customer.UpdatedBy,
		UpdatedAt:       customer.UpdatedAt,
		Active:          customer.Active,
		Region:          customer.GetRegion(),
	}
}

func mapCreateCustomerDTOToCustomer(customerDTO CreateCustomerDTO) (domain.Customer, error) {
	return domain.NewCustomer(
		customerDTO.FirstName,
		customerDTO.LastName,
		customerDTO.CompanyName,
		customerDTO.LegalName,
		customerDTO.Document,
		customerDTO.DocumentType,
		customerDTO.CreatedBy,
		mapContactDTOToContact(customerDTO.PersonalContact),
		mapContactDTOToContact(customerDTO.BusinessContact),
		mapAddressDTOToAddress(customerDTO.ShippingAddress),
		mapAddressDTOToAddress(customerDTO.BillingAddress),
	)
}

func mapCustomersToCustomerDTOs(customers []domain.Customer) []CustomerDTO {
	customerDTOs := make([]CustomerDTO, 0, len(customers))
	for _, customer := range customers {
		customerDTO := mapCustomerToCustomerDTO(customer)
		customerDTOs = append(customerDTOs, customerDTO)
	}

	return customerDTOs
}

func mapUpdateCustomerDTOToUpdateCustomer(updateCustomerDTO UpdateCustomerDTO) domain.UpdateCustomer {
	var parsedShippingAddress *domain.UpdateAddress
	if updateCustomerDTO.ShippingAddress != nil {
		shippingAddress := mapUpdateAddressDTOToUpdateAddress(*updateCustomerDTO.ShippingAddress)
		parsedShippingAddress = &shippingAddress
	}

	var parsedBillingAddress *domain.UpdateAddress
	if updateCustomerDTO.BillingAddress != nil {
		billingAddress := mapUpdateAddressDTOToUpdateAddress(*updateCustomerDTO.BillingAddress)
		parsedBillingAddress = &billingAddress
	}

	var parsedPersonalContact *domain.UpdateContact
	if updateCustomerDTO.PersonalContact != nil {
		personalContact := mapUpdateContactDTOToUpdateContact(*updateCustomerDTO.PersonalContact)
		parsedPersonalContact = &personalContact
	}

	var parsedBusinessContact *domain.UpdateContact
	if updateCustomerDTO.BusinessContact != nil {
		businessContact := mapUpdateContactDTOToUpdateContact(*updateCustomerDTO.BusinessContact)
		parsedBusinessContact = &businessContact
	}

	return domain.UpdateCustomer{
		FirstName:       updateCustomerDTO.FirstName,
		LastName:        updateCustomerDTO.LastName,
		CompanyName:     updateCustomerDTO.CompanyName,
		LegalName:       updateCustomerDTO.LegalName,
		Document:        updateCustomerDTO.Document,
		DocumentType:    updateCustomerDTO.DocumentType,
		ShippingAddress: parsedShippingAddress,
		BillingAddress:  parsedBillingAddress,
		PersonalContact: parsedPersonalContact,
		BusinessContact: parsedBusinessContact,
		UpdatedBy:       updateCustomerDTO.UpdatedBy,
	}
}
