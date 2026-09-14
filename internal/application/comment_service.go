package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type commentService struct {
	commentRepository        domain.CommentRepository
	attachmentRepository     domain.AttachmentRepository
	attachmentBucket         domain.AttachmentBucket
	commentHistoryRepository domain.CommentHistoryRepository
	transactionManager       domain.TransactionManager
}

//go:generate mockgen -source=comment_service.go -destination=mock_application/mock_comment_service.go -package=mock_application
type CommentService interface {
	Create(ctx context.Context, comment domain.Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (*domain.Comment, error)
	GetByCaseID(ctx context.Context, caseID string) ([]domain.Comment, error)
	DeleteByCaseID(ctx context.Context, caseID string) error
	AddAttachment(ctx context.Context, commentID string, attachment domain.Attachment) (domain.Attachment, error)
	UpdateContent(ctx context.Context, commentID string, content string, updatedBy string) error
}

func NewCommentService(
	commentRepository domain.CommentRepository,
	attachmentRepository domain.AttachmentRepository,
	attachmentBucket domain.AttachmentBucket,
	commentHistoryRepository domain.CommentHistoryRepository,
	transactionManager domain.TransactionManager,
) CommentService {
	return &commentService{
		commentRepository:        commentRepository,
		attachmentRepository:     attachmentRepository,
		attachmentBucket:         attachmentBucket,
		commentHistoryRepository: commentHistoryRepository,
		transactionManager:       transactionManager,
	}
}

func (s *commentService) Create(ctx context.Context, comment domain.Comment) (string, error) {
	var commentID string

	err := s.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		createdID, err := s.commentRepository.Create(txCtx, comment)
		if err != nil {
			return err
		}
		commentID = createdID

		if len(comment.Attachments) == 0 {
			return nil
		}

		for idx := range comment.Attachments {
			comment.Attachments[idx].CommentID = commentID
		}

		return s.attachmentRepository.SaveBatch(txCtx, comment.Attachments)
	})
	if err != nil {
		return "", err
	}

	return commentID, nil
}

func (s *commentService) GetByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("commentID is required", nil)
	}

	comment, err := s.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	comment.Attachments, err = s.attachmentRepository.GetByCommentID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) GetByCaseID(ctx context.Context, caseID string) ([]domain.Comment, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("caseID is required", nil)
	}

	comments, err := s.commentRepository.GetByCaseID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	for idx, comment := range comments {
		attachments, err := s.attachmentRepository.GetByCommentID(ctx, comment.CommentID)
		if err != nil {
			return nil, err
		}

		comments[idx].Attachments = attachments
	}

	return comments, nil
}

func (s *commentService) DeleteByCaseID(ctx context.Context, caseID string) error {
	if caseID == "" {
		return domain.NewValidationError("caseID is required", nil)
	}

	return s.commentRepository.DeleteManyByCaseID(ctx, caseID)
}

func (s *commentService) AddAttachment(ctx context.Context, commentID string, attachment domain.Attachment) (domain.Attachment, error) {
	if commentID == "" {
		return domain.Attachment{}, domain.NewValidationError("commentID is required", nil)
	}

	existingComment, err := s.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return domain.Attachment{}, err
	}

	attachment.CommentID = commentID

	err = s.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.attachmentRepository.Save(txCtx, attachment); err != nil {
			return err
		}

		if err := s.commentRepository.UpdateContent(txCtx, commentID, existingComment.Content, attachment.CreatedBy); err != nil {
			return err
		}

		history, err := domain.NewCommentHistory(
			commentID,
			domain.CommentAttachmentAddedEvent,
			attachment.CreatedBy,
			map[string]any{},
			map[string]any{
				"attachment_id":  attachment.AttachmentID,
				"file_name":      attachment.FileName,
				"attachment_url": attachment.AttachmentURL,
			},
		)
		if err != nil {
			return err
		}

		return s.commentHistoryRepository.Create(txCtx, history)
	})
	if err != nil {
		return domain.Attachment{}, err
	}

	return attachment, nil
}

func (s *commentService) UpdateContent(ctx context.Context, commentID string, content string, updatedBy string) error {
	if commentID == "" {
		return domain.NewValidationError("commentID is required", nil)
	}

	existingComment, err := s.commentRepository.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	if existingComment.Content == content {
		return nil
	}

	return s.transactionManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.commentRepository.UpdateContent(txCtx, commentID, content, updatedBy); err != nil {
			return err
		}

		history, err := domain.NewCommentHistory(
			commentID,
			domain.CommentContentUpdatedEvent,
			updatedBy,
			map[string]any{"content": existingComment.Content},
			map[string]any{"content": content},
		)
		if err != nil {
			return err
		}

		return s.commentHistoryRepository.Create(txCtx, history)
	})
}
