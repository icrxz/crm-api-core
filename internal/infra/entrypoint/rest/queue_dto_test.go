package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapUpdateQueueDTOToUpdateQueue(t *testing.T) {
	update := mapUpdateQueueDTOToUpdateQueue(UpdateQueueDTO{
		Criteria:  map[string]any{"category": "digital"},
		UpdatedBy: "author-2",
	})

	assert.NotNil(t, update.Criteria)
	assert.Equal(t, "digital", update.Criteria["category"])
	assert.Equal(t, "author-2", update.UpdatedBy)
}
