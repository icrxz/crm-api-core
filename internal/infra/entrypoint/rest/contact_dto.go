package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type ContactDTO struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type UpdateContactDTO struct {
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
}

func mapContactDTOToContact(contactDTO ContactDTO) domain.Contact {
	return domain.Contact{
		PhoneNumber: contactDTO.PhoneNumber,
		Email:       contactDTO.Email,
	}
}

func mapContactToContactDTO(contact domain.Contact) ContactDTO {
	return ContactDTO{
		PhoneNumber: contact.PhoneNumber,
		Email:       contact.Email,
	}
}

func mapUpdateContactDTOToUpdateContact(updateContactDTO UpdateContactDTO) domain.UpdateContact {
	return domain.UpdateContact{
		PhoneNumber: updateContactDTO.PhoneNumber,
		Email:       updateContactDTO.Email,
	}
}
