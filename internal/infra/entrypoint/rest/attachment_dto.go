package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type AttachmentDTO struct {
	AttachmentID  string    `json:"attachment_id"`
	FileExtension string    `json:"file_extension"`
	FileName      string    `json:"file_name"`
	Key           string    `json:"key"`
	Size          int       `json:"size"`
	AttachmentURL string    `json:"url"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateAttachmentDTO struct {
	FileExtension string `json:"file_extension"`
	FileName      string `json:"file_name"`
	Key           string `json:"key"`
	Size          int    `json:"size"`
	AttachmentURL string `json:"url"`
	CreatedBy     string `json:"created_by"`
}

func mapCreateAttachmentDTOToAttachment(createAttachmentDTO CreateAttachmentDTO) (domain.Attachment, error) {
	return domain.NewAttachment(
		createAttachmentDTO.FileName,
		createAttachmentDTO.AttachmentURL,
		createAttachmentDTO.FileExtension,
		createAttachmentDTO.Key,
		createAttachmentDTO.CreatedBy,
		createAttachmentDTO.Size,
	)
}

func mapAttachmentToAttachmentDTO(attachment domain.Attachment) AttachmentDTO {
	return AttachmentDTO{
		AttachmentID:  attachment.AttachmentID,
		FileExtension: attachment.FileExtension,
		FileName:      attachment.FileName,
		Size:          attachment.Size,
		AttachmentURL: attachment.AttachmentURL,
		CreatedBy:     attachment.CreatedBy,
		CreatedAt:     attachment.CreatedAt,
	}
}

func mapAttachmentsToAttachmentDTOs(attachments []domain.Attachment) []AttachmentDTO {
	var attachmentDTOs []AttachmentDTO
	for _, attachment := range attachments {
		attachmentDTOs = append(attachmentDTOs, mapAttachmentToAttachmentDTO(attachment))
	}
	return attachmentDTOs
}

func mapCreateAttachmentDTOsToAttachments(attachmentDTOs []CreateAttachmentDTO) ([]domain.Attachment, error) {
	var attachments []domain.Attachment
	for _, attachmentDTO := range attachmentDTOs {
		attachment, err := mapCreateAttachmentDTOToAttachment(attachmentDTO)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, attachment)
	}
	return attachments, nil
}
