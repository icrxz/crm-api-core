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
	CaseID       string
	ContractorID string
	CustomerID   string
	PartnerID    string
	OwnerID      string
	Origin       string
	Type         string
	Subject      string
	Priority     CasePriority
	Transactions []Transaction
	Comments     []Comment
	Status       CaseStatus
	DueDate      time.Time
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    string
	UpdatedAt    time.Time
}

type CaseStatus string

const (
	PENDING   CaseStatus = "Pending"
	PROGRESS  CaseStatus = "Progress"
	HOLD      CaseStatus = "Hold"
	COMPLETED CaseStatus = "Completed"
	REJECTED  CaseStatus = "Rejected"
)

type CasePriority string

const (
	LOW    CasePriority = "Low"
	MEDIUM CasePriority = "Medium"
	HIGH   CasePriority = "High"
)
