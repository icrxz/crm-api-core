package database

import (
	"github.com/icrxz/crm-api-core/internal/domain"
	"time"
)

type CaseDTO struct {
	CaseID        string    `db:"case_id"`
	ContractorID  string    `db:"contractor_id"`
	CustomerID    string    `db:"customer_id"`
	PartnerID     string    `db:"partner_id"`
	OwnerID       string    `db:"owner_id"`
	OriginChannel string    `db:"origin"`
	Type          string    `db:"type"`
	Subject       string    `db:"subject"`
	Priority      string    `db:"priority"`
	Status        string    `db:"status"`
	DueDate       time.Time `db:"due_date"`
	CreatedBy     string    `db:"created_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedBy     string    `db:"updated_by"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func mapCaseToCaseDTO(crmCase domain.Case) CaseDTO {
	return CaseDTO{
		CaseID:        crmCase.CaseID,
		ContractorID:  crmCase.ContractorID,
		CustomerID:    crmCase.CustomerID,
		PartnerID:     crmCase.PartnerID,
		OwnerID:       crmCase.OwnerID,
		OriginChannel: crmCase.OriginChannel,
		Type:          crmCase.Type,
		Subject:       crmCase.Subject,
		Priority:      string(crmCase.Priority),
		Status:        string(crmCase.Status),
		DueDate:       crmCase.DueDate,
		CreatedBy:     crmCase.CreatedBy,
		CreatedAt:     crmCase.CreatedAt,
		UpdatedBy:     crmCase.UpdatedBy,
		UpdatedAt:     crmCase.UpdatedAt,
	}
}

func mapCaseDTOToCase(crmCaseDTO CaseDTO) domain.Case {
	return domain.Case{
		CaseID:        crmCaseDTO.CaseID,
		ContractorID:  crmCaseDTO.ContractorID,
		CustomerID:    crmCaseDTO.CustomerID,
		PartnerID:     crmCaseDTO.PartnerID,
		OwnerID:       crmCaseDTO.OwnerID,
		OriginChannel: crmCaseDTO.OriginChannel,
		Type:          crmCaseDTO.Type,
		Subject:       crmCaseDTO.Subject,
		Priority:      domain.CasePriority(crmCaseDTO.Priority),
		Status:        domain.CaseStatus(crmCaseDTO.Status),
		DueDate:       crmCaseDTO.DueDate,
		CreatedBy:     crmCaseDTO.CreatedBy,
		CreatedAt:     crmCaseDTO.CreatedAt,
		UpdatedBy:     crmCaseDTO.UpdatedBy,
		UpdatedAt:     crmCaseDTO.UpdatedAt,
	}
}

func mapCaseDTOsToCases(crmCaseDTOs []CaseDTO) []domain.Case {
	crmCases := make([]domain.Case, 0, len(crmCaseDTOs))
	for _, crmCaseDTO := range crmCaseDTOs {
		crmCase := mapCaseDTOToCase(crmCaseDTO)
		crmCases = append(crmCases, crmCase)
	}

	return crmCases
}
