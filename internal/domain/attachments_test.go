package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestAttachment(documentType AttachmentDocumentType) (Attachment, error) {
	return NewAttachment("os.png", "https://bucket/os.png", "image/png", "key-1", "partner-1", 1024, documentType)
}

func TestNewAttachment_KeepsTheGivenDocumentType(t *testing.T) {
	attachment, err := newTestAttachment(DOCUMENT_SERVICE_ORDER)

	assert.NoError(t, err)
	assert.Equal(t, DOCUMENT_SERVICE_ORDER, attachment.DocumentType)
}

// Attachments created before document_type existed, and callers that leave the
// field out, must stay valid and fall back to a plain photo.
func TestNewAttachment_DefaultsEmptyDocumentTypeToPhoto(t *testing.T) {
	attachment, err := newTestAttachment("")

	assert.NoError(t, err)
	assert.Equal(t, DOCUMENT_PHOTO, attachment.DocumentType)
}

func TestNewAttachment_RejectsUnknownDocumentType(t *testing.T) {
	_, err := newTestAttachment("contrato")

	assert.Error(t, err)
}
