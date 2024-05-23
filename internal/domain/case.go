package domain

import (
	"context"
	"time"
)

type CaseRepository interface {
	Create(ctx context.Context, crmCase Case) (string, error)
	GetByID(ctx context.Context, caseID string) (*Case, error)
	Search(ctx context.Context) ([]Case, error)
	Update(ctx context.Context, crmCase Case) error
}

type Case struct {
	CaseID        string
	ContractorID  string
	CustomerID    string
	PartnerID     string
	OwnerID       string
	OriginChannel string
	Type          string
	Subject       string
	Priority      CasePriority
	Transactions  []Transaction
	Comments      []Comment
	Status        CaseStatus
	DueDate       time.Time
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedBy     string
	UpdatedAt     time.Time
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
