package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/lib/pq"
)

type QueueDTO struct {
	QueueID   string         `db:"queue_id"`
	Name      string         `db:"name"`
	Category  string         `db:"category"`
	States    pq.StringArray `db:"states"`
	Active    bool           `db:"active"`
	CreatedBy string         `db:"created_by"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedBy string         `db:"updated_by"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func mapQueueToQueueDTO(queue domain.Queue) QueueDTO {
	return QueueDTO{
		QueueID:   queue.QueueID,
		Name:      queue.Name,
		Category:  string(queue.Category),
		States:    pq.StringArray(queue.States),
		Active:    queue.Active,
		CreatedBy: queue.CreatedBy,
		CreatedAt: queue.CreatedAt,
		UpdatedBy: queue.UpdatedBy,
		UpdatedAt: queue.UpdatedAt,
	}
}

func mapQueueDTOToQueue(queueDTO QueueDTO) domain.Queue {
	return domain.Queue{
		QueueID:   queueDTO.QueueID,
		Name:      queueDTO.Name,
		Category:  domain.QueueCategory(queueDTO.Category),
		States:    []string(queueDTO.States),
		Active:    queueDTO.Active,
		CreatedBy: queueDTO.CreatedBy,
		CreatedAt: queueDTO.CreatedAt,
		UpdatedBy: queueDTO.UpdatedBy,
		UpdatedAt: queueDTO.UpdatedAt,
	}
}

func mapQueueDTOsToQueues(queueDTOs []QueueDTO) []domain.Queue {
	queues := make([]domain.Queue, 0, len(queueDTOs))
	for _, queueDTO := range queueDTOs {
		queues = append(queues, mapQueueDTOToQueue(queueDTO))
	}

	return queues
}

// QueueOptionalDTO scans a LEFT JOIN'd queue that may not exist for a given case.
type QueueOptionalDTO struct {
	QueueID   *string        `db:"queue_id"`
	Name      *string        `db:"name"`
	Category  *string        `db:"category"`
	States    pq.StringArray `db:"states"`
	Active    *bool          `db:"active"`
	CreatedBy *string        `db:"created_by"`
	CreatedAt *time.Time     `db:"created_at"`
	UpdatedBy *string        `db:"updated_by"`
	UpdatedAt *time.Time     `db:"updated_at"`
}

func mapQueueOptionalDTOToQueue(queueDTO QueueOptionalDTO) domain.Queue {
	if queueDTO.QueueID == nil {
		return domain.Queue{}
	}

	queue := domain.Queue{
		QueueID: *queueDTO.QueueID,
		States:  []string(queueDTO.States),
	}

	if queueDTO.Name != nil {
		queue.Name = *queueDTO.Name
	}

	if queueDTO.Category != nil {
		queue.Category = domain.QueueCategory(*queueDTO.Category)
	}

	if queueDTO.Active != nil {
		queue.Active = *queueDTO.Active
	}

	if queueDTO.CreatedBy != nil {
		queue.CreatedBy = *queueDTO.CreatedBy
	}

	if queueDTO.CreatedAt != nil {
		queue.CreatedAt = *queueDTO.CreatedAt
	}

	if queueDTO.UpdatedBy != nil {
		queue.UpdatedBy = *queueDTO.UpdatedBy
	}

	if queueDTO.UpdatedAt != nil {
		queue.UpdatedAt = *queueDTO.UpdatedAt
	}

	return queue
}
