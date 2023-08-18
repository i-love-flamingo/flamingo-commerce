package forms

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/customer/domain"
)

type (
	// AddressForm defines the checkout address form data
	AddressForm struct {
		Vat              string `form:"vat"`
		Firstname        string `form:"firstname" validate:"required" conform:"ucfirst,trim"`
		Lastname         string `form:"lastname" validate:"required" conform:"ucfirst,trim"`
		MiddleName       string `form:"middlename" conform:"ucfirst,trim"`
		Title            string `form:"title" conform:"trim"`
		Salutation       string `form:"salutation" conform:"trim"`
		Street           string `form:"street" conform:"ucfirst,trim"`
		StreetNr         string `form:"streetNr" conform:"trim"`
		AddressLine1     string `form:"addressLine1" conform:"trim"`
		AddressLine2     string `form:"addressLine2" conform:"trim"`
		Company          string `form:"company" conform:"trim"`
		PostCode         string `form:"postCode" conform:"trim"`
		City             string `form:"city" conform:"ucfirst,trim"`
		State            string `form:"state" conform:"ucfirst,trim"`
		RegionCode       string `form:"regionCode" conform:"trim"`
		Country          string `form:"country" conform:"ucfirst,trim"`
		CountryCode      string `form:"countryCode" conform:"trim"`
		PhoneAreaCode    string `form:"phoneAreaCode" conform:"trim"`
		PhoneCountryCode string `form:"phoneCountryCode" conform:"trim"`
		PhoneNumber      string `form:"phoneNumber" conform:"trim"`
		Email            string `form:"email" validate:"required,email" conform:"trim,lowercase"`
	}
)

// MapToDomainAddress - returns the cart Address Object
func (a *AddressForm) MapToDomainAddress() cart.Address {
	lines := make([]string, 2)
	lines[0] = a.AddressLine1
	lines[1] = a.AddressLine2

	return cart.Address{
		Vat:                    a.Vat,
		Firstname:              a.Firstname,
		Lastname:               a.Lastname,
		MiddleName:             a.MiddleName,
		Title:                  a.Title,
		Salutation:             a.Salutation,
		Street:                 a.Street,
		StreetNr:               a.StreetNr,
		AdditionalAddressLines: lines,
		Company:                a.Company,
		PostCode:               a.PostCode,
		City:                   a.City,
		State:                  a.State,
		RegionCode:             a.RegionCode,
		Country:                a.Country,
		CountryCode:            a.CountryCode,
		Email:                  a.Email,
		TelephoneCountryCode:   a.PhoneCountryCode,
		TelephoneAreaCode:      a.PhoneAreaCode,
		TelephoneNumber:        a.PhoneNumber,
		Telephone:              a.PhoneCountryCode + a.PhoneAreaCode + a.PhoneNumber,
	}
}

// LoadFromCustomerAddress - fills the form from data in the address object (from customer module)
func (a *AddressForm) LoadFromCustomerAddress(address domain.Address) {

	if a.Email == "" || a.Email == "@" {
		a.Email = address.Email
	}
	if a.Firstname == "" {
		a.Firstname = address.Firstname
	}
	if a.Lastname == "" {
		a.Lastname = address.Lastname
	}
	if a.CountryCode == "" {
		a.CountryCode = address.CountryCode
	}
	if a.PhoneNumber == "" {
		a.PhoneNumber = address.Telephone
	}

	if a.Street == "" && a.City == "" {
		a.Street = address.Street
		a.StreetNr = address.StreetNr
		a.City = address.City
	}

}

// LoadFromCartAddress - loads the form data from cart address
func (a *AddressForm) LoadFromCartAddress(address cart.Address) {
	if address.Firstname != "" {
		a.Firstname = address.Firstname
	}
	if address.PostCode != "" {
		a.PostCode = address.PostCode
	}

	if address.State != "" {
		a.State = address.State
	}

	if len(address.AdditionalAddressLines) > 0 {
		a.AddressLine1 = address.AdditionalAddressLines[0]
	}
	if len(address.AdditionalAddressLines) > 1 {
		a.AddressLine2 = address.AdditionalAddressLines[1]
	}

	if address.Lastname != "" {
		a.Lastname = address.Lastname
	}

	if address.Email != "" {
		a.Email = address.Email
	}

	if address.Street != "" {
		a.Street = address.Street
	}

	if address.StreetNr != "" {
		a.StreetNr = address.StreetNr
	}

	if address.Title != "" {
		a.Title = address.Title
	}
	if address.Salutation != "" {
		a.Title = address.Salutation
	}

	if address.City != "" {
		a.City = address.City
	}

	//nolint:staticcheck // deprecated since 10.08
	if address.Telephone != "" {
		a.PhoneNumber = address.Telephone
	}

	if address.TelephoneCountryCode != "" {
		a.PhoneCountryCode = address.TelephoneCountryCode
	}

	if address.TelephoneAreaCode != "" {
		a.PhoneAreaCode = address.TelephoneAreaCode
	}

	// Overrides if used new field
	if address.TelephoneNumber != "" {
		a.PhoneNumber = address.TelephoneNumber
	}

	if address.CountryCode != "" {
		a.CountryCode = address.CountryCode
	}

	if address.Company != "" {
		a.Company = address.Company
	}
}
