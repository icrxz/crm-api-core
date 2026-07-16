package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateQueueDTO struct {
	Name      string   `json:"name" validate:"required"`
	Category  string   `json:"category" validate:"required"`
	States    []string `json:"states"`
	CreatedBy string   `json:"created_by"`
}

type QueueDTO struct {
	QueueID   string    `json:"queue_id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	States    []string  `json:"states"`
	Active    bool      `json:"active"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateQueueDTO struct {
	Name      *string  `json:"name"`
	Category  *string  `json:"category"`
	States    []string `json:"states"`
	Active    *bool    `json:"active"`
	UpdatedBy string   `json:"updated_by"`
}

type AddQueueMemberDTO struct {
	UserID string `json:"user_id" validate:"required"`
}

func mapQueueToQueueDTO(queue domain.Queue) QueueDTO {
	return QueueDTO{
		QueueID:   queue.QueueID,
		Name:      queue.Name,
		Category:  string(queue.Category),
		States:    queue.States,
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
		domain.QueueCategory(queueDTO.Category),
		queueDTO.States,
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
	var category *domain.QueueCategory
	if queueDTO.Category != nil {
		parsedCategory := domain.QueueCategory(*queueDTO.Category)
		category = &parsedCategory
	}

	return domain.UpdateQueue{
		Name:      queueDTO.Name,
		Category:  category,
		States:    queueDTO.States,
		Active:    queueDTO.Active,
		UpdatedBy: queueDTO.UpdatedBy,
	}
}
