package database

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMapCaseFullDTOToCaseFull(t *testing.T) {
	t.Run("maps the joined queue onto the case full", func(t *testing.T) {
		queueID := "queue-1"
		name := "SP Mobile"

		dto := CaseFullDTO{
			CaseID: "case-1",
			Queue: QueueOptionalDTO{
				QueueID: &queueID,
				Name:    &name,
				States:  pq.StringArray{"SP"},
			},
		}

		caseFull := mapCaseFullDTOToCaseFull(dto)

		assert.Equal(t, "case-1", caseFull.CaseID)
		assert.Equal(t, "queue-1", caseFull.Queue.QueueID)
		assert.Equal(t, "SP Mobile", caseFull.Queue.Name)
	})

	t.Run("leaves the queue empty when there is no joined row", func(t *testing.T) {
		dto := CaseFullDTO{CaseID: "case-1"}

		caseFull := mapCaseFullDTOToCaseFull(dto)

		assert.Empty(t, caseFull.Queue.QueueID)
	})
}
