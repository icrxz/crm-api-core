package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateContractorDTO struct {
	CompanyName                   string                        `json:"company_name"`
	LegalName                     string                        `json:"legal_name"`
	Document                      string                        `json:"document"`
	BusinessContact               ContactDTO                    `json:"business_contact"`
	ContractorPlatformTemplateDTO ContractorPlatformTemplateDTO `json:"template"`
	CreatedBy                     string                        `json:"created_by"`
}

type ContractorDTO struct {
	ContractorID    string     `json:"contractor_id"`
	CompanyName     string     `json:"company_name"`
	LegalName       string     `json:"legal_name"`
	Document        string     `json:"document"`
	DocumentType    string     `json:"document_type"`
	BusinessContact ContactDTO `json:"business_contact"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       string     `json:"updated_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Active          bool       `json:"active"`
}

type UpdateContractorDTO struct {
	CompanyName     *string     `json:"company_name"`
	LegalName       *string     `json:"legal_name"`
	Document        *string     `json:"document"`
	DocumentType    *string     `json:"document_type"`
	BusinessContact *ContactDTO `json:"business_contact"`
	UpdatedBy       string      `json:"updated_by"`
}

func mapContractorToContractorDTO(contractor domain.Contractor) ContractorDTO {
	return ContractorDTO{
		ContractorID:    contractor.ContractorID,
		CompanyName:     contractor.CompanyName,
		LegalName:       contractor.LegalName,
		Document:        contractor.Document,
		DocumentType:    string(contractor.DocumentType),
		BusinessContact: mapContactToContactDTO(contractor.BusinessContact),
		CreatedBy:       contractor.CreatedBy,
		CreatedAt:       contractor.CreatedAt,
		UpdatedBy:       contractor.UpdatedBy,
		UpdatedAt:       contractor.UpdatedAt,
		Active:          contractor.Active,
	}
}

func mapCreateContractorDTOToContractor(contractorDTO CreateContractorDTO) (domain.Contractor, error) {
	contractorPlatformTemplate, err := mapContractorPlatformTemplateDTOToContractorPlatformTemplate(contractorDTO.ContractorPlatformTemplateDTO)
	if err != nil {
		return domain.Contractor{}, err
	}

	return domain.NewContractor(
		contractorDTO.LegalName,
		contractorDTO.CompanyName,
		contractorDTO.Document,
		contractorDTO.CreatedBy,
		mapContactDTOToContact(contractorDTO.BusinessContact),
		contractorPlatformTemplate,
	)
}

func mapContractorsToContractorDTOs(contractors []domain.Contractor) []ContractorDTO {
	contractorDTOs := make([]ContractorDTO, 0, len(contractors))
	for _, contractor := range contractors {
		contractorDTO := mapContractorToContractorDTO(contractor)
		contractorDTOs = append(contractorDTOs, contractorDTO)
	}

	return contractorDTOs
}

func mapUpdateContractorDTOToUpdateContractor(updateContractorDTO UpdateContractorDTO) domain.UpdateContractor {
	var parsedDocumentType *domain.DocumentType
	if updateContractorDTO.DocumentType != nil {
		documentType := domain.DocumentType(*updateContractorDTO.DocumentType)
		parsedDocumentType = &documentType
	}

	var parsedBusinessContact *domain.Contact
	if updateContractorDTO.BusinessContact != nil {
		businessContact := mapContactDTOToContact(*updateContractorDTO.BusinessContact)
		parsedBusinessContact = &businessContact
	}

	return domain.UpdateContractor{
		CompanyName:     updateContractorDTO.CompanyName,
		LegalName:       updateContractorDTO.LegalName,
		Document:        updateContractorDTO.Document,
		DocumentType:    parsedDocumentType,
		BusinessContact: parsedBusinessContact,
		UpdatedBy:       updateContractorDTO.UpdatedBy,
	}
}
