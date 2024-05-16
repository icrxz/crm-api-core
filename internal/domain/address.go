package domain

type Address struct {
	Address string
	City    string
	State   string
	ZipCode string
	Country string
	Region  int
}

func NewAddress(address1, address2, city, state, zipCode, country string, region int) Address {
	return Address{
		Address: address1,
		City:    city,
		State:   state,
		ZipCode: zipCode,
		Country: country,
		Region:  region,
	}
}
