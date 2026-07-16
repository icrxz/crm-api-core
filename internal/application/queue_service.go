package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type queueService struct {
	queueRepository domain.QueueRepository
}

//go:generate mockgen -source=queue_service.go -destination=mock_application/mock_queue_service.go -package=mock_application
type QueueService interface {
	Create(ctx context.Context, queue domain.Queue) (string, error)
	GetByID(ctx context.Context, queueID string) (*domain.Queue, error)
	Update(ctx context.Context, queueID string, update domain.UpdateQueue) error
	Delete(ctx context.Context, queueID string) error
	Search(ctx context.Context, filters domain.QueueFilters) (domain.PagingResult[domain.Queue], error)
	AddMember(ctx context.Context, queueID string, userID string) error
	RemoveMember(ctx context.Context, queueID string, userID string) error
	GetMembers(ctx context.Context, queueID string) ([]domain.User, error)
	GetQueuesByUser(ctx context.Context, userID string) ([]domain.Queue, error)
}

func NewQueueService(queueRepository domain.QueueRepository) QueueService {
	return &queueService{
		queueRepository: queueRepository,
	}
}

func (s *queueService) Create(ctx context.Context, queue domain.Queue) (string, error) {
	return s.queueRepository.Create(ctx, queue)
}

func (s *queueService) GetByID(ctx context.Context, queueID string) (*domain.Queue, error) {
	if queueID == "" {
		return nil, domain.NewValidationError("queueID cannot be empty", nil)
	}

	return s.queueRepository.GetByID(ctx, queueID)
}

func (s *queueService) Update(ctx context.Context, queueID string, update domain.UpdateQueue) error {
	if queueID == "" {
		return domain.NewValidationError("queueID cannot be empty", nil)
	}

	queue, err := s.GetByID(ctx, queueID)
	if err != nil {
		return err
	}

	queue.MergeUpdate(update)

	return s.queueRepository.Update(ctx, *queue)
}

func (s *queueService) Delete(ctx context.Context, queueID string) error {
	if queueID == "" {
		return domain.NewValidationError("queueID cannot be empty", nil)
	}

	return s.queueRepository.Delete(ctx, queueID)
}

func (s *queueService) Search(ctx context.Context, filters domain.QueueFilters) (domain.PagingResult[domain.Queue], error) {
	return s.queueRepository.Search(ctx, filters)
}

func (s *queueService) AddMember(ctx context.Context, queueID string, userID string) error {
	if queueID == "" || userID == "" {
		return domain.NewValidationError("queueID and userID cannot be empty", nil)
	}

	return s.queueRepository.AddMember(ctx, queueID, userID)
}

func (s *queueService) RemoveMember(ctx context.Context, queueID string, userID string) error {
	if queueID == "" || userID == "" {
		return domain.NewValidationError("queueID and userID cannot be empty", nil)
	}

	return s.queueRepository.RemoveMember(ctx, queueID, userID)
}

func (s *queueService) GetMembers(ctx context.Context, queueID string) ([]domain.User, error) {
	if queueID == "" {
		return nil, domain.NewValidationError("queueID cannot be empty", nil)
	}

	return s.queueRepository.GetMembers(ctx, queueID)
}

func (s *queueService) GetQueuesByUser(ctx context.Context, userID string) ([]domain.Queue, error) {
	if userID == "" {
		return nil, domain.NewValidationError("userID cannot be empty", nil)
	}

	return s.queueRepository.GetQueuesByUser(ctx, userID)
}
