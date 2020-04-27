package customer

import (
	"context"
	"errors"
	"time"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo/v3/core/auth"
	OAuthDomain "flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
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
	_ domain.CustomerService         = new(FakeService)
	_ domain.CustomerIdentityService = new(FakeService)
	_ domain.Customer                = new(customer)

	birthdayOfFlamingoCommerce, _ = time.Parse("2006-01-02", "2019-04-02")
)

func (f FakeService) GetByAuth(_ context.Context, auth OAuthDomain.Auth) (domain.Customer, error) {
	if auth.IDToken != nil {
		return getCustomer(auth.IDToken.Subject, birthdayOfFlamingoCommerce), nil
	}

	return nil, errors.New("not logged in")
}

func (f FakeService) GetByIdentity(_ context.Context, identity auth.Identity) (domain.Customer, error) {
	if identity != nil {
		return getCustomer(identity.Subject(), birthdayOfFlamingoCommerce), nil
	}

	return nil, errors.New("not logged in")
}

func (c customer) GetId() string {
	return c.id
}

func (c customer) GetPersonalData() domain.PersonData {
	return c.personalData
}

func (c customer) GetAddresses() []domain.Address {
	return c.addresses
}

func (c customer) GetDefaultShippingAddress() *domain.Address {
	return c.defaultShippingAddress
}

func (c customer) GetDefaultBillingAddress() *domain.Address {
	return c.defaultBillingAddress
}

func getCustomer(id string, birthday time.Time) *customer {
	return &customer{
		id: id,
		personalData: domain.PersonData{
			Gender:      domain.GenderMale,
			FirstName:   "Flamingo",
			LastName:    "Commerce",
			MiddleName:  "C.",
			MainEmail:   "flamingo@aoe.com",
			Prefix:      "Mr.",
			Birthday:    birthday,
			Nationality: "DE",
		},
		addresses:              nil,
		defaultShippingAddress: nil,
		defaultBillingAddress:  nil,
	}
}
