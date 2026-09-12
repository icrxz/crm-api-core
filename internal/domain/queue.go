package domain

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=queue.go -destination=mock_domain/mock_queue_repository.go -package=mock_domain

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

// Criteria is a generic AND-matcher evaluated against a case's typed fields
// plus its metadata. A value can be a scalar (equality) or a list (membership -
// matches if the case field equals one of the values, or, when the case field
// is itself a list, if the two lists intersect).
type Criteria map[string]any

func (c Criteria) Matches(fields map[string]any) bool {
	for key, expected := range c {
		actual, ok := fields[key]
		if !ok || !criteriaValueMatches(expected, actual) {
			return false
		}
	}

	return true
}

func criteriaValueMatches(expected any, actual any) bool {
	expectedList, expectedIsList := toAnySlice(expected)
	actualList, actualIsList := toAnySlice(actual)

	switch {
	case expectedIsList && actualIsList:
		return anySliceIntersects(expectedList, actualList)
	case expectedIsList:
		return anySliceContains(expectedList, actual)
	case actualIsList:
		return anySliceContains(actualList, expected)
	default:
		return fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
	}
}

func toAnySlice(value any) ([]any, bool) {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, false
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}

	return result, true
}

func anySliceContains(list []any, value any) bool {
	for _, item := range list {
		if fmt.Sprintf("%v", item) == fmt.Sprintf("%v", value) {
			return true
		}
	}

	return false
}

func anySliceIntersects(a []any, b []any) bool {
	for _, item := range a {
		if anySliceContains(b, item) {
			return true
		}
	}

	return false
}

type Queue struct {
	QueueID        string
	Name           string
	Criteria       Criteria
	Active         bool
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedBy      string
	UpdatedAt      time.Time
	OrganizationID *string
}

type QueueFilters struct {
	QueueID  []string
	Criteria Criteria
	Active   *bool
	PagingFilter
}

type UpdateQueue struct {
	Name           *string
	Criteria       Criteria
	Active         *bool
	OrganizationID *string
	UpdatedBy      string
}

func NewQueue(name string, criteria Criteria, author string) (Queue, error) {
	if name == "" {
		return Queue{}, NewValidationError("name cannot be empty", nil)
	}

	if len(criteria) == 0 {
		return Queue{}, NewValidationError("criteria cannot be empty", nil)
	}

	now := time.Now().UTC()
	queueID, err := uuid.NewUUID()
	if err != nil {
		return Queue{}, err
	}

	return Queue{
		QueueID:   queueID.String(),
		Name:      name,
		Criteria:  criteria,
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

	if update.Criteria != nil {
		q.Criteria = update.Criteria
	}

	if update.Active != nil {
		q.Active = *update.Active
	}

	if update.OrganizationID != nil {
		q.OrganizationID = update.OrganizationID
	}
}
