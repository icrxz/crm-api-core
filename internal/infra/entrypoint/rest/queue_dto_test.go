package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapUpdateQueueDTOToUpdateQueue(t *testing.T) {
	category := "digital"

	update := mapUpdateQueueDTOToUpdateQueue(UpdateQueueDTO{
		Category:  &category,
		UpdatedBy: "author-2",
	})

	assert.NotNil(t, update.Category)
	assert.Equal(t, "digital", string(*update.Category))
	assert.Equal(t, "author-2", update.UpdatedBy)
}
