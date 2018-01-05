package formDto

import (
	"net/url"

	"errors"

	"strings"

	"github.com/go-playground/form"
	"github.com/leebenson/conform"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	domain2 "go.aoe.com/flamingo/core/customer/domain"
	"go.aoe.com/flamingo/core/form/application"
	"go.aoe.com/flamingo/core/form/domain"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/web"
	"gopkg.in/go-playground/validator.v9"
)

type (
	CheckoutFormData struct {
		BillingAddress                     AddressFormData `form:"billingAddress"`
		ShippingAddress                    AddressFormData `form:"shippingAddress" validate:"-"`
		UseBillingAddressAsShippingAddress bool            `form:"billingAsShipping"`
		TermsAndConditions                 bool            `form:"termsAndConditions" validate:"required"`
	}

	AddressFormData struct {
		RegionCode    string `form:"regionCode" conform:"name"`
		CountryCode   string `form:"countryCode" conform:"name"`
		Company       string `form:"company" conform:"trim"`
		Street        string `form:"street" conform:"trim"`
		StreetNr      string `form:"streetNr" conform:"trim"`
		AddressLine1  string `form:"addressLine1" conform:"trim"`
		AddressLine2  string `form:"addressLine2" conform:"trim"`
		PhoneAreaCode string `form:"phoneAreaCode" conform:"num"`
		PhoneNumber   string `form:"phoneNumber"  conform:"num"`
		PostCode      string `form:"postCode" conform:"trim"`
		City          string `form:"city" conform:"name"`
		Firstname     string `form:"firstname" validate:"required" conform:"name"`
		Lastname      string `form:"lastname" validate:"required" conform:"name"`
		Email         string `form:"email" validate:"required,email" conform:"email"`
	}

	CheckoutFormService struct {
		DefaultFormValues  config.Map    `inject:"config:checkout.checkoutForm.defaultValues,optional"`
		OverrideFormValues config.Map    `inject:"config:checkout.checkoutForm.overrideValues,optional"`
		Decoder            *form.Decoder `inject:""`
		//Customer  might be passed by the controller - we use it to initialize the form
		Customer domain2.Customer
	}
)

// ParseFormData - from FormService interface
func (fs *CheckoutFormService) ParseFormData(ctx web.Context, formValues url.Values) (interface{}, error) {
	if formValues == nil {
		formValues = make(map[string][]string)
	}

	// Preset eMail when email parameter is given:
	email, e := ctx.Query("email")
	if e == nil {
		formValues["billingAddress.email"] = email
	}

	//Merge in DefaultValues
	if fs.DefaultFormValues != nil {
		for k, v := range fs.DefaultFormValues {
			k = strings.Replace(k, "_", ".", -1)
			if _, ok := formValues[k]; ok {
				//value is passed - no override
				continue
			}
			stringV, ok := v.(string)
			if !ok {
				//value configured is no string - missconfiguration - continue
				continue
			}
			newStringSlice := make([]string, 1)
			newStringSlice[0] = stringV
			formValues[k] = newStringSlice
		}
	}

	//log.Printf("formValues before %#v", formValues)
	//log.Printf("fs.OverrideFormValues %#v", fs.OverrideFormValues)
	//OverrideValues
	if fs.OverrideFormValues != nil {
		for k, v := range fs.OverrideFormValues {
			k = strings.Replace(k, "_", ".", -1)
			stringV, ok := v.(string)
			if !ok {
				//value configured is no string - missconfiguration - continue
				continue
			}
			newStringSlice := make([]string, 1)
			newStringSlice[0] = stringV
			formValues[k] = newStringSlice
		}
	}
	//log.Printf("formValues %#v", formValues)
	var formData CheckoutFormData
	fs.Decoder.Decode(&formData, formValues)
	//If customer is given - get default values for the form if not empty yet
	if fs.Customer != nil {
		fillAddressFormWithCustomerAddress(&formData.BillingAddress, fs.Customer.GetDefaultBillingAddress())
		fillAddressFormWithCustomerAddress(&formData.ShippingAddress, fs.Customer.GetDefaultShippingAddress())
	}
	conform.Strings(&formData)
	return formData, nil
}

func fillAddressFormWithCustomerAddress(addressForm *AddressFormData, address *domain2.Address) {
	if address == nil || addressForm == nil {
		return
	}
	if addressForm.Email == "" {
		addressForm.Email = address.Email
	}
	if addressForm.Firstname == "" {
		addressForm.Firstname = address.Firstname
	}
	if addressForm.Lastname == "" {
		addressForm.Lastname = address.Lastname
	}
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
	if data.UseBillingAddressAsShippingAddress {
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
