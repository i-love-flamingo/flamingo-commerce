package domain

import (
	"time"

	"go.aoe.com/flamingo/core/auth/domain"
)

type (
	Customer interface {
		GetPersonalData() PersonalData
		GetAddresses() []Address
		GetDefaultShippingAddress() *Address
		GetDefaultBillingAddress() *Address
	}

	PersonalData struct {
		//Gender male, female, other, unknown
		Gender     string
		FirstName  string
		LastName   string
		MiddleName string
		MainEmail  string
		//Prefix such as (mr, mrs, ms) string
		Prefix      string
		Birthday    time.Time
		Nationality string
	}

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

	CustomerService interface {
		//GetByAuth - returns Customer by the provided Auth infos
		GetByAuth(auth domain.Auth) (Customer, error)
	}
)

const (
	GENDER_MALE    = "male"
	GENDER_FEMALE  = "female"
	GENDER_OTHER   = "other"
	GENDER_UNKNOWN = ""
)
