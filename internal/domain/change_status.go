package domain

type ChangeStatus struct {
	Status      CaseStatus
	UpdatedBy   string
	Content     *string
	Attachments []Attachment
}
