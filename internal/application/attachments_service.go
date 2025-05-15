package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type attachmentService struct {
	attachmentRepository domain.AttachmentRepository
}

type AttachmentService interface {
	GetByID(ctx context.Context, attachmentID string) (*domain.Attachment, error)
	SearchByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error)
	DeleteByComments(ctx context.Context, commentIDs []string) error
}

func NewAttachmentService(attachmentRepository domain.AttachmentRepository) AttachmentService {
	return &attachmentService{
		attachmentRepository: attachmentRepository,
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
		return domain.NewValidationError("commentIDs is required to delete attachments", nil)
	}

	return s.attachmentRepository.DeleteManyByComments(ctx, commentIDs)
}
