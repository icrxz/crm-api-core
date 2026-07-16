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
	SearchFull(ctx context.Context, filters CaseFilters) (PagingResult[CaseFull], error)
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
	QueueID           string
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
	QueueID           string
}

type CaseFilters struct {
	CaseID            []string
	OwnerID           []string
	PartnerID         []string
	ContractorID      []string
	CustomerID        []string
	Status            []string
	Region            []string
	ExternalReference []string
	StartDate         *string
	EndDate           *string
	ClosedAtStart     *string
	ClosedAtEnd       *string
	ShippingState     []string
	QueueID           []string
	PagingFilter
}

type CaseUpdate struct {
	Status     *CaseStatus
	PartnerID  *string
	OwnerID    *string
	TargetDate *time.Time
	ClosedAt   *time.Time
	Type       *string
	CustomerID *string
	ProductID  *string
	Subject    *string
	QueueID    *string
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
	REJECTED        CaseStatus = "Rejected"
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

	if updateCase.Type != nil {
		c.Type = *updateCase.Type
	}

	if updateCase.QueueID != nil {
		c.QueueID = *updateCase.QueueID
	}
}

// caseDiff accumulates the old/new values of fields changed by a CaseUpdate.
type caseDiff struct {
	oldValues map[string]any
	newValues map[string]any
	changed   map[string]bool
}

func newCaseDiff() *caseDiff {
	return &caseDiff{
		oldValues: map[string]any{},
		newValues: map[string]any{},
		changed:   map[string]bool{},
	}
}

func (d *caseDiff) recordString(field string, current string, update *string) {
	if update == nil || *update == current {
		return
	}

	d.oldValues[field] = current
	d.newValues[field] = *update
	d.changed[field] = true
}

func (d *caseDiff) recordTime(field string, current *time.Time, update *time.Time) {
	if update == nil || update.Equal(derefTime(current)) {
		return
	}

	d.oldValues[field] = current
	d.newValues[field] = *update
	d.changed[field] = true
}

func (d *caseDiff) recordStatus(current CaseStatus, update *CaseStatus) {
	if update == nil || *update == current {
		return
	}

	d.oldValues["status"] = current
	d.newValues["status"] = *update
	d.changed["status"] = true
}

// DetectChanges compares the current case state against an update payload and
// returns a semantic event name plus the old/new values of every field that
// actually changed. It must be called before MergeUpdate mutates the case.
func (c *Case) DetectChanges(update CaseUpdate) (string, map[string]any, map[string]any) {
	diff := newCaseDiff()

	diff.recordStatus(c.Status, update.Status)
	diff.recordString("owner_id", c.OwnerID, update.OwnerID)
	diff.recordString("partner_id", c.PartnerID, update.PartnerID)
	diff.recordTime("target_date", c.TargetDate, update.TargetDate)
	diff.recordTime("closed_at", c.ClosedAt, update.ClosedAt)
	diff.recordString("customer_id", c.CustomerID, update.CustomerID)
	diff.recordString("product_id", c.ProductID, update.ProductID)
	diff.recordString("subject", c.Subject, update.Subject)
	diff.recordString("type", c.Type, update.Type)
	diff.recordString("queue_id", c.QueueID, update.QueueID)

	return resolveCaseEventName(diff.changed), diff.oldValues, diff.newValues
}

// resolveCaseEventName picks a single semantic event name for a set of
// simultaneously changed fields, most specific business meaning first
// (e.g. status+owner changing together means the case was assigned).
func resolveCaseEventName(changed map[string]bool) string {
	switch {
	case changed["owner_id"] && changed["status"]:
		return CaseAssignedEvent
	case changed["owner_id"]:
		return CaseOwnerChangedEvent
	case changed["closed_at"]:
		return CaseClosedEvent
	case changed["status"]:
		return CaseStatusChangedEvent
	case changed["partner_id"]:
		return CasePartnerChangedEvent
	case changed["target_date"]:
		return CaseTargetDateChangedEvent
	case changed["queue_id"]:
		return CaseQueueChangedEvent
	case changed["customer_id"], changed["product_id"], changed["subject"], changed["type"]:
		return CaseDetailsUpdatedEvent
	default:
		return CaseUpdatedEvent
	}
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// Snapshot returns the initial field values recorded on case creation.
func (c Case) Snapshot() map[string]any {
	return map[string]any{
		"status":             c.Status,
		"owner_id":           c.OwnerID,
		"customer_id":        c.CustomerID,
		"contractor_id":      c.ContractorID,
		"partner_id":         c.PartnerID,
		"product_id":         c.ProductID,
		"type":               c.Type,
		"subject":            c.Subject,
		"priority":           c.Priority,
		"due_date":           c.DueDate,
		"external_reference": c.ExternalReference,
		"queue_id":           c.QueueID,
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
		UpdatedAt:         crmCase.UpdatedAt,
		Region:            crmCase.Region,
		Product:           product,
		ClosedAt:          crmCase.ClosedAt,
		ExternalReference: crmCase.ExternalReference,
		TargetDate:        crmCase.TargetDate,
		QueueID:           crmCase.QueueID,
	}
}
