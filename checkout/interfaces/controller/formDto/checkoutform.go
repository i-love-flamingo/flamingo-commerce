package formDto

import (
	"net/url"

	"errors"

	"strings"

	"github.com/go-playground/form"
	"github.com/leebenson/conform"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	customerDomain "go.aoe.com/flamingo/core/customer/domain"
	"go.aoe.com/flamingo/core/form/application"
	"go.aoe.com/flamingo/core/form/domain"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/flamingo"
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
		Title            string `form:"title"`
		RegionCode       string `form:"regionCode" conform:"name"`
		CountryCode      string `form:"countryCode" conform:"name"`
		Company          string `form:"company" conform:"trim"`
		Street           string `form:"street" conform:"trim"`
		StreetNr         string `form:"streetNr" conform:"trim"`
		AddressLine1     string `form:"addressLine1" conform:"trim"`
		AddressLine2     string `form:"addressLine2" conform:"trim"`
		PhoneAreaCode    string `form:"phoneAreaCode"`
		PhoneCountryCode string `form:"phoneCountryCode"`
		PhoneNumber      string `form:"phoneNumber" conform:"num"`
		PostCode         string `form:"postCode" conform:"trim"`
		City             string `form:"city" conform:"name"`
		Firstname        string `form:"firstname" validate:"required" conform:"name"`
		Lastname         string `form:"lastname" validate:"required" conform:"name"`
		Email            string `form:"email" validate:"required,email" conform:"email"`
	}

	CheckoutFormService struct {
		DefaultFormValues  config.Map    `inject:"config:checkout.checkoutForm.defaultValues,optional"`
		OverrideFormValues config.Map    `inject:"config:checkout.checkoutForm.overrideValues,optional"`
		Decoder            *form.Decoder `inject:""`
		//Customer  might be passed by the controller - we use it to initialize the form
		Customer customerDomain.Customer
		Logger   flamingo.Logger `inject:""`
	}
)

// ParseFormData - from FormService interface
func (fs *CheckoutFormService) ParseFormData(ctx web.Context, formValues url.Values) (interface{}, error) {
	if formValues == nil {
		formValues = make(map[string][]string)
	}

	// Preset eMail when email parameter is given:
	if ctx != nil {
		email, e := ctx.Query("email")
		if e == nil && len(formValues["billingAddress.email"]) == 0 {
			formValues["billingAddress.email"] = email
		}
	}

	fs.Logger.WithField("category", "checkout").Debugf("passed formValues before modifications: %#v", formValues)

	//Merge in DefaultValues that are configured
	formValues = fs.setDefaultFormValuesFromCustomer(formValues)

	//Merge in configured DefaultValues that are configured
	formValues = fs.setConfiguredDefaultFormValues(formValues)

	//OverrideValues
	formValues = fs.overrideConfiguredDefaultFormValues(formValues)

	fs.Logger.WithField("category", "checkout").Debugf("formValues after modifications: %#v", formValues)

	var formData CheckoutFormData
	fs.Decoder.Decode(&formData, formValues)

	conform.Strings(&formData)
	return formData, nil
}

func (fs *CheckoutFormService) setDefaultFormValuesFromCustomer(formValues url.Values) url.Values {
	//If customer is given - get default values for the form if not empty yet
	if fs.Customer != nil {
		formValues["billingAddress.email"] = make([]string, 1)
		formValues["billingAddress.email"][0] = fs.Customer.GetDefaultBillingAddress().Email
		formValues["billingAddress.firstname"] = make([]string, 1)
		formValues["billingAddress.firstname"][0] = fs.Customer.GetDefaultBillingAddress().Firstname
		formValues["billingAddress.lastname"] = make([]string, 1)
		formValues["billingAddress.lastname"][0] = fs.Customer.GetDefaultBillingAddress().Lastname

		formValues["shippingAddress.email"] = make([]string, 1)
		formValues["shippingAddress.email"][0] = fs.Customer.GetDefaultShippingAddress().Email
		formValues["shippingAddress.firstname"] = make([]string, 1)
		formValues["shippingAddress.firstname"][0] = fs.Customer.GetDefaultShippingAddress().Firstname
		formValues["shippingAddress.lastname"] = make([]string, 1)
		formValues["shippingAddress.lastname"][0] = fs.Customer.GetDefaultShippingAddress().Lastname
	}
	return formValues
}

func (fs *CheckoutFormService) setConfiguredDefaultFormValues(formValues url.Values) url.Values {
	if fs.DefaultFormValues != nil {
		for k, v := range fs.DefaultFormValues {
			k = strings.Replace(k, "_", ".", -1)
			if _, ok := formValues[k]; ok {
				//value is passed - dont set default
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
	return formValues
}

func (fs *CheckoutFormService) overrideConfiguredDefaultFormValues(formValues url.Values) url.Values {
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
	return formValues
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

	firstName := addressData.Firstname
	if addressData.Title != "" {
		firstName = addressData.Title + " " + firstName
	}
	address := cart.Address{
		CountryCode: addressData.CountryCode,
		Company:     addressData.Company,
		Lastname:    addressData.Lastname,
		Firstname:   firstName,
		Email:       addressData.Email,
		City:        addressData.City,
		AdditionalAddressLines: lines,
		RegionCode:             addressData.RegionCode,
		Street:                 addressData.Street,
		PostCode:               addressData.PostCode,
		StreetNr:               addressData.StreetNr,
		Telephone:              addressData.PhoneCountryCode + addressData.PhoneAreaCode + addressData.PhoneNumber,
	}
	return &address
}
