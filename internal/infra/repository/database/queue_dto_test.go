package database

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMapQueueOptionalDTOToQueue(t *testing.T) {
	t.Run("returns an empty queue when the left join found nothing", func(t *testing.T) {
		queue := mapQueueOptionalDTOToQueue(QueueOptionalDTO{})

		assert.Empty(t, queue.QueueID)
		assert.Empty(t, queue.Name)
	})

	t.Run("maps a fully populated joined queue", func(t *testing.T) {
		queueID := "queue-1"
		name := "SP Mobile"
		category := "mobile"
		active := true
		createdBy := "author-1"
		updatedBy := "author-2"
		createdAt := time.Now().UTC()
		updatedAt := createdAt.Add(time.Hour)

		dto := QueueOptionalDTO{
			QueueID:   &queueID,
			Name:      &name,
			Category:  &category,
			States:    pq.StringArray{"SP", "RJ"},
			Active:    &active,
			CreatedBy: &createdBy,
			CreatedAt: &createdAt,
			UpdatedBy: &updatedBy,
			UpdatedAt: &updatedAt,
		}

		queue := mapQueueOptionalDTOToQueue(dto)

		assert.Equal(t, queueID, queue.QueueID)
		assert.Equal(t, name, queue.Name)
		assert.EqualValues(t, category, queue.Category)
		assert.Equal(t, []string{"SP", "RJ"}, queue.States)
		assert.True(t, queue.Active)
		assert.Equal(t, createdBy, queue.CreatedBy)
		assert.Equal(t, updatedBy, queue.UpdatedBy)
		assert.Equal(t, createdAt, queue.CreatedAt)
		assert.Equal(t, updatedAt, queue.UpdatedAt)
	})
}
