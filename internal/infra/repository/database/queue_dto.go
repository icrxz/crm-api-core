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
