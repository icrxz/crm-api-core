package application

import (
	"context"
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/icrxz/crm-api-core/internal/domain/mock_domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testCommentID    = "comment-1"
	testAttachmentID = "attachment-1"
)

func newCommentServiceForTest(t *testing.T) (CommentService, *mock_domain.MockCommentRepository, *mock_domain.MockAttachmentRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)

	commentRepository := mock_domain.NewMockCommentRepository(ctrl)
	attachmentRepository := mock_domain.NewMockAttachmentRepository(ctrl)

	service := NewCommentService(
		commentRepository,
		attachmentRepository,
		nil,
		mock_domain.NewMockTransactionManager(ctrl),
	)

	return service, commentRepository, attachmentRepository
}

func TestCommentService_AddAttachment(t *testing.T) {
	t.Run("returns a validation error when commentID is empty", func(t *testing.T) {
		service, _, _ := newCommentServiceForTest(t)

		_, err := service.AddAttachment(context.Background(), "", domain.Attachment{})

		require.Error(t, err)
	})

	t.Run("returns the comment repository error when the comment does not exist", func(t *testing.T) {
		service, commentRepository, _ := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(nil, domain.NewNotFoundError("comment not found", nil))

		_, err := service.AddAttachment(context.Background(), testCommentID, domain.Attachment{AttachmentID: testAttachmentID})

		require.Error(t, err)
	})

	t.Run("saves the attachment linked to the comment", func(t *testing.T) {
		service, commentRepository, attachmentRepository := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID}, nil)

		attachmentRepository.EXPECT().
			Save(gomock.Any(), domain.Attachment{AttachmentID: testAttachmentID, CommentID: testCommentID}).
			Return(nil)

		saved, err := service.AddAttachment(context.Background(), testCommentID, domain.Attachment{AttachmentID: testAttachmentID})

		require.NoError(t, err)
		require.Equal(t, testCommentID, saved.CommentID)
	})
}

func TestCommentService_UpdateContent(t *testing.T) {
	t.Run("returns a validation error when commentID is empty", func(t *testing.T) {
		service, _, _ := newCommentServiceForTest(t)

		err := service.UpdateContent(context.Background(), "", "pagamento realizado no dia 01/01/2024", testAuthor)

		require.Error(t, err)
	})

	t.Run("updates the comment content", func(t *testing.T) {
		service, commentRepository, _ := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "pagamento realizado no dia 01/01/2024", testAuthor).
			Return(nil)

		err := service.UpdateContent(context.Background(), testCommentID, "pagamento realizado no dia 01/01/2024", testAuthor)

		require.NoError(t, err)
	})
}
