package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CaseRepository interface {
	Create(ctx context.Context, crmCase Case) (string, error)
	GetByID(ctx context.Context, caseID string) (*Case, error)
	Search(ctx context.Context, filters CaseFilters) (PagingResult[Case], error)
	Update(ctx context.Context, crmCase Case) error
	CreateBatch(ctx context.Context, cases []Case) ([]string, error)
}

type CreateCase struct {
	Case    Case
	Product Product
}

type Case struct {
	CaseID            string
	ContractorID      string
	CustomerID        string
	PartnerID         string
	OwnerID           string
	OriginChannel     string
	Type              string
	Subject           string
	Priority          CasePriority
	Transactions      []Transaction
	Comments          []Comment
	Status            CaseStatus
	DueDate           time.Time
	CreatedBy         string
	CreatedAt         time.Time
	UpdatedBy         string
	UpdatedAt         time.Time
	Region            int
	ProductID         string
	ClosedAt          *time.Time
	ExternalReference string
	TargetDate        *time.Time
}

type CaseFull struct {
	CaseID            string
	Contractor        Contractor
	Customer          Customer
	Partner           Partner
	OwnerID           string
	OriginChannel     string
	Type              string
	Subject           string
	Priority          CasePriority
	Transactions      []Transaction
	Comments          []Comment
	Status            CaseStatus
	DueDate           time.Time
	CreatedBy         string
	CreatedAt         time.Time
	UpdatedBy         string
	UpdatedAt         time.Time
	Region            int
	Product           Product
	ClosedAt          *time.Time
	ExternalReference string
	TargetDate        *time.Time
}

type CaseFilters struct {
	OwnerID           []string
	PartnerID         []string
	ContractorID      []string
	CustomerID        []string
	Status            []string
	Region            []string
	ExternalReference []string
	PagingFilter
}

type CaseUpdate struct {
	Status     *CaseStatus
	PartnerID  *string
	OwnerID    *string
	TargetDate *time.Time
	ClosedAt   *time.Time
	CustomerID *string
	ProductID  *string
	Subject    *string
	UpdatedBy  string
}

type CaseStatus string

const (
	NEW             CaseStatus = "New"
	CUSTOMER_INFO   CaseStatus = "CustomerInfo"
	WAITING_PARTNER CaseStatus = "WaitingPartner"
	ONGOING         CaseStatus = "Ongoing"
	REPORT          CaseStatus = "Report"
	PAYMENT         CaseStatus = "Payment"
	RECEIPT         CaseStatus = "Receipt"
	CLOSED          CaseStatus = "Closed"
	CANCELED        CaseStatus = "Canceled"
	DRAFT           CaseStatus = "Draft"
)

type CasePriority string

const (
	LOW    CasePriority = "Low"
	MEDIUM CasePriority = "Medium"
	HIGH   CasePriority = "High"
)

func NewCase(
	contractorID string,
	customerID string,
	originChannel string,
	caseType string,
	subject string,
	dueDate time.Time,
	author string,
	externalReference string,
) (Case, error) {
	now := time.Now().UTC()

	caseID, err := uuid.NewUUID()
	if err != nil {
		return Case{}, err
	}

	return Case{
		CaseID:            caseID.String(),
		ContractorID:      contractorID,
		CustomerID:        customerID,
		OriginChannel:     originChannel,
		Type:              caseType,
		Subject:           subject,
		Priority:          MEDIUM,
		Status:            NEW,
		DueDate:           dueDate,
		CreatedAt:         now,
		CreatedBy:         author,
		UpdatedAt:         now,
		UpdatedBy:         author,
		ExternalReference: externalReference,
	}, nil
}

func (c *Case) MergeUpdate(updateCase CaseUpdate) {
	c.UpdatedAt = time.Now().UTC()
	c.UpdatedBy = updateCase.UpdatedBy

	if updateCase.Status != nil {
		c.Status = *updateCase.Status
	}

	if updateCase.OwnerID != nil {
		c.OwnerID = *updateCase.OwnerID
	}

	if updateCase.PartnerID != nil {
		c.PartnerID = *updateCase.PartnerID
	}

	if updateCase.TargetDate != nil {
		c.TargetDate = updateCase.TargetDate
	}

	if updateCase.ClosedAt != nil {
		c.ClosedAt = updateCase.ClosedAt
	}

	if updateCase.CustomerID != nil {
		c.CustomerID = *updateCase.CustomerID
	}

	if updateCase.ProductID != nil {
		c.ProductID = *updateCase.ProductID
	}

	if updateCase.Subject != nil {
		c.Subject = *updateCase.Subject
	}
}

func NewCaseFull(crmCase Case, comments []Comment, transactions []Transaction, product Product, customer Customer, partner Partner, contractor Contractor) CaseFull {
	return CaseFull{
		CaseID:            crmCase.CaseID,
		Contractor:        contractor,
		Customer:          customer,
		Partner:           partner,
		OwnerID:           crmCase.OwnerID,
		OriginChannel:     crmCase.OriginChannel,
		Type:              crmCase.Type,
		Subject:           crmCase.Subject,
		Priority:          crmCase.Priority,
		Transactions:      transactions,
		Comments:          comments,
		Status:            crmCase.Status,
		DueDate:           crmCase.DueDate,
		CreatedBy:         crmCase.CreatedBy,
		CreatedAt:         crmCase.CreatedAt,
		UpdatedBy:         crmCase.UpdatedBy,
		UpdatedAt:         crmCase.CreatedAt,
		Region:            crmCase.Region,
		Product:           product,
		ClosedAt:          crmCase.ClosedAt,
		ExternalReference: crmCase.ExternalReference,
		TargetDate:        crmCase.TargetDate,
	}
}
