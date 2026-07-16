package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=queue.go -destination=mock_domain/mock_queue_repository.go -package=mock_domain

type QueueCategory string

const (
	MobileQueueCategory  QueueCategory = "mobile"
	DigitalQueueCategory QueueCategory = "digital"
)

type QueueRepository interface {
	Create(ctx context.Context, queue Queue) (string, error)
	GetByID(ctx context.Context, queueID string) (*Queue, error)
	Search(ctx context.Context, filters QueueFilters) (PagingResult[Queue], error)
	Update(ctx context.Context, queue Queue) error
	Delete(ctx context.Context, queueID string) error
	AddMember(ctx context.Context, queueID string, userID string) error
	RemoveMember(ctx context.Context, queueID string, userID string) error
	GetMembers(ctx context.Context, queueID string) ([]User, error)
	GetQueuesByUser(ctx context.Context, userID string) ([]Queue, error)
}

type Queue struct {
	QueueID   string
	Name      string
	Category  QueueCategory
	States    []string
	Active    bool
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
}

type QueueFilters struct {
	QueueID  []string
	Category []string
	State    []string
	Active   *bool
	PagingFilter
}

type UpdateQueue struct {
	Name      *string
	Category  *QueueCategory
	States    []string
	Active    *bool
	UpdatedBy string
}

func NewQueue(name string, category QueueCategory, states []string, author string) (Queue, error) {
	if name == "" {
		return Queue{}, NewValidationError("name cannot be empty", nil)
	}

	if category == "" {
		return Queue{}, NewValidationError("category cannot be empty", nil)
	}

	now := time.Now().UTC()
	queueID, err := uuid.NewUUID()
	if err != nil {
		return Queue{}, err
	}

	return Queue{
		QueueID:   queueID.String(),
		Name:      name,
		Category:  category,
		States:    states,
		Active:    true,
		CreatedBy: author,
		CreatedAt: now,
		UpdatedBy: author,
		UpdatedAt: now,
	}, nil
}

func (q *Queue) MergeUpdate(update UpdateQueue) {
	q.UpdatedBy = update.UpdatedBy
	q.UpdatedAt = time.Now().UTC()

	if update.Name != nil {
		q.Name = *update.Name
	}

	if update.Category != nil {
		q.Category = *update.Category
	}

	if update.States != nil {
		q.States = update.States
	}

	if update.Active != nil {
		q.Active = *update.Active
	}
}
