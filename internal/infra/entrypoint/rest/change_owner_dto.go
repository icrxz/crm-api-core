package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type ChangeOwnerDTO struct {
	OwnerID   string            `json:"owner_id"`
	UpdatedBy string            `json:"updated_by"`
	Status    domain.CaseStatus `json:"status"`
}

func mapChangeOwnerDTOToChangeOwner(changeOwnerDTO ChangeOwnerDTO) domain.ChangeOwner {
	return domain.ChangeOwner{
		OwnerID:   changeOwnerDTO.OwnerID,
		UpdatedBy: changeOwnerDTO.UpdatedBy,
		Status:    changeOwnerDTO.Status,
	}
}
