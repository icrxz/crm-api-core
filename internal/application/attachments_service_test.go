package application

import (
	"context"
	"errors"
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/icrxz/crm-api-core/internal/domain/mock_domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const testDeletedBy = "user-2"

func newAttachmentServiceForTest(t *testing.T) (
	AttachmentService,
	*mock_domain.MockAttachmentRepository,
	*mock_domain.MockAttachmentBucket,
	*mock_domain.MockCommentRepository,
	*mock_domain.MockCommentHistoryRepository,
) {
	t.Helper()

	ctrl := gomock.NewController(t)
	attachmentRepository := mock_domain.NewMockAttachmentRepository(ctrl)
	attachmentBucket := mock_domain.NewMockAttachmentBucket(ctrl)
	commentRepository := mock_domain.NewMockCommentRepository(ctrl)
	commentHistoryRepository := mock_domain.NewMockCommentHistoryRepository(ctrl)
	transactionManager := mock_domain.NewMockTransactionManager(ctrl)
	transactionManager.EXPECT().
		WithinTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()

	service := NewAttachmentService(
		attachmentRepository,
		attachmentBucket,
		commentRepository,
		commentHistoryRepository,
		transactionManager,
	)

	return service, attachmentRepository, attachmentBucket, commentRepository, commentHistoryRepository
}

func TestAttachmentService_DeleteByID(t *testing.T) {
	t.Run("returns a validation error when attachmentID is empty", func(t *testing.T) {
		service, _, _, _, _ := newAttachmentServiceForTest(t)

		err := service.DeleteByID(context.Background(), "", testDeletedBy)

		require.Error(t, err)
	})

	t.Run("deletes the database record, touches the comment, records history and deletes the bucket file", func(t *testing.T) {
		service, attachmentRepository, attachmentBucket, commentRepository, commentHistoryRepository := newAttachmentServiceForTest(t)

		attachmentRepository.EXPECT().
			GetByID(gomock.Any(), testAttachmentID).
			Return(domain.Attachment{AttachmentID: testAttachmentID, CommentID: testCommentID, Key: "file.png"}, nil)
		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, Content: "original content"}, nil)
		attachmentRepository.EXPECT().DeleteByID(gomock.Any(), testAttachmentID).Return(nil)
		commentRepository.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "original content", testDeletedBy).
			Return(nil)
		commentHistoryRepository.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, history domain.CommentHistory) error {
				require.Equal(t, domain.CommentAttachmentRemovedEvent, history.EventName)
				require.Equal(t, testCommentID, history.CommentID)
				require.Equal(t, testDeletedBy, history.AuthorID)
				return nil
			})
		attachmentBucket.EXPECT().Delete(gomock.Any(), "file.png").Return(nil)

		err := service.DeleteByID(context.Background(), testAttachmentID, testDeletedBy)

		require.NoError(t, err)
	})

	t.Run("still succeeds when the bucket file cannot be deleted", func(t *testing.T) {
		service, attachmentRepository, attachmentBucket, commentRepository, commentHistoryRepository := newAttachmentServiceForTest(t)

		attachmentRepository.EXPECT().
			GetByID(gomock.Any(), testAttachmentID).
			Return(domain.Attachment{AttachmentID: testAttachmentID, CommentID: testCommentID, Key: "file.png"}, nil)
		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID}, nil)
		attachmentRepository.EXPECT().DeleteByID(gomock.Any(), testAttachmentID).Return(nil)
		commentRepository.EXPECT().UpdateContent(gomock.Any(), testCommentID, "", testDeletedBy).Return(nil)
		commentHistoryRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		attachmentBucket.EXPECT().Delete(gomock.Any(), "file.png").Return(errors.New("s3 unavailable"))

		err := service.DeleteByID(context.Background(), testAttachmentID, testDeletedBy)

		require.NoError(t, err)
	})

	t.Run("does not delete the database record when the attachment is not found", func(t *testing.T) {
		service, attachmentRepository, _, _, _ := newAttachmentServiceForTest(t)

		attachmentRepository.EXPECT().
			GetByID(gomock.Any(), testAttachmentID).
			Return(domain.Attachment{}, domain.NewNotFoundError("attachment not found", nil))

		err := service.DeleteByID(context.Background(), testAttachmentID, testDeletedBy)

		require.Error(t, err)
	})
}
