package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type QueueDTO struct {
	QueueID        string    `db:"queue_id"`
	Name           string    `db:"name"`
	Criteria       JSONMap   `db:"criteria"`
	Active         bool      `db:"active"`
	CreatedBy      string    `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedBy      string    `db:"updated_by"`
	UpdatedAt      time.Time `db:"updated_at"`
	OrganizationID *string   `db:"organization_id"`
}

func mapQueueToQueueDTO(queue domain.Queue) QueueDTO {
	return QueueDTO{
		QueueID:        queue.QueueID,
		Name:           queue.Name,
		Criteria:       JSONMap(queue.Criteria),
		Active:         queue.Active,
		CreatedBy:      queue.CreatedBy,
		CreatedAt:      queue.CreatedAt,
		UpdatedBy:      queue.UpdatedBy,
		UpdatedAt:      queue.UpdatedAt,
		OrganizationID: queue.OrganizationID,
	}
}

func mapQueueDTOToQueue(queueDTO QueueDTO) domain.Queue {
	return domain.Queue{
		QueueID:        queueDTO.QueueID,
		Name:           queueDTO.Name,
		Criteria:       domain.Criteria(queueDTO.Criteria),
		Active:         queueDTO.Active,
		CreatedBy:      queueDTO.CreatedBy,
		CreatedAt:      queueDTO.CreatedAt,
		UpdatedBy:      queueDTO.UpdatedBy,
		UpdatedAt:      queueDTO.UpdatedAt,
		OrganizationID: queueDTO.OrganizationID,
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
	QueueID   *string    `db:"queue_id"`
	Name      *string    `db:"name"`
	Criteria  JSONMap    `db:"criteria"`
	Active    *bool      `db:"active"`
	CreatedBy *string    `db:"created_by"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedBy *string    `db:"updated_by"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func mapQueueOptionalDTOToQueue(queueDTO QueueOptionalDTO) domain.Queue {
	if queueDTO.QueueID == nil {
		return domain.Queue{}
	}

	queue := domain.Queue{
		QueueID:  *queueDTO.QueueID,
		Criteria: domain.Criteria(queueDTO.Criteria),
	}

	if queueDTO.Name != nil {
		queue.Name = *queueDTO.Name
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
