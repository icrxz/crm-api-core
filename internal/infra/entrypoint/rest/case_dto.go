package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateCaseDTO struct {
	ContractorID       string    `json:"contractor_id" validate:"required"`
	CustomerID         string    `json:"customer_id" validate:"required"`
	OriginChannel      string    `json:"origin_channel" validate:"required"`
	CaseType           string    `json:"case_type" validate:"required"`
	Subject            string    `json:"subject" validate:"required"`
	DueDate            time.Time `json:"due_date" validate:"required"`
	CreatedBy          string    `json:"created_by" validate:"required"`
	ExternalReference  string    `json:"external_reference"`
	ProductName        string    `json:"product_name"`
	Brand              string    `json:"brand" validate:"required"`
	Model              string    `json:"model" validate:"required"`
	ProductDescription string    `json:"product_description"`
	Value              float64   `json:"value"`
	SerialNumber       string    `json:"serial_number"`
}

type CaseDTO struct {
	CaseID            string              `json:"case_id"`
	ContractorID      string              `json:"contractor_id"`
	CustomerID        string              `json:"customer_id"`
	PartnerID         string              `json:"partner_id"`
	OwnerID           string              `json:"owner_id"`
	OriginChannel     string              `json:"origin_channel"`
	Type              string              `json:"type"`
	Subject           string              `json:"subject"`
	Priority          domain.CasePriority `json:"priority"`
	Status            domain.CaseStatus   `json:"status"`
	DueDate           time.Time           `json:"due_date"`
	CreatedBy         string              `json:"created_by"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedBy         string              `json:"updated_by"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Region            int                 `json:"region"`
	ExternalReference string              `json:"external_reference"`
	ProductID         string              `json:"product_id"`
	ClosedAt          *time.Time          `json:"closed_at"`
	TargetDate        *time.Time          `json:"target_date"`
}

type CaseFullDTO struct {
	CaseID            string              `json:"case_id"`
	Contractor        ContractorDTO       `json:"contractor"`
	Customer          CustomerDTO         `json:"customer"`
	Partner           PartnerDTO          `json:"partner"`
	Product           ProductDTO          `json:"product"`
	Comments          []CommentDTO        `json:"comments"`
	Transactions      []TransactionDTO    `json:"transactions"`
	OwnerID           string              `json:"owner_id"`
	OriginChannel     string              `json:"origin_channel"`
	Type              string              `json:"type"`
	Subject           string              `json:"subject"`
	Priority          domain.CasePriority `json:"priority"`
	Status            domain.CaseStatus   `json:"status"`
	DueDate           time.Time           `json:"due_date"`
	CreatedBy         string              `json:"created_by"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedBy         string              `json:"updated_by"`
	UpdatedAt         time.Time           `json:"updated_at"`
	Region            int                 `json:"region"`
	ExternalReference string              `json:"external_reference"`
	ClosedAt          *time.Time          `json:"closed_at"`
	TargetDate        *time.Time          `json:"target_date"`
}

type UpdateCaseDTO struct {
	TargetDate *time.Time `json:"target_date"`
	Status     *string    `json:"status"`
	PartnerID  *string    `json:"partner_id"`
	OwnerID    *string    `json:"owner_id"`
	CustomerID *string    `json:"customer_id"`
	ProductID  *string    `json:"product_id"`
	Subject    *string    `json:"subject"`
	UpdatedBy  string     `json:"updated_by" validate:"required"`
}

func mapCreateCaseDTOToCreateCase(createCaseDTO CreateCaseDTO) (domain.CreateCase, error) {
	crmCase, err := domain.NewCase(
		createCaseDTO.ContractorID,
		createCaseDTO.CustomerID,
		createCaseDTO.OriginChannel,
		createCaseDTO.CaseType,
		createCaseDTO.Subject,
		createCaseDTO.DueDate,
		createCaseDTO.CreatedBy,
		createCaseDTO.ExternalReference,
	)
	if err != nil {
		return domain.CreateCase{}, err
	}

	product, err := domain.NewProduct(
		createCaseDTO.ProductName,
		createCaseDTO.ProductDescription,
		createCaseDTO.Value,
		createCaseDTO.Brand,
		createCaseDTO.Model,
		createCaseDTO.SerialNumber,
		createCaseDTO.CreatedBy,
	)
	if err != nil {
		return domain.CreateCase{}, err
	}

	return domain.CreateCase{
		Case:    crmCase,
		Product: product,
	}, nil
}

func mapCaseToCaseDTO(crmCase domain.Case) CaseDTO {
	return CaseDTO{
		CaseID:            crmCase.CaseID,
		ContractorID:      crmCase.ContractorID,
		CustomerID:        crmCase.CustomerID,
		PartnerID:         crmCase.PartnerID,
		OwnerID:           crmCase.OwnerID,
		OriginChannel:     crmCase.OriginChannel,
		Type:              crmCase.Type,
		Subject:           crmCase.Subject,
		Priority:          crmCase.Priority,
		Status:            crmCase.Status,
		DueDate:           crmCase.DueDate,
		CreatedBy:         crmCase.CreatedBy,
		CreatedAt:         crmCase.CreatedAt,
		UpdatedBy:         crmCase.UpdatedBy,
		UpdatedAt:         crmCase.UpdatedAt,
		Region:            crmCase.Region,
		ExternalReference: crmCase.ExternalReference,
		ProductID:         crmCase.ProductID,
		ClosedAt:          crmCase.ClosedAt,
		TargetDate:        crmCase.TargetDate,
	}
}

func mapCasesToCaseDTOs(crmCases []domain.Case) []CaseDTO {
	crmCasesDTO := make([]CaseDTO, len(crmCases))
	for i, crmCase := range crmCases {
		crmCasesDTO[i] = mapCaseToCaseDTO(crmCase)
	}
	return crmCasesDTO
}

func mapUpdateCaseDTOToUpdateCase(dto UpdateCaseDTO) domain.CaseUpdate {
	var status *domain.CaseStatus
	if dto.Status != nil {
		newStatus := domain.CaseStatus(*dto.Status)
		status = &newStatus
	}

	return domain.CaseUpdate{
		TargetDate: dto.TargetDate,
		UpdatedBy:  dto.UpdatedBy,
		Status:     status,
		PartnerID:  dto.PartnerID,
		OwnerID:    dto.OwnerID,
		CustomerID: dto.CustomerID,
		Subject:    dto.Subject,
		ProductID:  dto.ProductID,
	}
}

func mapCaseFullToCaseFullDTO(caseFull domain.CaseFull) CaseFullDTO {
	contractor := mapContractorToContractorDTO(caseFull.Contractor)
	customer := mapCustomerToCustomerDTO(caseFull.Customer)
	partner := mapPartnerToPartnerDTO(caseFull.Partner)
	product := mapProductToProductDTO(caseFull.Product)

	comments := mapCommentsToCommentDTOs(caseFull.Comments)
	transactions := mapTransactionsToTransactionsDTO(caseFull.Transactions)

	return CaseFullDTO{
		CaseID:            caseFull.CaseID,
		Contractor:        contractor,
		Customer:          customer,
		Partner:           partner,
		Product:           product,
		Comments:          comments,
		Transactions:      transactions,
		OwnerID:           caseFull.OwnerID,
		OriginChannel:     caseFull.OriginChannel,
		Type:              caseFull.Type,
		Subject:           caseFull.Subject,
		Priority:          caseFull.Priority,
		Status:            caseFull.Status,
		DueDate:           caseFull.DueDate,
		CreatedBy:         caseFull.CreatedBy,
		CreatedAt:         caseFull.CreatedAt,
		UpdatedBy:         caseFull.UpdatedBy,
		UpdatedAt:         caseFull.UpdatedAt,
		Region:            caseFull.Region,
		ExternalReference: caseFull.ExternalReference,
		ClosedAt:          caseFull.ClosedAt,
		TargetDate:        caseFull.TargetDate,
	}
}

func mapCasesFullToCasesFullDTOs(crmCases []domain.CaseFull) []CaseFullDTO {
	crmCasesDTO := make([]CaseFullDTO, len(crmCases))
	for i, crmCase := range crmCases {
		crmCasesDTO[i] = mapCaseFullToCaseFullDTO(crmCase)
	}
	return crmCasesDTO
}
