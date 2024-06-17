package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type CaseRepository interface {
	Create(ctx context.Context, crmCase Case) (string, error)
	GetByID(ctx context.Context, caseID string) (*Case, error)
	Search(ctx context.Context, filters CaseFilters) ([]Case, error)
	Update(ctx context.Context, crmCase Case) error
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
}

type CaseFilters struct {
	OwnerID      []string
	PartnerID    []string
	ContractorID []string
	CustomerID   []string
	Status       []string
	Region       []string
}

type CaseStatus string

const (
	NEW             CaseStatus = "New"
	CUSTOMER_INFO   CaseStatus = "CustomerInfo"
	WAITING_PARTNER CaseStatus = "WaitingPartner"
	ONGOING         CaseStatus = "Ongoing"
	REPORT          CaseStatus = "Report"
	PAYMENT         CaseStatus = "Payment"
	CLOSED          CaseStatus = "Closed"
	CANCELED        CaseStatus = "Canceled"
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
