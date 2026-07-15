package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CaseHistoryRepository interface {
	Create(ctx context.Context, history CaseHistory) error
	GetByCaseID(ctx context.Context, caseID string) ([]CaseHistory, error)
}

type CaseHistory struct {
	HistoryID string
	CaseID    string
	EventName string
	AuthorID  string
	OldValues map[string]any
	NewValues map[string]any
	CreatedAt time.Time
}

const (
	CaseCreatedEvent           = "case_created"
	CaseAssignedEvent          = "case_assigned"
	CaseOwnerChangedEvent      = "case_owner_changed"
	CaseStatusChangedEvent     = "case_status_changed"
	CaseClosedEvent            = "case_closed"
	CasePartnerChangedEvent    = "case_partner_changed"
	CaseTargetDateChangedEvent = "case_target_date_changed"
	CaseDetailsUpdatedEvent    = "case_details_updated"
	CaseUpdatedEvent           = "case_updated"
	CaseResetEvent             = "case_reset"
)

func NewCaseHistory(
	caseID string,
	eventName string,
	authorID string,
	oldValues map[string]any,
	newValues map[string]any,
) (CaseHistory, error) {
	historyID, err := uuid.NewUUID()
	if err != nil {
		return CaseHistory{}, err
	}

	return CaseHistory{
		HistoryID: historyID.String(),
		CaseID:    caseID,
		EventName: eventName,
		AuthorID:  authorID,
		OldValues: oldValues,
		NewValues: newValues,
		CreatedAt: time.Now().UTC(),
	}, nil
}
