package domain

type Address struct {
	Address string
	City    string
	State   string
	ZipCode string
	Country string
}

type UpdateAddress struct {
	Address *string
	City    *string
	State   *string
	ZipCode *string
	Country *string
}

func (a *Address) MergeUpdate(updateAddress UpdateAddress) {
	if updateAddress.Address != nil {
		a.Address = *updateAddress.Address
	}

	if updateAddress.City != nil {
		a.City = *updateAddress.City
	}

	if updateAddress.State != nil {
		a.State = *updateAddress.State
	}

	if updateAddress.ZipCode != nil {
		a.ZipCode = *updateAddress.ZipCode
	}

	if updateAddress.Country != nil {
		a.Country = *updateAddress.Country
	}
}
