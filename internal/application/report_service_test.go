package application

import (
	"testing"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/require"
)

func reportTestComment(commentType domain.CommentType, createdAt time.Time, content string, attachments ...domain.Attachment) domain.Comment {
	return domain.Comment{
		CommentType: commentType,
		CreatedAt:   createdAt,
		Content:     content,
		Attachments: attachments,
	}
}

func reportTestAttachment(key string, documentType domain.AttachmentDocumentType) domain.Attachment {
	return domain.Attachment{Key: key, DocumentType: documentType}
}

func attachmentKeys(attachments []domain.Attachment) []string {
	keys := make([]string, 0, len(attachments))
	for _, attachment := range attachments {
		keys = append(keys, attachment.Key)
	}
	return keys
}

func TestReportPhotos(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 10, d, 10, 0, 0, 0, time.UTC) }
	v2Case := domain.Case{Metadata: map[string]any{ongoingFlowVersionMetadataKey: ongoingFlowVersionV2}}

	t.Run("original flow uses only the latest resolution, whatever the slice order", func(t *testing.T) {
		comments := []domain.Comment{
			reportTestComment(domain.COMMENT_RESOLUTION, day(5), "", reportTestAttachment("new.png", domain.DOCUMENT_PHOTO)),
			reportTestComment(domain.COMMENT_RESOLUTION, day(1), "", reportTestAttachment("rejected.png", domain.DOCUMENT_PHOTO)),
		}

		require.Equal(t, []string{"new.png"}, attachmentKeys(reportPhotos(domain.Case{}, comments)))
	})

	t.Run("leaves the paperwork out of the image slots", func(t *testing.T) {
		comments := []domain.Comment{
			reportTestComment(domain.COMMENT_RESOLUTION, day(1), "",
				reportTestAttachment("photo.png", domain.DOCUMENT_PHOTO),
				reportTestAttachment("os.png", domain.DOCUMENT_SERVICE_ORDER),
				reportTestAttachment("legacy.png", ""),
				reportTestAttachment("term.png", domain.DOCUMENT_WARRANTY_TERM),
			),
		}

		require.Equal(t, []string{"photo.png", "legacy.png"}, attachmentKeys(reportPhotos(domain.Case{}, comments)))
	})

	t.Run("v2 flow collects photos from every resolution in creation order", func(t *testing.T) {
		comments := []domain.Comment{
			reportTestComment(domain.COMMENT_RESOLUTION, day(3), "Envio para o laudo"),
			reportTestComment(domain.COMMENT_RESOLUTION, day(2), "Ordem de Serviço", reportTestAttachment("os.png", domain.DOCUMENT_SERVICE_ORDER)),
			reportTestComment(domain.COMMENT_RESOLUTION, day(1), "Fotos", reportTestAttachment("a.png", domain.DOCUMENT_PHOTO)),
			reportTestComment(domain.COMMENT, day(1), "avulso", reportTestAttachment("comment.png", domain.DOCUMENT_PHOTO)),
			reportTestComment(domain.COMMENT_RESOLUTION, day(2), "Mais fotos", reportTestAttachment("b.png", domain.DOCUMENT_PHOTO)),
		}

		require.Equal(t, []string{"a.png", "b.png"}, attachmentKeys(reportPhotos(v2Case, comments)))
	})

	t.Run("v2 flow does not repeat a file sent twice", func(t *testing.T) {
		comments := []domain.Comment{
			reportTestComment(domain.COMMENT_RESOLUTION, day(1), "", reportTestAttachment("a.png", domain.DOCUMENT_PHOTO)),
			reportTestComment(domain.COMMENT_RESOLUTION, day(2), "", reportTestAttachment("a.png", domain.DOCUMENT_PHOTO)),
		}

		require.Equal(t, []string{"a.png"}, attachmentKeys(reportPhotos(v2Case, comments)))
	})

	t.Run("no resolution yields no photos", func(t *testing.T) {
		comments := []domain.Comment{reportTestComment(domain.COMMENT, day(1), "", reportTestAttachment("a.png", domain.DOCUMENT_PHOTO))}

		require.Empty(t, reportPhotos(v2Case, comments))
		require.Empty(t, reportPhotos(domain.Case{}, nil))
	})
}

func TestReportConclusion(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 10, d, 10, 0, 0, 0, time.UTC) }

	comments := []domain.Comment{
		reportTestComment(domain.COMMENT_REPORT, day(4), "segunda conclusão"),
		reportTestComment(domain.COMMENT_RESOLUTION, day(5), "não é conclusão"),
		reportTestComment(domain.COMMENT_REPORT, day(2), "primeira conclusão"),
	}

	require.Equal(t, "segunda conclusão", reportConclusion(comments))
	require.Empty(t, reportConclusion(nil))
}
