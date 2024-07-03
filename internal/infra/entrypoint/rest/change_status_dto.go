package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type ChangeStatusDTO struct {
	Status      string          `json:"status"`
	UpdatedBy   string          `json:"updated_by"`
	Content     *string         `json:"content"`
	Attachments []AttachmentDTO `json:"attachments"`
}

func mapChangeStatusDTOToChangeStatus(changeStatusDTO ChangeStatusDTO) domain.ChangeStatus {
	return domain.ChangeStatus{
		Status:      domain.CaseStatus(changeStatusDTO.Status),
		UpdatedBy:   changeStatusDTO.UpdatedBy,
		Content:     changeStatusDTO.Content,
		Attachments: mapAttachmentDTOsToAttachments(changeStatusDTO.Attachments),
	}
}
