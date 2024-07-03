package domain

import "time"

type ChangePartner struct {
	PartnerID  string
	TargetDate time.Time
	Status     CaseStatus
	UpdatedBy  string
}
