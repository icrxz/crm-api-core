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

func newCommentServiceForTest(t *testing.T) (
	CommentService,
	*mock_domain.MockCommentRepository,
	*mock_domain.MockAttachmentRepository,
	*mock_domain.MockCommentHistoryRepository,
) {
	t.Helper()

	ctrl := gomock.NewController(t)

	commentRepository := mock_domain.NewMockCommentRepository(ctrl)
	attachmentRepository := mock_domain.NewMockAttachmentRepository(ctrl)
	commentHistoryRepository := mock_domain.NewMockCommentHistoryRepository(ctrl)
	transactionManager := mock_domain.NewMockTransactionManager(ctrl)
	transactionManager.EXPECT().
		WithinTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()

	service := NewCommentService(
		commentRepository,
		attachmentRepository,
		nil,
		commentHistoryRepository,
		transactionManager,
	)

	return service, commentRepository, attachmentRepository, commentHistoryRepository
}

func TestCommentService_AddAttachment(t *testing.T) {
	t.Run("returns a validation error when commentID is empty", func(t *testing.T) {
		service, _, _, _ := newCommentServiceForTest(t)

		_, err := service.AddAttachment(context.Background(), "", domain.Attachment{})

		require.Error(t, err)
	})

	t.Run("returns the comment repository error when the comment does not exist", func(t *testing.T) {
		service, commentRepository, _, _ := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(nil, domain.NewNotFoundError("comment not found", nil))

		_, err := service.AddAttachment(context.Background(), testCommentID, domain.Attachment{AttachmentID: testAttachmentID})

		require.Error(t, err)
	})

	t.Run("saves the attachment, touches the comment and records history", func(t *testing.T) {
		service, commentRepository, attachmentRepository, commentHistoryRepository := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, Content: "original content"}, nil)

		attachmentRepository.EXPECT().
			Save(gomock.Any(), domain.Attachment{AttachmentID: testAttachmentID, CommentID: testCommentID, CreatedBy: testAuthor}).
			Return(nil)

		commentRepository.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "original content", testAuthor).
			Return(nil)

		commentHistoryRepository.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, history domain.CommentHistory) error {
				require.Equal(t, domain.CommentAttachmentAddedEvent, history.EventName)
				require.Equal(t, testCommentID, history.CommentID)
				require.Equal(t, testAuthor, history.AuthorID)
				return nil
			})

		saved, err := service.AddAttachment(
			context.Background(),
			testCommentID,
			domain.Attachment{AttachmentID: testAttachmentID, CreatedBy: testAuthor},
		)

		require.NoError(t, err)
		require.Equal(t, testCommentID, saved.CommentID)
	})
}

func TestCommentService_UpdateContent(t *testing.T) {
	t.Run("returns a validation error when commentID is empty", func(t *testing.T) {
		service, _, _, _ := newCommentServiceForTest(t)

		err := service.UpdateContent(context.Background(), "", "pagamento realizado no dia 01/01/2024", testAuthor)

		require.Error(t, err)
	})

	t.Run("updates the comment content and records history", func(t *testing.T) {
		service, commentRepository, _, commentHistoryRepository := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, Content: "old content"}, nil)

		commentRepository.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "pagamento realizado no dia 01/01/2024", testAuthor).
			Return(nil)

		commentHistoryRepository.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ context.Context, history domain.CommentHistory) error {
				require.Equal(t, domain.CommentContentUpdatedEvent, history.EventName)
				require.Equal(t, "old content", history.OldValues["content"])
				require.Equal(t, "pagamento realizado no dia 01/01/2024", history.NewValues["content"])
				return nil
			})

		err := service.UpdateContent(context.Background(), testCommentID, "pagamento realizado no dia 01/01/2024", testAuthor)

		require.NoError(t, err)
	})

	t.Run("is a no-op when the content did not change", func(t *testing.T) {
		service, commentRepository, _, _ := newCommentServiceForTest(t)

		commentRepository.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, Content: "same content"}, nil)

		err := service.UpdateContent(context.Background(), testCommentID, "same content", testAuthor)

		require.NoError(t, err)
	})
}
