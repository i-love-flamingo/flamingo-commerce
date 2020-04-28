package domain

import (
	"context"
	"errors"
	"time"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	// Customer data interface
	Customer interface {
		GetID() string
		GetPersonalData() PersonData
		GetAddresses() []Address
		GetDefaultShippingAddress() *Address
		GetDefaultBillingAddress() *Address
	}

	// PersonData contains personal data
	PersonData struct {
		// Gender male, female, other, unknown
		Gender     string
		FirstName  string
		LastName   string
		MiddleName string
		MainEmail  string
		// Prefix such as (mr, mrs, ms) string
		Prefix      string
		Birthday    time.Time
		Nationality string
	}

	// Address data of a customer
	Address struct {
		RegionCode             string
		CountryCode            string
		Company                string
		Street                 string
		StreetNr               string
		AdditionalAddressLines []string
		Telephone              string
		PostCode               string
		City                   string
		Firstname              string
		Lastname               string
		Email                  string
	}

	// CustomerService to retrieve customers
	// Deprecated: use CustomerIdentityService and core/auth module instead
	CustomerService interface {
		// GetByAuth - returns Customer by the provided Auth infos
		GetByAuth(ctx context.Context, auth domain.Auth) (Customer, error)
	}

	// CustomerIdentityService to retrieve customers by Identity
	CustomerIdentityService interface {
		GetByIdentity(ctx context.Context, identity auth.Identity) (Customer, error)
	}
)

var (
	// ErrCustomerNotFoundError - semantic error returned if no customer was found
	ErrCustomerNotFoundError = errors.New("Customer not found")
)

const (
	// GenderMale for male customers
	GenderMale = "male"
	// GenderFemale for female customers
	GenderFemale = "female"
	// GenderOther for other customers
	GenderOther = "other"
	// GenderUnknown unknown
	GenderUnknown = ""
)
