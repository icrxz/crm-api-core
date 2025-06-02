package domain

type ChangeStatus struct {
	Status      CaseStatus
	UpdatedBy   string
	Content     *string
	Type        *string
	Attachments []Attachment
}
