package application

import (
	"context"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type attachmentService struct {
	attachmentRepository domain.AttachmentRepository
	attachmentBucket     domain.AttachmentBucket
}

type AttachmentService interface {
	GetByID(ctx context.Context, attachmentID string) (*domain.Attachment, error)
	SearchByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error)
	DeleteByComments(ctx context.Context, commentIDs []string) error
}

func NewAttachmentService(
	attachmentRepository domain.AttachmentRepository,
	attachmentBucket domain.AttachmentBucket,
) AttachmentService {
	return &attachmentService{
		attachmentRepository: attachmentRepository,
		attachmentBucket:     attachmentBucket,
	}
}

func (s *attachmentService) GetByID(ctx context.Context, attachmentID string) (*domain.Attachment, error) {
	if attachmentID == "" {
		return nil, domain.NewValidationError("attachmentID is required", nil)
	}

	attachment, err := s.attachmentRepository.GetByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}

func (s *attachmentService) SearchByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("commentID is required to search attachments", nil)
	}

	foundAttachments, err := s.attachmentRepository.GetByCommentID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return foundAttachments, nil
}

func (s *attachmentService) DeleteByComments(ctx context.Context, commentIDs []string) error {
	if len(commentIDs) == 0 {
		fmt.Println("DeleteByComments called with empty commentIDs, nothing to delete")
		return nil
	}

	return s.attachmentRepository.DeleteManyByComments(ctx, commentIDs)
}
