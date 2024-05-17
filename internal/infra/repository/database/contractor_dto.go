package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type ContractorDTO struct {
	ContractorID  string    `db:"contractor_id"`
	CompanyName   string    `db:"company_name"`
	LegalName     string    `db:"legal_name"`
	Document      string    `db:"document"`
	DocumentType  string    `db:"document_type"`
	BusinessPhone string    `db:"business_phone"`
	BusinessEmail string    `db:"business_email"`
	CreatedBy     string    `db:"created_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedBy     string    `db:"updated_by"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func mapContractorToContractorDTO(contractor domain.Contractor) ContractorDTO {
	return ContractorDTO{
		ContractorID:  contractor.ContractorID,
		CompanyName:   contractor.CompanyName,
		LegalName:     contractor.LegalName,
		Document:      contractor.Document,
		DocumentType:  string(contractor.DocumentType),
		BusinessPhone: contractor.BusinessContact.PhoneNumber,
		BusinessEmail: contractor.BusinessContact.Email,
		CreatedBy:     contractor.CreatedBy,
		CreatedAt:     contractor.CreatedAt,
		UpdatedBy:     contractor.UpdatedBy,
		UpdatedAt:     contractor.UpdatedAt,
	}
}

func mapContractorDTOToContractor(contractorDTO ContractorDTO) domain.Contractor {
	return domain.Contractor{
		ContractorID: contractorDTO.ContractorID,
		CompanyName:  contractorDTO.CompanyName,
		LegalName:    contractorDTO.LegalName,
		Document:     contractorDTO.Document,
		DocumentType: domain.DocumentType(contractorDTO.DocumentType),
		BusinessContact: domain.Contact{
			PhoneNumber: contractorDTO.BusinessPhone,
			Email:       contractorDTO.BusinessEmail,
		},
		CreatedBy: contractorDTO.CreatedBy,
		CreatedAt: contractorDTO.CreatedAt,
		UpdatedBy: contractorDTO.UpdatedBy,
		UpdatedAt: contractorDTO.UpdatedAt,
	}
}

func mapContractorDTOsToContractors(contractorDTOs []ContractorDTO) []domain.Contractor {
	contractors := make([]domain.Contractor, 0, len(contractorDTOs))
	for _, contractorDTO := range contractorDTOs {
		contractor := mapContractorDTOToContractor(contractorDTO)
		contractors = append(contractors, contractor)
	}

	return contractors
}
