package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type ChangeStatusDTO struct {
	Status      string                `json:"status"`
	UpdatedBy   string                `json:"updated_by"`
	Content     *string               `json:"content"`
	Attachments []CreateAttachmentDTO `json:"attachments"`
}

func mapChangeStatusDTOToChangeStatus(changeStatusDTO ChangeStatusDTO) (domain.ChangeStatus, error) {
	attachmentDTO, err := mapCreateAttachmentDTOsToAttachments(changeStatusDTO.Attachments)
	if err != nil {
		return domain.ChangeStatus{}, err
	}

	return domain.ChangeStatus{
		Status:      domain.CaseStatus(changeStatusDTO.Status),
		UpdatedBy:   changeStatusDTO.UpdatedBy,
		Content:     changeStatusDTO.Content,
		Attachments: attachmentDTO,
	}, nil
}
