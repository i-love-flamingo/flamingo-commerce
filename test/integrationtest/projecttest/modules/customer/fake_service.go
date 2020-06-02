package customer

import (
	"context"
	"errors"
	"time"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo/v3/core/auth"
)

type (
	// FakeService returns the hard configured customer with ID from given auth or identity
	FakeService struct{}

	customer struct {
		id                     string
		personalData           domain.PersonData
		addresses              []domain.Address
		defaultShippingAddress *domain.Address
		defaultBillingAddress  *domain.Address
	}
)

var (
	_ domain.CustomerIdentityService = new(FakeService)
	_ domain.Customer                = new(customer)
)

// GetByIdentity retrieves the authenticated customer by Identity
func (f FakeService) GetByIdentity(_ context.Context, identity auth.Identity) (domain.Customer, error) {
	if identity != nil {
		return getCustomer(identity.Subject()), nil
	}

	return nil, errors.New("not logged in")
}

// GetID of the customer
func (c customer) GetID() string {
	return c.id
}

// GetPersonalData of the customer
func (c customer) GetPersonalData() domain.PersonData {
	return c.personalData
}

// GetAddresses of the customer
func (c customer) GetAddresses() []domain.Address {
	return c.addresses
}

// GetDefaultShippingAddress of the customer
func (c customer) GetDefaultShippingAddress() *domain.Address {
	return c.defaultShippingAddress
}

// GetDefaultBillingAddress of the customer
func (c customer) GetDefaultBillingAddress() *domain.Address {
	return c.defaultBillingAddress
}

func getCustomer(id string) *customer {
	birthdayOfFlamingoCommerce, _ := time.Parse("2006-01-02", "2019-04-02")
	return &customer{
		id: id,
		personalData: domain.PersonData{
			Gender:      domain.GenderMale,
			FirstName:   "Flamingo",
			LastName:    "Commerce",
			MiddleName:  "C.",
			MainEmail:   "flamingo@aoe.com",
			Prefix:      "Mr.",
			Birthday:    birthdayOfFlamingoCommerce,
			Nationality: "DE",
		},
		addresses:              nil,
		defaultShippingAddress: nil,
		defaultBillingAddress:  nil,
	}
}
