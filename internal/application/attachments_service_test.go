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

func newAttachmentServiceForTest(t *testing.T) (AttachmentService, *mock_domain.MockAttachmentRepository, *mock_domain.MockAttachmentBucket) {
	t.Helper()

	ctrl := gomock.NewController(t)
	attachmentRepository := mock_domain.NewMockAttachmentRepository(ctrl)
	attachmentBucket := mock_domain.NewMockAttachmentBucket(ctrl)

	service := NewAttachmentService(attachmentRepository, attachmentBucket)

	return service, attachmentRepository, attachmentBucket
}

func TestAttachmentService_DeleteByID(t *testing.T) {
	t.Run("returns a validation error when attachmentID is empty", func(t *testing.T) {
		service, _, _ := newAttachmentServiceForTest(t)

		err := service.DeleteByID(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("deletes the database record and the bucket file", func(t *testing.T) {
		service, attachmentRepository, attachmentBucket := newAttachmentServiceForTest(t)

		attachmentRepository.EXPECT().
			GetByID(gomock.Any(), testAttachmentID).
			Return(domain.Attachment{AttachmentID: testAttachmentID, Key: "file.png"}, nil)
		attachmentRepository.EXPECT().DeleteByID(gomock.Any(), testAttachmentID).Return(nil)
		attachmentBucket.EXPECT().Delete(gomock.Any(), "file.png").Return(nil)

		err := service.DeleteByID(context.Background(), testAttachmentID)

		require.NoError(t, err)
	})

	t.Run("still succeeds when the bucket file cannot be deleted", func(t *testing.T) {
		service, attachmentRepository, attachmentBucket := newAttachmentServiceForTest(t)

		attachmentRepository.EXPECT().
			GetByID(gomock.Any(), testAttachmentID).
			Return(domain.Attachment{AttachmentID: testAttachmentID, Key: "file.png"}, nil)
		attachmentRepository.EXPECT().DeleteByID(gomock.Any(), testAttachmentID).Return(nil)
		attachmentBucket.EXPECT().Delete(gomock.Any(), "file.png").Return(errors.New("s3 unavailable"))

		err := service.DeleteByID(context.Background(), testAttachmentID)

		require.NoError(t, err)
	})

	t.Run("does not delete the database record when the attachment is not found", func(t *testing.T) {
		service, attachmentRepository, _ := newAttachmentServiceForTest(t)

		attachmentRepository.EXPECT().
			GetByID(gomock.Any(), testAttachmentID).
			Return(domain.Attachment{}, domain.NewNotFoundError("attachment not found", nil))

		err := service.DeleteByID(context.Background(), testAttachmentID)

		require.Error(t, err)
	})
}
