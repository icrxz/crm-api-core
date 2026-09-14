package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=comment.go -destination=mock_domain/mock_comment_repository.go -package=mock_domain
type CommentRepository interface {
	Create(ctx context.Context, comment Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (*Comment, error)
	GetByCaseID(ctx context.Context, caseID string) ([]Comment, error)
	DeleteManyByCaseID(ctx context.Context, caseID string) error
	UpdateContent(ctx context.Context, commentID string, content string, updatedBy string) error
}

type Comment struct {
	CommentID   string
	CaseID      string
	Content     string
	CommentType CommentType
	Attachments []Attachment
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type CommentType string

const (
	COMMENT_CONTENT       CommentType = "Content"
	COMMENT               CommentType = "Comment"
	COMMENT_RESOLUTION    CommentType = "Resolution"
	COMMENT_REPORT        CommentType = "Report"
	COMMENT_REJECTION     CommentType = "Rejection"
	COMMENT_PAYMENT_PROOF CommentType = "PaymentProof"
)

func NewComment(
	caseID string,
	content string,
	createdBy string,
	commentType CommentType,
	attachments []Attachment,
) (Comment, error) {
	now := time.Now().UTC()
	commentID, err := uuid.NewUUID()
	if err != nil {
		return Comment{}, err
	}

	return Comment{
		CommentID:   commentID.String(),
		CommentType: commentType,
		CaseID:      caseID,
		Content:     content,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		Attachments: attachments,
	}, nil
}
