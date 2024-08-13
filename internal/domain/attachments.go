package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AttachmentRepository interface {
	Save(ctx context.Context, attachment Attachment) error
	SaveBatch(ctx context.Context, attachments []Attachment) error
	GetByID(ctx context.Context, attachmentID string) (Attachment, error)
	GetByCommentID(ctx context.Context, commentID string) ([]Attachment, error)
}

type AttachmentBucket interface {
	Download(ctx context.Context, attachmentID string) ([]byte, error)
}

type Attachment struct {
	AttachmentID  string
	CommentID     string
	Key           string
	FileName      string
	AttachmentURL string
	FileExtension string
	CreatedAt     time.Time
	Size          int
	CreatedBy     string
}

func NewAttachment(
	fileName string,
	attachmentURL string,
	fileExtension string,
	fileKey string,
	createdBy string,
	size int,
) (Attachment, error) {
	now := time.Now().UTC()
	attachmentID, err := uuid.NewUUID()
	if err != nil {
		return Attachment{}, err
	}

	return Attachment{
		AttachmentID:  attachmentID.String(),
		FileName:      fileName,
		AttachmentURL: attachmentURL,
		FileExtension: fileExtension,
		CreatedBy:     createdBy,
		CreatedAt:     now,
		Size:          size,
		Key:           fileKey,
	}, nil
}
