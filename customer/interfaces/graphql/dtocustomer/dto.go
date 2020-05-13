package dtocustomer

import (
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
