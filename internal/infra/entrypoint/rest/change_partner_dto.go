package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type ChangePartnerDTO struct {
	PartnerID  string            `json:"partner_id"`
	TargetDate time.Time         `json:"target_date"`
	Status     domain.CaseStatus `json:"status"`
	UpdatedBy  string            `json:"updated_by"`
}

func mapChangePartnerDTOToChangePartner(c ChangePartnerDTO) domain.ChangePartner {
	return domain.ChangePartner{
		PartnerID:  c.PartnerID,
		TargetDate: c.TargetDate,
		Status:     c.Status,
		UpdatedBy:  c.UpdatedBy,
	}
}
