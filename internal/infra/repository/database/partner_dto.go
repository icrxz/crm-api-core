package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type PartnerDTO struct {
	PartnerID       string    `db:"partner_id"`
	FirstName       string    `db:"first_name"`
	LastName        string    `db:"last_name"`
	CompanyName     string    `db:"company_name"`
	LegalName       string    `db:"legal_name"`
	PartnerType     string    `db:"partner_type"`
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
	Region          int       `db:"region"`
	CreatedBy       string    `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedBy       string    `db:"updated_by"`
	UpdatedAt       time.Time `db:"updated_at"`
	Active          bool      `db:"active"`
	Description     *string   `db:"description"`
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
		ShippingAddress: partner.ShippingAddress.Address,
		ShippingCity:    partner.ShippingAddress.City,
		ShippingState:   partner.ShippingAddress.State,
		ShippingZipCode: partner.ShippingAddress.ZipCode,
		ShippingCountry: partner.ShippingAddress.Country,
		BillingAddress:  partner.BillingAddress.Address,
		BillingCity:     partner.BillingAddress.City,
		BillingState:    partner.BillingAddress.State,
		BillingZipCode:  partner.BillingAddress.ZipCode,
		BillingCountry:  partner.BillingAddress.Country,
		PersonalPhone:   partner.PersonalContact.PhoneNumber,
		BusinessPhone:   partner.BusinessContact.PhoneNumber,
		PersonalEmail:   partner.PersonalContact.Email,
		BusinessEmail:   partner.BusinessContact.Email,
		CreatedBy:       partner.CreatedBy,
		CreatedAt:       partner.CreatedAt,
		UpdatedBy:       partner.UpdatedBy,
		UpdatedAt:       partner.UpdatedAt,
		Active:          partner.Active,
		Region:          partner.GetRegion(),
		Description:     &partner.Description,
	}
}

func mapPartnerDTOToPartner(partnerDTO PartnerDTO) domain.Partner {
	var descriptionString string
	if partnerDTO.Description != nil {
		descriptionString = *partnerDTO.Description
	}

	return domain.Partner{
		PartnerID:    partnerDTO.PartnerID,
		FirstName:    partnerDTO.FirstName,
		LastName:     partnerDTO.LastName,
		CompanyName:  partnerDTO.CompanyName,
		LegalName:    partnerDTO.LegalName,
		PartnerType:  domain.EntityType(partnerDTO.PartnerType),
		Document:     partnerDTO.Document,
		DocumentType: domain.DocumentType(partnerDTO.DocumentType),
		ShippingAddress: domain.Address{
			Address: partnerDTO.ShippingAddress,
			City:    partnerDTO.ShippingCity,
			State:   partnerDTO.ShippingState,
			ZipCode: partnerDTO.ShippingZipCode,
			Country: partnerDTO.ShippingCountry,
		},
		BillingAddress: domain.Address{
			Address: partnerDTO.BillingAddress,
			City:    partnerDTO.BillingCity,
			State:   partnerDTO.BillingState,
			ZipCode: partnerDTO.BillingZipCode,
			Country: partnerDTO.BillingCountry,
		},
		PersonalContact: domain.Contact{
			PhoneNumber: partnerDTO.PersonalPhone,
			Email:       partnerDTO.PersonalEmail,
		},
		BusinessContact: domain.Contact{
			PhoneNumber: partnerDTO.BusinessPhone,
			Email:       partnerDTO.BusinessEmail,
		},
		CreatedBy:   partnerDTO.CreatedBy,
		CreatedAt:   partnerDTO.CreatedAt,
		UpdatedBy:   partnerDTO.UpdatedBy,
		UpdatedAt:   partnerDTO.UpdatedAt,
		Active:      partnerDTO.Active,
		Description: descriptionString,
	}
}

func mapPartnerDTOsToPartners(partnerDTOs []PartnerDTO) []domain.Partner {
	partners := make([]domain.Partner, 0, len(partnerDTOs))
	for _, partnerDTO := range partnerDTOs {
		partner := mapPartnerDTOToPartner(partnerDTO)
		partners = append(partners, partner)
	}

	return partners
}

func mapPartnersToPartnerDTOs(partners []domain.Partner) []PartnerDTO {
	partnerDTOs := make([]PartnerDTO, 0, len(partners))
	for _, partner := range partners {
		partnerDTO := mapPartnerToPartnerDTO(partner)
		partnerDTOs = append(partnerDTOs, partnerDTO)
	}

	return partnerDTOs
}
