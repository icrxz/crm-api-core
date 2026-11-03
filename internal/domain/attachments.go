package domain

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=attachments.go -destination=mock_domain/mock_attachment_repository.go -package=mock_domain
type AttachmentRepository interface {
	Save(ctx context.Context, attachment Attachment) error
	SaveBatch(ctx context.Context, attachments []Attachment) error
	GetByID(ctx context.Context, attachmentID string) (Attachment, error)
	GetByCommentID(ctx context.Context, commentID string) ([]Attachment, error)
	DeleteManyByComments(ctx context.Context, commentIDs []string) error
	DeleteByID(ctx context.Context, attachmentID string) error
}

type AttachmentBucket interface {
	Download(ctx context.Context, attachmentID string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

// AttachmentDocumentType tells apart the paperwork the partner has to collect on site
// (signed service order, customer confirmation, warranty term) from the plain
// pictures of the product. Everything uploaded before this existed is a PHOTO.
type AttachmentDocumentType string

const (
	DOCUMENT_PHOTO                 AttachmentDocumentType = "photo"
	DOCUMENT_SERVICE_ORDER         AttachmentDocumentType = "service_order"
	DOCUMENT_CUSTOMER_CONFIRMATION AttachmentDocumentType = "customer_confirmation"
	DOCUMENT_WARRANTY_TERM         AttachmentDocumentType = "warranty_term"
)

var attachmentDocumentTypes = []AttachmentDocumentType{
	DOCUMENT_PHOTO,
	DOCUMENT_SERVICE_ORDER,
	DOCUMENT_CUSTOMER_CONFIRMATION,
	DOCUMENT_WARRANTY_TERM,
}

func (d AttachmentDocumentType) IsValid() bool {
	return slices.Contains(attachmentDocumentTypes, d)
}

type Attachment struct {
	AttachmentID  string
	CommentID     string
	Key           string
	FileName      string
	AttachmentURL string
	FileExtension string
	DocumentType  AttachmentDocumentType
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
	documentType AttachmentDocumentType,
) (Attachment, error) {
	now := time.Now().UTC()
	attachmentID, err := uuid.NewUUID()
	if err != nil {
		return Attachment{}, err
	}

	if documentType == "" {
		documentType = DOCUMENT_PHOTO
	}

	if !documentType.IsValid() {
		return Attachment{}, NewValidationError("invalid document_type", map[string]any{"document_type": documentType})
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
		DocumentType:  documentType,
	}, nil
}
