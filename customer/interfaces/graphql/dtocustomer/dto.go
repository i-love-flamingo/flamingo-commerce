package dtocustomer

import (
	"errors"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
)

type (
	// CustomerStatusResult is a dto
	CustomerStatusResult struct {
		IsLoggedIn bool
		UserID     string
	}

	// CustomerResult is a dto
	CustomerResult struct {
		ID                     string
		PersonalData           domain.PersonData
		Addresses              []domain.Address
		DefaultShippingAddress domain.Address
		DefaultBillingAddress  domain.Address
	}
)

// GetAddress returns address by id
func (cr *CustomerResult) GetAddress(ID string) (*domain.Address, error) {
	for _, address := range cr.Addresses {
		if address.ID == ID {
			return &address, nil
		}
	}
	return nil, errors.New("address not found")
}
