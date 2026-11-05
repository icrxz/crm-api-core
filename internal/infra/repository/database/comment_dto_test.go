package database

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
)

const testOngoingStep = "service_order"

func TestCommentDTOMetadata(t *testing.T) {
	t.Run("round-trips the comment metadata", func(t *testing.T) {
		comment := domain.Comment{
			CommentID: "comment-1",
			Metadata:  map[string]any{"ongoing_step": testOngoingStep},
		}

		dto := mapCommentToCommentDTO(comment)
		assert.Equal(t, JSONMap{"ongoing_step": testOngoingStep}, dto.Metadata)

		stored, err := dto.Metadata.Value()
		assert.NoError(t, err)
		assert.JSONEq(t, `{"ongoing_step":"service_order"}`, stored.(string))

		assert.Equal(t, comment.Metadata, mapCommentDTOToComment(dto).Metadata)
	})

	t.Run("stores an empty object for a comment without metadata", func(t *testing.T) {
		stored, err := mapCommentToCommentDTO(domain.Comment{}).Metadata.Value()

		assert.NoError(t, err)
		assert.Equal(t, "{}", stored)
	})
}
