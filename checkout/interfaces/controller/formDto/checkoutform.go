package formDto

import (
	"net/url"

	"errors"

	"github.com/go-playground/form"
	"github.com/leebenson/conform"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/form/application"
	"go.aoe.com/flamingo/core/form/domain"
	"go.aoe.com/flamingo/framework/web"
	"gopkg.in/go-playground/validator.v9"
)

type (
	CheckoutFormData struct {
		BillingAddress                    AddressFormData `form:"billingAddress"`
		ShippingAddress                   AddressFormData `form:"shippingAddress" validate:"-"`
		UseBillingAddressAsShippinAddress bool            `form:"billingAsShipping"`
		TermsAndConditions                bool            `form:"termsAndConditions" validate:"required"`
	}

	AddressFormData struct {
		RegionCode    string `form:"regionCode" conform:"name"`
		CountryCode   string `form:"countryCode" conform:"name"`
		Company       string `form:"company" conform:"name"`
		Street        string `form:"street" conform:"name"`
		StreetNr      string `form:"streetNr" conform:"name"`
		AddressLine1  string `form:"addressLine1" conform:"name"`
		AddressLine2  string `form:"addressLine2" conform:"name"`
		PhoneAreaCode string `form:"phoneAreaCode" conform:"name"`
		PhoneNumber   string `form:"phoneNumber" conform:"name"`
		PostCode      string `form:"postCode" conform:"name"`
		City          string `form:"city" conform:"name"`
		Firstname     string `form:"firstname" validate:"required" conform:"name"`
		Lastname      string `form:"lastname" validate:"required" conform:"name"`
		Email         string `form:"email" validate:"required" conform:"name"`
	}

	CheckoutFormService struct{}
)

// use a single instance of Decoder, it caches struct info
var decoder *form.Decoder

// ParseFormData - from FormService interface
func (fs *CheckoutFormService) ParseFormData(ctx web.Context, formValues url.Values) (interface{}, error) {
	decoder = form.NewDecoder()
	var formData CheckoutFormData
	decoder.Decode(&formData, formValues)
	conform.Strings(&formData)
	return formData, nil
}

//ValidateFormData - from FormService interface
func (fs *CheckoutFormService) ValidateFormData(data interface{}) (domain.ValidationInfo, error) {
	if formData, ok := data.(CheckoutFormData); ok {
		validate := validator.New()
		return application.ValidationErrorsToValidationInfo(validate.Struct(formData)), nil
	} else {
		return domain.ValidationInfo{}, errors.New("Cannot convert to AddressFormData")
	}
}

func MapAddresses(data CheckoutFormData) (billingAddress *cart.Address, shippingAddress *cart.Address) {
	billingAddress = mapAddress(data.BillingAddress)
	if data.UseBillingAddressAsShippinAddress {
		shippingAddress = billingAddress
	} else {
		shippingAddress = mapAddress(data.ShippingAddress)
	}
	return billingAddress, shippingAddress
}

func mapAddress(addressData AddressFormData) *cart.Address {

	lines := make([]string, 2)
	lines[0] = addressData.AddressLine1
	lines[1] = addressData.AddressLine2

	address := cart.Address{
		CountryCode: addressData.CountryCode,
		Company:     addressData.Company,
		Lastname:    addressData.Lastname,
		Firstname:   addressData.Firstname,
		Email:       addressData.Email,
		City:        addressData.City,
		AdditionalAddressLines: lines,
		RegionCode:             addressData.RegionCode,
		Street:                 addressData.Street,
		PostCode:               addressData.PostCode,
		StreetNr:               addressData.StreetNr,
		Telephone:              addressData.PhoneAreaCode + addressData.PhoneNumber,
	}
	return &address
}
