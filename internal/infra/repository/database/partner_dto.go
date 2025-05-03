package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type PartnerDTO struct {
	PartnerID              string    `db:"partner_id"`
	FirstName              string    `db:"first_name"`
	LastName               string    `db:"last_name"`
	CompanyName            string    `db:"company_name"`
	LegalName              string    `db:"legal_name"`
	PartnerType            string    `db:"partner_type"`
	Document               string    `db:"document"`
	DocumentType           string    `db:"document_type"`
	ShippingAddress        string    `db:"shipping_address"`
	ShippingCity           string    `db:"shipping_city"`
	ShippingState          string    `db:"shipping_state"`
	ShippingZipCode        string    `db:"shipping_zip_code"`
	ShippingCountry        string    `db:"shipping_country"`
	BillingAddress         string    `db:"billing_address"`
	BillingCity            string    `db:"billing_city"`
	BillingState           string    `db:"billing_state"`
	BillingZipCode         string    `db:"billing_zip_code"`
	BillingCountry         string    `db:"billing_country"`
	PersonalPhone          string    `db:"personal_phone"`
	BusinessPhone          string    `db:"business_phone"`
	PersonalEmail          string    `db:"personal_email"`
	BusinessEmail          string    `db:"business_email"`
	Region                 *int      `db:"region"`
	CreatedBy              string    `db:"created_by"`
	CreatedAt              time.Time `db:"created_at"`
	UpdatedBy              string    `db:"updated_by"`
	UpdatedAt              time.Time `db:"updated_at"`
	Active                 bool      `db:"active"`
	Description            *string   `db:"description"`
	PaymentKey             *string   `db:"payment_key"`
	PaymentKeyOption       *string   `db:"payment_key_option"`
	PaymentType            *string   `db:"payment_type"`
	PaymentOwner           *string   `db:"payment_owner"`
	PaymentIsSameFromOwner *bool     `db:"payment_is_same_from_owner"`
}

type PartnerOptionalDTO struct {
	PartnerID              *string    `db:"partner_id"`
	FirstName              *string    `db:"first_name"`
	LastName               *string    `db:"last_name"`
	CompanyName            *string    `db:"company_name"`
	LegalName              *string    `db:"legal_name"`
	PartnerType            *string    `db:"partner_type"`
	Document               *string    `db:"document"`
	DocumentType           *string    `db:"document_type"`
	ShippingAddress        *string    `db:"shipping_address"`
	ShippingCity           *string    `db:"shipping_city"`
	ShippingState          *string    `db:"shipping_state"`
	ShippingZipCode        *string    `db:"shipping_zip_code"`
	ShippingCountry        *string    `db:"shipping_country"`
	BillingAddress         *string    `db:"billing_address"`
	BillingCity            *string    `db:"billing_city"`
	BillingState           *string    `db:"billing_state"`
	BillingZipCode         *string    `db:"billing_zip_code"`
	BillingCountry         *string    `db:"billing_country"`
	PersonalPhone          *string    `db:"personal_phone"`
	BusinessPhone          *string    `db:"business_phone"`
	PersonalEmail          *string    `db:"personal_email"`
	BusinessEmail          *string    `db:"business_email"`
	Region                 *int       `db:"region"`
	CreatedBy              *string    `db:"created_by"`
	CreatedAt              *time.Time `db:"created_at"`
	UpdatedBy              *string    `db:"updated_by"`
	UpdatedAt              *time.Time `db:"updated_at"`
	Active                 *bool      `db:"active"`
	Description            *string    `db:"description"`
	PaymentKey             *string    `db:"payment_key"`
	PaymentKeyOption       *string    `db:"payment_key_option"`
	PaymentType            *string    `db:"payment_type"`
	PaymentOwner           *string    `db:"payment_owner"`
	PaymentIsSameFromOwner *bool      `db:"payment_is_same_from_owner"`
}

func mapPartnerToPartnerDTO(partner domain.Partner) PartnerDTO {
	paymentTypeString := string(partner.Billing.Type)

	return PartnerDTO{
		PartnerID:              partner.PartnerID,
		FirstName:              partner.FirstName,
		LastName:               partner.LastName,
		CompanyName:            partner.CompanyName,
		LegalName:              partner.LegalName,
		PartnerType:            partner.PartnerType,
		Document:               partner.Document,
		DocumentType:           string(partner.DocumentType),
		ShippingAddress:        partner.ShippingAddress.Address,
		ShippingCity:           partner.ShippingAddress.City,
		ShippingState:          partner.ShippingAddress.State,
		ShippingZipCode:        partner.ShippingAddress.ZipCode,
		ShippingCountry:        partner.ShippingAddress.Country,
		BillingAddress:         partner.BillingAddress.Address,
		BillingCity:            partner.BillingAddress.City,
		BillingState:           partner.BillingAddress.State,
		BillingZipCode:         partner.BillingAddress.ZipCode,
		BillingCountry:         partner.BillingAddress.Country,
		PersonalPhone:          partner.PersonalContact.PhoneNumber,
		BusinessPhone:          partner.BusinessContact.PhoneNumber,
		PersonalEmail:          partner.PersonalContact.Email,
		BusinessEmail:          partner.BusinessContact.Email,
		CreatedBy:              partner.CreatedBy,
		CreatedAt:              partner.CreatedAt,
		UpdatedBy:              partner.UpdatedBy,
		UpdatedAt:              partner.UpdatedAt,
		Active:                 partner.Active,
		Description:            &partner.Description,
		PaymentKey:             &partner.Billing.Key,
		PaymentKeyOption:       &partner.Billing.Option,
		PaymentType:            &paymentTypeString,
		PaymentOwner:           &partner.Billing.Name,
		PaymentIsSameFromOwner: &partner.Billing.IsSameFromOwner,
	}
}

func mapPartnerDTOToPartner(partnerDTO PartnerDTO) domain.Partner {
	var descriptionString string
	if partnerDTO.Description != nil {
		descriptionString = *partnerDTO.Description
	}

	var paymentKeyString string
	if partnerDTO.PaymentKey != nil {
		paymentKeyString = *partnerDTO.PaymentKey
	}

	var paymentOptionString string
	if partnerDTO.PaymentKeyOption != nil {
		paymentOptionString = *partnerDTO.PaymentKeyOption
	}

	var paymentTypeString string
	if partnerDTO.PaymentType != nil {
		paymentTypeString = *partnerDTO.PaymentType
	}

	var paymentOwnerString string
	if partnerDTO.PaymentOwner != nil {
		paymentOwnerString = *partnerDTO.PaymentOwner
	}

	var paymentIsSameFromOwner bool
	if partnerDTO.PaymentIsSameFromOwner != nil {
		paymentIsSameFromOwner = *partnerDTO.PaymentIsSameFromOwner
	}

	return domain.Partner{
		PartnerID:    partnerDTO.PartnerID,
		FirstName:    partnerDTO.FirstName,
		LastName:     partnerDTO.LastName,
		CompanyName:  partnerDTO.CompanyName,
		LegalName:    partnerDTO.LegalName,
		PartnerType:  partnerDTO.PartnerType,
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
		Billing: domain.Billing{
			Key:             paymentKeyString,
			Option:          paymentOptionString,
			Type:            domain.BillingType(paymentTypeString),
			Name:            paymentOwnerString,
			IsSameFromOwner: paymentIsSameFromOwner,
		},
	}
}

func mapPartnerOptionalDTOToPartner(partnerDTO PartnerOptionalDTO) domain.Partner {
	if partnerDTO.PartnerID == nil {
		return domain.Partner{}
	}

	var descriptionString string
	if partnerDTO.Description != nil {
		descriptionString = *partnerDTO.Description
	}

	var paymentKeyString string
	if partnerDTO.PaymentKey != nil {
		paymentKeyString = *partnerDTO.PaymentKey
	}

	var paymentOptionString string
	if partnerDTO.PaymentKeyOption != nil {
		paymentOptionString = *partnerDTO.PaymentKeyOption
	}

	var paymentTypeString string
	if partnerDTO.PaymentType != nil {
		paymentTypeString = *partnerDTO.PaymentType
	}

	var paymentOwnerString string
	if partnerDTO.PaymentOwner != nil {
		paymentOwnerString = *partnerDTO.PaymentOwner
	}

	var paymentIsSameFromOwner bool
	if partnerDTO.PaymentIsSameFromOwner != nil {
		paymentIsSameFromOwner = *partnerDTO.PaymentIsSameFromOwner
	}

	return domain.Partner{
		PartnerID:    *partnerDTO.PartnerID,
		FirstName:    *partnerDTO.FirstName,
		LastName:     *partnerDTO.LastName,
		CompanyName:  *partnerDTO.CompanyName,
		LegalName:    *partnerDTO.LegalName,
		PartnerType:  *partnerDTO.PartnerType,
		Document:     *partnerDTO.Document,
		DocumentType: domain.DocumentType(*partnerDTO.DocumentType),
		ShippingAddress: domain.Address{
			Address: *partnerDTO.ShippingAddress,
			City:    *partnerDTO.ShippingCity,
			State:   *partnerDTO.ShippingState,
			ZipCode: *partnerDTO.ShippingZipCode,
			Country: *partnerDTO.ShippingCountry,
		},
		BillingAddress: domain.Address{
			Address: *partnerDTO.BillingAddress,
			City:    *partnerDTO.BillingCity,
			State:   *partnerDTO.BillingState,
			ZipCode: *partnerDTO.BillingZipCode,
			Country: *partnerDTO.BillingCountry,
		},
		PersonalContact: domain.Contact{
			PhoneNumber: *partnerDTO.PersonalPhone,
			Email:       *partnerDTO.PersonalEmail,
		},
		BusinessContact: domain.Contact{
			PhoneNumber: *partnerDTO.BusinessPhone,
			Email:       *partnerDTO.BusinessEmail,
		},
		CreatedBy:   *partnerDTO.CreatedBy,
		CreatedAt:   *partnerDTO.CreatedAt,
		UpdatedBy:   *partnerDTO.UpdatedBy,
		UpdatedAt:   *partnerDTO.UpdatedAt,
		Active:      *partnerDTO.Active,
		Description: descriptionString,
		Billing: domain.Billing{
			Key:             paymentKeyString,
			Option:          paymentOptionString,
			Type:            domain.BillingType(paymentTypeString),
			Name:            paymentOwnerString,
			IsSameFromOwner: paymentIsSameFromOwner,
		},
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
