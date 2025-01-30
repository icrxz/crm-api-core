package domain

type Contact struct {
	PhoneNumber string
	Email       string
}

type UpdateContact struct {
	PhoneNumber *string
	Email       *string
}

func (c *Contact) MergeUpdate(updateContact UpdateContact) {
	if updateContact.PhoneNumber != nil {
		c.PhoneNumber = *updateContact.PhoneNumber
	}

	if updateContact.Email != nil {
		c.Email = *updateContact.Email
	}
}
