package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CaseDTO struct {
	CaseID            string     `db:"case_id"`
	ContractorID      string     `db:"contractor_id"`
	CustomerID        *string    `db:"customer_id"`
	PartnerID         *string    `db:"partner_id"`
	OwnerID           *string    `db:"owner_id"`
	OriginChannel     string     `db:"origin"`
	Type              string     `db:"type"`
	Subject           string     `db:"subject"`
	Priority          string     `db:"priority"`
	Status            string     `db:"status"`
	DueDate           time.Time  `db:"due_date"`
	CreatedBy         string     `db:"created_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedBy         string     `db:"updated_by"`
	UpdatedAt         time.Time  `db:"updated_at"`
	ExternalReference string     `db:"external_reference"`
	ProductID         *string    `db:"product_id"`
	Region            int        `db:"region"`
	ClosedAt          *time.Time `db:"closed_at"`
	TargetDate        *time.Time `db:"target_date"`
}

type CaseFullDTO struct {
	CaseID            string              `db:"case_id"`
	OwnerID           *string             `db:"owner_id"`
	OriginChannel     string              `db:"origin"`
	Type              string              `db:"type"`
	Subject           string              `db:"subject"`
	Priority          string              `db:"priority"`
	Status            string              `db:"status"`
	DueDate           time.Time           `db:"due_date"`
	CreatedBy         string              `db:"created_by"`
	CreatedAt         time.Time           `db:"created_at"`
	UpdatedBy         string              `db:"updated_by"`
	UpdatedAt         time.Time           `db:"updated_at"`
	ExternalReference string              `db:"external_reference"`
	Region            int                 `db:"region"`
	ClosedAt          *time.Time          `db:"closed_at"`
	TargetDate        *time.Time          `db:"target_date"`
	Contractor        ContractorDTO       `db:"contractor"`
	Customer          CustomerOptionalDTO `db:"customers"`
	Partner           PartnerOptionalDTO  `db:"partners"`
	Product           ProductDTO          `db:"products"`
	Transactions      []TransactionDTO    `db:"transactions"`
}

func mapCaseToCaseDTO(crmCase domain.Case) CaseDTO {
	var partnerID *string
	if crmCase.PartnerID != "" {
		partnerID = &crmCase.PartnerID
	}

	var ownerID *string
	if crmCase.OwnerID != "" {
		ownerID = &crmCase.OwnerID
	}

	var customerID *string
	if crmCase.CustomerID != "" {
		customerID = &crmCase.CustomerID
	}

	var productID *string
	if crmCase.ProductID != "" {
		productID = &crmCase.ProductID
	}

	return CaseDTO{
		CaseID:            crmCase.CaseID,
		ContractorID:      crmCase.ContractorID,
		CustomerID:        customerID,
		PartnerID:         partnerID,
		OwnerID:           ownerID,
		OriginChannel:     crmCase.OriginChannel,
		Type:              crmCase.Type,
		Subject:           crmCase.Subject,
		Priority:          string(crmCase.Priority),
		Status:            string(crmCase.Status),
		DueDate:           crmCase.DueDate,
		CreatedBy:         crmCase.CreatedBy,
		CreatedAt:         crmCase.CreatedAt,
		UpdatedBy:         crmCase.UpdatedBy,
		UpdatedAt:         crmCase.UpdatedAt,
		ExternalReference: crmCase.ExternalReference,
		Region:            crmCase.Region,
		ProductID:         productID,
		TargetDate:        crmCase.TargetDate,
		ClosedAt:          crmCase.ClosedAt,
	}
}

func mapCaseDTOToCase(crmCaseDTO CaseDTO) domain.Case {
	var partnerID string
	if crmCaseDTO.PartnerID != nil {
		partnerID = *crmCaseDTO.PartnerID
	}

	var ownerID string
	if crmCaseDTO.OwnerID != nil {
		ownerID = *crmCaseDTO.OwnerID
	}

	var customerID string
	if crmCaseDTO.CustomerID != nil {
		customerID = *crmCaseDTO.CustomerID
	}

	var productID string
	if crmCaseDTO.ProductID != nil {
		productID = *crmCaseDTO.ProductID
	}

	return domain.Case{
		CaseID:            crmCaseDTO.CaseID,
		ContractorID:      crmCaseDTO.ContractorID,
		CustomerID:        customerID,
		PartnerID:         partnerID,
		OwnerID:           ownerID,
		OriginChannel:     crmCaseDTO.OriginChannel,
		Type:              crmCaseDTO.Type,
		Subject:           crmCaseDTO.Subject,
		Priority:          domain.CasePriority(crmCaseDTO.Priority),
		Status:            domain.CaseStatus(crmCaseDTO.Status),
		DueDate:           crmCaseDTO.DueDate,
		CreatedBy:         crmCaseDTO.CreatedBy,
		CreatedAt:         crmCaseDTO.CreatedAt,
		UpdatedBy:         crmCaseDTO.UpdatedBy,
		UpdatedAt:         crmCaseDTO.UpdatedAt,
		ExternalReference: crmCaseDTO.ExternalReference,
		Region:            crmCaseDTO.Region,
		ProductID:         productID,
		ClosedAt:          crmCaseDTO.ClosedAt,
		TargetDate:        crmCaseDTO.TargetDate,
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

func mapCasesToCaseDTOs(cases []domain.Case) []CaseDTO {
	crmCaseDTOs := make([]CaseDTO, 0, len(cases))
	for _, crmCase := range cases {
		caseDTO := mapCaseToCaseDTO(crmCase)
		crmCaseDTOs = append(crmCaseDTOs, caseDTO)
	}

	return crmCaseDTOs
}

func mapCaseFullDTOToCaseFull(crmCaseFullDTO CaseFullDTO) domain.CaseFull {
	var ownerID string
	if crmCaseFullDTO.OwnerID != nil {
		ownerID = *crmCaseFullDTO.OwnerID
	}

	return domain.CaseFull{
		CaseID:            crmCaseFullDTO.CaseID,
		OwnerID:           ownerID,
		OriginChannel:     crmCaseFullDTO.OriginChannel,
		Type:              crmCaseFullDTO.Type,
		Subject:           crmCaseFullDTO.Subject,
		Priority:          domain.CasePriority(crmCaseFullDTO.Priority),
		Status:            domain.CaseStatus(crmCaseFullDTO.Status),
		DueDate:           crmCaseFullDTO.DueDate,
		CreatedBy:         crmCaseFullDTO.CreatedBy,
		CreatedAt:         crmCaseFullDTO.CreatedAt,
		UpdatedBy:         crmCaseFullDTO.UpdatedBy,
		UpdatedAt:         crmCaseFullDTO.UpdatedAt,
		ExternalReference: crmCaseFullDTO.ExternalReference,
		Region:            crmCaseFullDTO.Region,
		ClosedAt:          crmCaseFullDTO.ClosedAt,
		TargetDate:        crmCaseFullDTO.TargetDate,
		Contractor:        mapContractorDTOToContractor(crmCaseFullDTO.Contractor),
		Customer:          mapCustomerOptionalDTOToCustomer(crmCaseFullDTO.Customer),
		Partner:           mapPartnerOptionalDTOToPartner(crmCaseFullDTO.Partner),
		Product:           mapProductDTOToProduct(crmCaseFullDTO.Product),
		Transactions:      mapTransactionDTOsToTransactions(crmCaseFullDTO.Transactions),
	}
}

func mapCaseFullDTOsToCasesFull(crmCaseDTOs []CaseFullDTO) []domain.CaseFull {
	crmCases := make([]domain.CaseFull, 0, len(crmCaseDTOs))
	for _, crmCaseDTO := range crmCaseDTOs {
		crmCase := mapCaseFullDTOToCaseFull(crmCaseDTO)
		crmCases = append(crmCases, crmCase)
	}

	return crmCases
}
