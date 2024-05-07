package domain

type Address struct {
	Address1 string
	Address2 string
	City     string
	State    string
	ZipCode  string
	Region   int
}

func NewAddress(address1, address2, city, state, zipCode string, region int) Address {
	return Address{
		Address1: address1,
		Address2: address2,
		City:     city,
		State:    state,
		ZipCode:  zipCode,
		Region:   region,
	}
}
