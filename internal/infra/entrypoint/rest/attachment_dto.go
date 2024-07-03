package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type AttachmentDTO struct {
	AttachmentID  string `json:"attachment_id"`
	FileExtension string `json:"file_extension"`
	FileName      string `json:"file_name"`
	AttachmentURL string `json:"url"`
}

func mapAttachmentDTOToAttachment(attachmentDTO AttachmentDTO) domain.Attachment {
	return domain.Attachment{
		AttachmentID:  attachmentDTO.AttachmentID,
		FileExtension: attachmentDTO.FileExtension,
		FileName:      attachmentDTO.FileName,
		AttachmentURL: attachmentDTO.AttachmentURL,
	}
}

func mapAttachmentToAttachmentDTO(attachment domain.Attachment) AttachmentDTO {
	return AttachmentDTO{
		AttachmentID:  attachment.AttachmentID,
		FileExtension: attachment.FileExtension,
		FileName:      attachment.FileName,
		AttachmentURL: attachment.AttachmentURL,
	}
}

func mapAttachmentsToAttachmentDTOs(attachments []domain.Attachment) []AttachmentDTO {
	var attachmentDTOs []AttachmentDTO
	for _, attachment := range attachments {
		attachmentDTOs = append(attachmentDTOs, mapAttachmentToAttachmentDTO(attachment))
	}
	return attachmentDTOs
}

func mapAttachmentDTOsToAttachments(attachmentDTOs []AttachmentDTO) []domain.Attachment {
	var attachments []domain.Attachment
	for _, attachmentDTO := range attachmentDTOs {
		attachments = append(attachments, mapAttachmentDTOToAttachment(attachmentDTO))
	}
	return attachments
}
