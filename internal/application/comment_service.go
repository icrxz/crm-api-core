package application

import (
	"context"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type commentService struct {
	commentRepository    domain.CommentRepository
	attachmentRepository domain.AttachmentRepository
}

type CommentService interface {
	Create(ctx context.Context, comment domain.Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (*domain.Comment, error)
	GetByCaseID(ctx context.Context, caseID string) ([]domain.Comment, error)
}

func NewCommentService(commentRepository domain.CommentRepository, attachmentRepository domain.AttachmentRepository) CommentService {
	return &commentService{
		commentRepository:    commentRepository,
		attachmentRepository: attachmentRepository,
	}
}

func (s *commentService) Create(ctx context.Context, comment domain.Comment) (string, error) {
	commentID, err := s.commentRepository.Create(ctx, comment)
	if err != nil {
		return "", err
	}

	if comment.Attachments != nil && len(comment.Attachments) > 0 {
		for idx := range comment.Attachments {
			comment.Attachments[idx].CommentID = commentID
		}

		err = s.attachmentRepository.SaveBatch(ctx, comment.Attachments)
		if err != nil {
			return commentID, err
		}
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
