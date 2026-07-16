package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueue(t *testing.T) {
	t.Run("creates a queue successfully", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", MobileQueueCategory, []string{"SP", "RJ"}, "author-1")

		require.NoError(t, err)
		assert.NotEmpty(t, queue.QueueID)
		assert.Equal(t, "SP Mobile", queue.Name)
		assert.Equal(t, MobileQueueCategory, queue.Category)
		assert.Equal(t, []string{"SP", "RJ"}, queue.States)
		assert.True(t, queue.Active)
		assert.Equal(t, "author-1", queue.CreatedBy)
		assert.Equal(t, "author-1", queue.UpdatedBy)
		assert.False(t, queue.CreatedAt.IsZero())
		assert.Equal(t, queue.CreatedAt, queue.UpdatedAt)
	})

	t.Run("returns validation error when name is empty", func(t *testing.T) {
		_, err := NewQueue("", MobileQueueCategory, []string{"SP"}, "author-1")

		require.Error(t, err)
		assert.IsType(t, &CustomError{}, err)
	})

	t.Run("returns validation error when category is empty", func(t *testing.T) {
		_, err := NewQueue("SP Mobile", "", []string{"SP"}, "author-1")

		require.Error(t, err)
		assert.IsType(t, &CustomError{}, err)
	})

	t.Run("allows an empty states list", func(t *testing.T) {
		queue, err := NewQueue("Digital", DigitalQueueCategory, nil, "author-1")

		require.NoError(t, err)
		assert.Empty(t, queue.States)
	})
}

func TestQueue_MergeUpdate(t *testing.T) {
	t.Run("updates only the provided fields", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", MobileQueueCategory, []string{"SP"}, "author-1")
		require.NoError(t, err)

		newName := "SP Mobile Updated"
		inactive := false

		queue.MergeUpdate(UpdateQueue{
			Name:      &newName,
			Active:    &inactive,
			UpdatedBy: "author-2",
		})

		assert.Equal(t, newName, queue.Name)
		assert.False(t, queue.Active)
		assert.Equal(t, MobileQueueCategory, queue.Category)
		assert.Equal(t, []string{"SP"}, queue.States)
		assert.Equal(t, "author-2", queue.UpdatedBy)
	})

	t.Run("replaces category and states when provided", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", MobileQueueCategory, []string{"SP"}, "author-1")
		require.NoError(t, err)

		newCategory := DigitalQueueCategory
		newStates := []string{"RJ", "MG"}

		queue.MergeUpdate(UpdateQueue{
			Category:  &newCategory,
			States:    newStates,
			UpdatedBy: "author-2",
		})

		assert.Equal(t, DigitalQueueCategory, queue.Category)
		assert.Equal(t, newStates, queue.States)
	})

	t.Run("leaves fields untouched when update is empty", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", MobileQueueCategory, []string{"SP"}, "author-1")
		require.NoError(t, err)

		queue.MergeUpdate(UpdateQueue{UpdatedBy: "author-2"})

		assert.Equal(t, "SP Mobile", queue.Name)
		assert.Equal(t, MobileQueueCategory, queue.Category)
		assert.Equal(t, []string{"SP"}, queue.States)
		assert.True(t, queue.Active)
	})
}
