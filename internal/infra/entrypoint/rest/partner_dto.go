package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreatePartnerDTO struct {
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

type PartnerDTO struct {
	PartnerID       string     `json:"partner_id"`
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
	Region          int        `json:"region"`
	Cases           []any      `json:"cases"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       string     `json:"updated_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Active          bool       `json:"active"`
}

type EditPartnerDTO struct {
	FirstName       *string     `json:"first_name"`
	LastName        *string     `json:"last_name"`
	CompanyName     *string     `json:"company_name"`
	LegalName       *string     `json:"legal_name"`
	PartnerType     *string     `json:"partner_type"`
	Document        *string     `json:"document"`
	DocumentType    *string     `json:"document_type"`
	ShippingAddress *AddressDTO `json:"shipping"`
	BillingAddress  *AddressDTO `json:"billing"`
	PersonalContact *ContactDTO `json:"personal_contact"`
	BusinessContact *ContactDTO `json:"business_contact"`
	Active          *bool       `json:"active"`
	UpdatedBy       string      `json:"updated_by"`
}

func mapPartnerToPartnerDTO(partner domain.Partner) PartnerDTO {
	return PartnerDTO{
		PartnerID:       partner.PartnerID,
		FirstName:       partner.FirstName,
		LastName:        partner.LastName,
		CompanyName:     partner.CompanyName,
		LegalName:       partner.LegalName,
		PartnerType:     string(partner.PartnerType),
		Document:        partner.Document,
		DocumentType:    string(partner.DocumentType),
		ShippingAddress: mapAddressToAddressDTO(partner.ShippingAddress),
		BillingAddress:  mapAddressToAddressDTO(partner.BillingAddress),
		Region:          partner.Region,
		PersonalContact: mapContactToContactDTO(partner.PersonalContact),
		BusinessContact: mapContactToContactDTO(partner.BusinessContact),
		CreatedBy:       partner.CreatedBy,
		CreatedAt:       partner.CreatedAt,
		UpdatedBy:       partner.UpdatedBy,
		UpdatedAt:       partner.UpdatedAt,
		Active:          partner.Active,
	}
}

func mapCreatePartnerDTOToPartner(partnerDTO CreatePartnerDTO) (domain.Partner, error) {
	return domain.NewPartner(
		partnerDTO.FirstName,
		partnerDTO.LastName,
		partnerDTO.CompanyName,
		partnerDTO.LegalName,
		partnerDTO.Document,
		partnerDTO.DocumentType,
		partnerDTO.CreatedBy,
		mapContactDTOToContact(partnerDTO.PersonalContact),
		mapContactDTOToContact(partnerDTO.BusinessContact),
		mapAddressDTOToAddress(partnerDTO.ShippingAddress),
		mapAddressDTOToAddress(partnerDTO.BillingAddress),
	)
}

func mapPartnersToPartnerDTOs(partners []domain.Partner) []PartnerDTO {
	partnerDTOs := make([]PartnerDTO, 0, len(partners))
	for _, partner := range partners {
		partnerDTO := mapPartnerToPartnerDTO(partner)
		partnerDTOs = append(partnerDTOs, partnerDTO)
	}

	return partnerDTOs
}

func mapEditPartnerDTOToEditPartner(editPartnerDTO EditPartnerDTO) domain.EditPartner {
	var parsedPartnerType *domain.EntityType
	if editPartnerDTO.PartnerType != nil {
		partnerType := domain.EntityType(*editPartnerDTO.PartnerType)
		parsedPartnerType = &partnerType
	}

	var parsedDocumentType *domain.DocumentType
	if editPartnerDTO.DocumentType != nil {
		documentType := domain.DocumentType(*editPartnerDTO.DocumentType)
		parsedDocumentType = &documentType
	}

	var parsedShippingAddress *domain.Address
	if editPartnerDTO.ShippingAddress != nil {
		shippingAddress := mapAddressDTOToAddress(*editPartnerDTO.ShippingAddress)
		parsedShippingAddress = &shippingAddress
	}

	var parsedBillingAddress *domain.Address
	if editPartnerDTO.BillingAddress != nil {
		billingAddress := mapAddressDTOToAddress(*editPartnerDTO.BillingAddress)
		parsedBillingAddress = &billingAddress
	}

	var parsedPersonalContact *domain.Contact
	if editPartnerDTO.PersonalContact != nil {
		personalContact := mapContactDTOToContact(*editPartnerDTO.PersonalContact)
		parsedPersonalContact = &personalContact
	}

	var parsedBusinessContact *domain.Contact
	if editPartnerDTO.BusinessContact != nil {
		businessContact := mapContactDTOToContact(*editPartnerDTO.BusinessContact)
		parsedBusinessContact = &businessContact
	}

	return domain.EditPartner{
		FirstName:       editPartnerDTO.FirstName,
		LastName:        editPartnerDTO.LastName,
		CompanyName:     editPartnerDTO.CompanyName,
		LegalName:       editPartnerDTO.LegalName,
		PartnerType:     parsedPartnerType,
		Document:        editPartnerDTO.Document,
		DocumentType:    parsedDocumentType,
		ShippingAddress: parsedShippingAddress,
		BillingAddress:  parsedBillingAddress,
		PersonalContact: parsedPersonalContact,
		BusinessContact: parsedBusinessContact,
		Active:          editPartnerDTO.Active,
		UpdatedBy:       editPartnerDTO.UpdatedBy,
	}
}
