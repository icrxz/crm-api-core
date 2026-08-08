package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueue(t *testing.T) {
	t.Run("creates a queue successfully", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", Criteria{"category": "mobile", "state": []string{"SP", "RJ"}}, "author-1")

		require.NoError(t, err)
		assert.NotEmpty(t, queue.QueueID)
		assert.Equal(t, "SP Mobile", queue.Name)
		assert.Equal(t, "mobile", queue.Criteria["category"])
		assert.Equal(t, []string{"SP", "RJ"}, queue.Criteria["state"])
		assert.True(t, queue.Active)
		assert.Equal(t, "author-1", queue.CreatedBy)
		assert.Equal(t, "author-1", queue.UpdatedBy)
		assert.False(t, queue.CreatedAt.IsZero())
		assert.Equal(t, queue.CreatedAt, queue.UpdatedAt)
	})

	t.Run("returns validation error when name is empty", func(t *testing.T) {
		_, err := NewQueue("", Criteria{"category": "mobile"}, "author-1")

		require.Error(t, err)
		assert.IsType(t, &CustomError{}, err)
	})

	t.Run("returns validation error when criteria is empty", func(t *testing.T) {
		_, err := NewQueue("SP Mobile", nil, "author-1")

		require.Error(t, err)
		assert.IsType(t, &CustomError{}, err)
	})
}

func TestQueue_MergeUpdate(t *testing.T) {
	t.Run("updates only the provided fields", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", Criteria{"category": "mobile", "state": []string{"SP"}}, "author-1")
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
		assert.Equal(t, "mobile", queue.Criteria["category"])
		assert.Equal(t, []string{"SP"}, queue.Criteria["state"])
		assert.Equal(t, "author-2", queue.UpdatedBy)
	})

	t.Run("replaces criteria when provided", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", Criteria{"category": "mobile", "state": []string{"SP"}}, "author-1")
		require.NoError(t, err)

		newCriteria := Criteria{"category": "digital", "state": []string{"RJ", "MG"}}

		queue.MergeUpdate(UpdateQueue{
			Criteria:  newCriteria,
			UpdatedBy: "author-2",
		})

		assert.Equal(t, newCriteria, queue.Criteria)
	})

	t.Run("leaves fields untouched when update is empty", func(t *testing.T) {
		queue, err := NewQueue("SP Mobile", Criteria{"category": "mobile", "state": []string{"SP"}}, "author-1")
		require.NoError(t, err)

		queue.MergeUpdate(UpdateQueue{UpdatedBy: "author-2"})

		assert.Equal(t, "SP Mobile", queue.Name)
		assert.Equal(t, "mobile", queue.Criteria["category"])
		assert.Equal(t, []string{"SP"}, queue.Criteria["state"])
		assert.True(t, queue.Active)
	})
}

func TestCriteria_Matches(t *testing.T) {
	t.Run("matches scalar equality", func(t *testing.T) {
		criteria := Criteria{"category": "mobile"}
		assert.True(t, criteria.Matches(map[string]any{"category": "mobile"}))
		assert.False(t, criteria.Matches(map[string]any{"category": "digital"}))
	})

	t.Run("matches when case field is a member of the criteria list", func(t *testing.T) {
		criteria := Criteria{"status": []string{"New", "Ongoing"}}
		assert.True(t, criteria.Matches(map[string]any{"status": "Ongoing"}))
		assert.False(t, criteria.Matches(map[string]any{"status": "Closed"}))
	})

	t.Run("matches when the case field list intersects the criteria list", func(t *testing.T) {
		criteria := Criteria{"state": []string{"SP", "RJ"}}
		assert.True(t, criteria.Matches(map[string]any{"state": []string{"RJ", "MG"}}))
		assert.False(t, criteria.Matches(map[string]any{"state": []string{"BA"}}))
	})

	t.Run("requires every criteria key to match (AND)", func(t *testing.T) {
		criteria := Criteria{"category": "mobile", "state": []string{"SP"}}
		assert.True(t, criteria.Matches(map[string]any{"category": "mobile", "state": []string{"SP"}}))
		assert.False(t, criteria.Matches(map[string]any{"category": "mobile", "state": []string{"RJ"}}))
	})

	t.Run("fails when the field is missing entirely", func(t *testing.T) {
		criteria := Criteria{"partner_id": "partner-1"}
		assert.False(t, criteria.Matches(map[string]any{}))
	})
}
