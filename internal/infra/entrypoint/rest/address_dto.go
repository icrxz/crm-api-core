package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type AddressDTO struct {
	Address string `json:"address"`
	State   string `json:"state"`
	City    string `json:"city"`
	Country string `json:"country"`
	ZipCode string `json:"zip_code"`
}

func mapAddressDTOToAddress(addressDTO AddressDTO) domain.Address {
	return domain.Address{
		Address: addressDTO.Address,
		State:   addressDTO.State,
		City:    addressDTO.City,
		Country: addressDTO.Country,
		ZipCode: addressDTO.ZipCode,
	}
}

func mapAddressToAddressDTO(address domain.Address) AddressDTO {
	return AddressDTO{
		Address: address.Address,
		State:   address.State,
		City:    address.City,
		Country: address.Country,
		ZipCode: address.ZipCode,
	}
}
