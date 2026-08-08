package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

//go:generate mockgen -source=queue_resolver.go -destination=mock_application/mock_queue_resolver.go -package=mock_application
type QueueResolver interface {
	Resolve(ctx context.Context, crmCase domain.Case) (*domain.Queue, error)
}

type queueResolver struct {
	queueService QueueService
}

func NewQueueResolver(queueService QueueService) QueueResolver {
	return &queueResolver{
		queueService: queueService,
	}
}

// Resolve finds the first active queue whose criteria matches the case's
// typed fields and metadata. Returns nil, nil when no queue matches.
func (r *queueResolver) Resolve(ctx context.Context, crmCase domain.Case) (*domain.Queue, error) {
	active := true
	queues, err := r.queueService.Search(ctx, domain.QueueFilters{
		Active: &active,
		PagingFilter: domain.PagingFilter{
			Limit:  100,
			Offset: 0,
		},
	})
	if err != nil {
		return nil, err
	}

	fields := crmCase.MatchableFields()
	for _, queue := range queues.Result {
		if queue.Criteria.Matches(fields) {
			matched := queue
			return &matched, nil
		}
	}

	return nil, nil
}
