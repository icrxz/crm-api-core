package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type commentService struct {
	commentRepository domain.CommentRepository
}

type CommentService interface {
	Create(ctx context.Context, comment domain.Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (*domain.Comment, error)
	GetByCaseID(ctx context.Context, caseID string) ([]domain.Comment, error)
}

func NewCommentService(commentRepository domain.CommentRepository) CommentService {
	return &commentService{
		commentRepository: commentRepository,
	}
}

func (s *commentService) Create(ctx context.Context, comment domain.Comment) (string, error) {
	return s.commentRepository.Create(ctx, comment)
}

func (s *commentService) GetByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("commentID is required", nil)
	}
	return s.commentRepository.GetByID(ctx, commentID)
}

func (s *commentService) GetByCaseID(ctx context.Context, caseID string) ([]domain.Comment, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("caseID is required", nil)
	}
	return s.commentRepository.GetByCaseID(ctx, caseID)
}
