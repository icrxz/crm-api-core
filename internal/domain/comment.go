package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context, comment Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (*Comment, error)
	GetByCaseID(ctx context.Context, caseID string) ([]Comment, error)
}

type Attachment struct {
	AttachmentID  string
	FileName      string
	AttachmentURL string
	FileExtension string
	CreatedAt     time.Time
}

type Comment struct {
	CommentID   string
	CaseID      string
	Content     string
	CommentType CommentType
	Attachment  []Attachment
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type CommentType string

const (
	CONTENT    CommentType = "Content"
	COMMENT    CommentType = "Comment"
	RESOLUTION CommentType = "Resolution"
	REJECTION  CommentType = "Rejection"
)

func NewComment(
	caseID string,
	content string,
	createdBy string,
	commentType CommentType,
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
	}, nil
}
