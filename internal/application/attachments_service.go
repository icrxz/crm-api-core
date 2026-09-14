package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type attachmentService struct {
	attachmentRepository     domain.AttachmentRepository
	attachmentBucket         domain.AttachmentBucket
	commentRepository        domain.CommentRepository
	commentHistoryRepository domain.CommentHistoryRepository
	transactionManager       domain.TransactionManager
}

//go:generate mockgen -source=attachments_service.go -destination=mock_application/mock_attachments_service.go -package=mock_application
type AttachmentService interface {
	GetByID(ctx context.Context, attachmentID string) (*domain.Attachment, error)
	SearchByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error)
	DeleteByComments(ctx context.Context, commentIDs []string) error
	DeleteByID(ctx context.Context, attachmentID string, deletedBy string) error
}

func NewAttachmentService(
	attachmentRepository domain.AttachmentRepository,
	attachmentBucket domain.AttachmentBucket,
	commentRepository domain.CommentRepository,
	commentHistoryRepository domain.CommentHistoryRepository,
	transactionManager domain.TransactionManager,
) AttachmentService {
	return &attachmentService{
		attachmentRepository:     attachmentRepository,
		attachmentBucket:         attachmentBucket,
		commentRepository:        commentRepository,
		commentHistoryRepository: commentHistoryRepository,
		transactionManager:       transactionManager,
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

func (s *attachmentService) DeleteByID(ctx context.Context, attachmentID string, deletedBy string) error {
	if attachmentID == "" {
		return domain.NewValidationError("attachmentID is required", nil)
	}

	attachment, err := s.attachmentRepository.GetByID(ctx, attachmentID)
	if err != nil {
		return err
	}

	comment, err := s.commentRepository.GetByID(ctx, attachment.CommentID)
	if err != nil {
		return err
	}

	err = s.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.attachmentRepository.DeleteByID(txCtx, attachmentID); err != nil {
			return err
		}

		if err := s.commentRepository.UpdateContent(txCtx, comment.CommentID, comment.Content, deletedBy); err != nil {
			return err
		}

		history, err := domain.NewCommentHistory(
			attachment.CommentID,
			domain.CommentAttachmentRemovedEvent,
			deletedBy,
			map[string]any{
				"attachment_id":  attachment.AttachmentID,
				"file_name":      attachment.FileName,
				"attachment_url": attachment.AttachmentURL,
			},
			map[string]any{},
		)
		if err != nil {
			return err
		}

		return s.commentHistoryRepository.Create(txCtx, history)
	})
	if err != nil {
		return err
	}

	if err := s.attachmentBucket.Delete(ctx, attachment.Key); err != nil {
		slog.WarnContext(ctx, "failed to delete attachment file from bucket, leaving it orphaned",
			"attachment_id", attachmentID, "key", attachment.Key, "error", err)
	}

	return nil
}
