package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateQueueDTO struct {
	Name      string         `json:"name" validate:"required"`
	Criteria  map[string]any `json:"criteria" validate:"required"`
	CreatedBy string         `json:"created_by"`
}

type QueueDTO struct {
	QueueID   string         `json:"queue_id"`
	Name      string         `json:"name"`
	Criteria  map[string]any `json:"criteria"`
	Active    bool           `json:"active"`
	CreatedBy string         `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedBy string         `json:"updated_by"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type UpdateQueueDTO struct {
	Name      *string        `json:"name"`
	Criteria  map[string]any `json:"criteria"`
	Active    *bool          `json:"active"`
	UpdatedBy string         `json:"updated_by"`
}

type AddQueueMemberDTO struct {
	UserID string `json:"user_id" validate:"required"`
}

func mapQueueToQueueDTO(queue domain.Queue) QueueDTO {
	return QueueDTO{
		QueueID:   queue.QueueID,
		Name:      queue.Name,
		Criteria:  map[string]any(queue.Criteria),
		Active:    queue.Active,
		CreatedBy: queue.CreatedBy,
		CreatedAt: queue.CreatedAt,
		UpdatedBy: queue.UpdatedBy,
		UpdatedAt: queue.UpdatedAt,
	}
}

func mapCreateQueueDTOToQueue(queueDTO CreateQueueDTO) (domain.Queue, error) {
	return domain.NewQueue(
		queueDTO.Name,
		domain.Criteria(queueDTO.Criteria),
		queueDTO.CreatedBy,
	)
}

func mapQueuesToQueueDTOs(queues []domain.Queue) []QueueDTO {
	queueDTOs := make([]QueueDTO, 0, len(queues))
	for _, queue := range queues {
		queueDTOs = append(queueDTOs, mapQueueToQueueDTO(queue))
	}

	return queueDTOs
}

func mapUpdateQueueDTOToUpdateQueue(queueDTO UpdateQueueDTO) domain.UpdateQueue {
	return domain.UpdateQueue{
		Name:      queueDTO.Name,
		Criteria:  domain.Criteria(queueDTO.Criteria),
		Active:    queueDTO.Active,
		UpdatedBy: queueDTO.UpdatedBy,
	}
}
