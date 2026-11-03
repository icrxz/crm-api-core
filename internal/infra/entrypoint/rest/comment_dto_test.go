package rest

import (
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testStepMetadataKey = "ongoing_step"

func TestCommentDTOMetadata(t *testing.T) {
	t.Run("keeps the metadata sent on creation", func(t *testing.T) {
		comment, err := mapCreateCommentDTOToComment(CreateCommentDTO{
			Content:     "Ordem de serviço",
			CommentType: domain.COMMENT_RESOLUTION,
			CreatedBy:   "tecnico",
			Metadata:    map[string]any{testStepMetadataKey: "service_order"},
		}, "case-1")

		require.NoError(t, err)
		assert.Equal(t, map[string]any{testStepMetadataKey: "service_order"}, comment.Metadata)
	})

	t.Run("defaults to empty metadata when none is sent", func(t *testing.T) {
		comment, err := mapCreateCommentDTOToComment(CreateCommentDTO{
			Content:     "avulso",
			CommentType: domain.COMMENT,
			CreatedBy:   "tecnico",
		}, "case-1")

		require.NoError(t, err)
		assert.Equal(t, map[string]any{}, comment.Metadata)
	})

	t.Run("exposes the metadata on the response", func(t *testing.T) {
		dto := mapCommentToCommentDTO(domain.Comment{
			CommentID: "comment-1",
			Metadata:  map[string]any{testStepMetadataKey: "technician_photos"},
		})

		assert.Equal(t, map[string]any{testStepMetadataKey: "technician_photos"}, dto.Metadata)
	})
}
