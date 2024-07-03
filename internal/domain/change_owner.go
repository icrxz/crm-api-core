package domain

type ChangeOwner struct {
	OwnerID   string
	UpdatedBy string
	Status    CaseStatus
}
