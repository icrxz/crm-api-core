package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type AttachmentDTO struct {
	AttachmentID  string    `db:"attachment_id"`
	CommentID     string    `db:"comment_id"`
	Key           string    `db:"key"`
	FileName      string    `db:"file_name"`
	AttachmentURL string    `db:"attachment_url"`
	FileExtension string    `db:"file_extension"`
	Size          int       `db:"size"`
	CreatedAt     time.Time `db:"created_at"`
	CreatedBy     string    `db:"created_by"`
}

func mapAttachmentToAttachmentDTO(attachment domain.Attachment) AttachmentDTO {
	return AttachmentDTO{
		AttachmentID:  attachment.AttachmentID,
		CommentID:     attachment.CommentID,
		Key:           attachment.Key,
		FileName:      attachment.FileName,
		AttachmentURL: attachment.AttachmentURL,
		FileExtension: attachment.FileExtension,
		Size:          attachment.Size,
		CreatedAt:     attachment.CreatedAt,
		CreatedBy:     attachment.CreatedBy,
	}
}

func mapAttachmentDTOToAttachment(attachmentDTO AttachmentDTO) domain.Attachment {
	return domain.Attachment{
		AttachmentID:  attachmentDTO.AttachmentID,
		CommentID:     attachmentDTO.CommentID,
		Key:           attachmentDTO.Key,
		FileName:      attachmentDTO.FileName,
		AttachmentURL: attachmentDTO.AttachmentURL,
		FileExtension: attachmentDTO.FileExtension,
		Size:          attachmentDTO.Size,
		CreatedAt:     attachmentDTO.CreatedAt,
		CreatedBy:     attachmentDTO.CreatedBy,
	}
}

func mapAttachmentsDTOToAttachments(attachmentsDTO []AttachmentDTO) []domain.Attachment {
	attachments := make([]domain.Attachment, 0, len(attachmentsDTO))
	for _, attachmentDTO := range attachmentsDTO {
		attachments = append(attachments, mapAttachmentDTOToAttachment(attachmentDTO))
	}
	return attachments
}

func mapAttachmentsToAttachmentsDTO(attachments []domain.Attachment) []AttachmentDTO {
	attachmentsDTO := make([]AttachmentDTO, 0, len(attachments))
	for _, attachment := range attachments {
		attachmentsDTO = append(attachmentsDTO, mapAttachmentToAttachmentDTO(attachment))
	}
	return attachmentsDTO
}
