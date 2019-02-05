package formDto

import (
	"context"
	"net/url"

	"errors"

	"strings"

	"strconv"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	customerDomain "flamingo.me/flamingo-commerce/customer/domain"
	"flamingo.me/flamingo/core/form/application"
	formDomain "flamingo.me/flamingo/core/form/domain"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"github.com/go-playground/form"
	"github.com/leebenson/conform"
	"gopkg.in/go-playground/validator.v9"
)

type (
	CheckoutFormData struct {
		BillingAddress                     AddressFormData `form:"billingAddress"`
		PersonalData                       PersonalData    `form:"personalData"`
		ShippingAddress                    AddressFormData `form:"shippingAddress" validate:"-"`
		UseBillingAddressAsShippingAddress bool            `form:"billingAsShipping"`
		TermsAndConditions                 bool            `form:"termsAndConditions" validate:"required"`
		LoyaltyMemberShipNumber            string          `form:"loyaltyPointsMembershipId"`
		SelectedPaymentProvider            string          `form:"selectedPaymentProvider" validate:"required"`
		SelectedPaymentProviderMethod      string          `form:"selectedPaymentProviderMethod" validate:"required"`
	}

	PersonalData struct {
		DateOfBirth     string          `form:"dateOfBirth"`
		PassportCountry string          `form:"passportCountry"`
		PassportNumber  string          `form:"passportNumber"`
		Address         AddressFormData `form:"address" validate:"-"`
	}

	// AddressFormData defines the checkout address form data
	AddressFormData struct {
		Vat              string `form:"vat"`
		Firstname        string `form:"firstname" validate:"required" conform:"name"`
		Lastname         string `form:"lastname" validate:"required" conform:"name"`
		MiddleName       string `form:"middlename" conform:"name"`
		Title            string `form:"title" conform:"trim"`
		Salutation       string `form:"salutation" conform:"trim"`
		Street           string `form:"street" conform:"trim"`
		StreetNr         string `form:"streetNr" conform:"trim"`
		AddressLine1     string `form:"addressLine1" conform:"trim"`
		AddressLine2     string `form:"addressLine2" conform:"trim"`
		Company          string `form:"company" conform:"trim"`
		PostCode         string `form:"postCode" conform:"trim"`
		City             string `form:"city" conform:"name"`
		State            string `form:"state" conform:"trim"`
		RegionCode       string `form:"regionCode" conform:"name"`
		Country          string `form:"country" conform:"trim"`
		CountryCode      string `form:"countryCode" conform:"name"`
		PhoneAreaCode    string `form:"phoneAreaCode"`
		PhoneCountryCode string `form:"phoneCountryCode"`
		PhoneNumber      string `form:"phoneNumber" conform:"num"`
		Email            string `form:"email" validate:"required,email" conform:"email"`
	}

	CheckoutFormService struct {
		DefaultFormValues    config.Map    `inject:"config:checkout.checkoutForm.defaultValues,optional"`
		OverrideFormValues   config.Map    `inject:"config:checkout.checkoutForm.overrideValues,optional"`
		AdditionalFormValues config.Slice  `inject:"config:checkout.checkoutForm.additionalFormValues,optional"`
		Decoder              *form.Decoder `inject:""`
		//Customer  might be passed by the controller - we use it to initialize the form
		Customer customerDomain.Customer
		//Cart might be passed by Controller - we use it to prefill the form in case it was not submitted
		Cart *cart.Cart

		Logger flamingo.Logger `inject:""`

		//A couple of configuration options for more flexible Validation
		PersonalData_DateOfBirthRequired     bool    `inject:"config:checkout.checkoutForm.validation.personalData.dateOfBirthRequired,optional"`
		PersonalData_MinAge                  float64 `inject:"config:checkout.checkoutForm.validation.personalData.minAge,optional"`
		PersonalData_PassportCountryRequired bool    `inject:"config:checkout.checkoutForm.validation.personalData.passportCountryRequired,optional"`
		PersonalData_PassportNumberRequired  bool    `inject:"config:checkout.checkoutForm.validation.personalData.passportNumberRequired,optional"`
		BillingAddress_PhoneNumberRequired   bool    `inject:"config:checkout.checkoutForm.validation.billingAddress.phoneNumberRequired,optional"`
	}
)

var _ formDomain.FormService = new(CheckoutFormService)
var _ formDomain.GetDefaultFormData = new(CheckoutFormService)

// ParseFormData - from FormService interface
// MEthod is Responsible to parse the Post Values and fill the FormData struct
func (fs *CheckoutFormService) ParseFormData(ctx context.Context, r *web.Request, formValues url.Values) (interface{}, error) {
	if formValues == nil {
		formValues = make(map[string][]string)
	}

	// Preset eMail when email parameter is given:
	if ctx != nil {
		email, ok := r.Form("email")
		if ok && len(formValues["billingAddress.email"]) == 0 {
			formValues["billingAddress.email"] = email
		}
	}

	fs.Logger.WithField("category", "checkout").Debug("passed formValues before modifications: %#v", formValues)

	//Merge in configured DefaultValues that are configured
	formValues = fs.setConfiguredDefaultFormValues(formValues)

	//OverrideValues
	formValues = fs.overrideConfiguredDefaultFormValues(formValues)

	fs.Logger.WithField("category", "checkout").Debug("formValues after modifications: %#v", formValues)

	var formData CheckoutFormData
	fs.Decoder.Decode(&formData, formValues)

	conform.Strings(&formData)

	return formData, nil
}

//GetDefaultFormData - from interface GetDefaultFormData
// Is called if form is not submitted - to get FormData with default values
func (fs *CheckoutFormService) GetDefaultFormData(parsedData interface{}) interface{} {
	if checkoutFormData, ok := parsedData.(CheckoutFormData); ok {
		checkoutFormData = fs.fillFormDataFromCart(checkoutFormData)
		checkoutFormData = fs.fillFormDataFromCustomer(checkoutFormData)
		return checkoutFormData
	}
	return parsedData

}

func (fs *CheckoutFormService) fillFormDataFromCustomer(formData CheckoutFormData) CheckoutFormData {
	//If customer is given - get default values for the form if not empty yet
	if fs.Customer != nil {
		billingAddress := fs.Customer.GetDefaultBillingAddress()
		if billingAddress != nil {
			fs.mapCustomerAddressToFormAddress(*billingAddress, &formData.BillingAddress)
		}
		shippingAddress := fs.Customer.GetDefaultShippingAddress()
		if shippingAddress != nil {
			fs.mapCustomerAddressToFormAddress(*shippingAddress, &formData.ShippingAddress)
		}
		if !fs.Customer.GetPersonalData().Birthday.IsZero() {
			formData.PersonalData.DateOfBirth = fs.Customer.GetPersonalData().Birthday.Format("2006-01-02")
		}
	}
	return formData
}

func (fs *CheckoutFormService) fillFormDataFromCart(formData CheckoutFormData) CheckoutFormData {
	if fs.Cart != nil {
		fs.mapCartAddressToFormAddress(fs.Cart.BillingAdress, &formData.BillingAddress)
		if len(fs.Cart.Deliveries) > 0 {
			if fs.Cart.Deliveries[0].DeliveryInfo.DeliveryLocation.Address != nil {
				fs.mapCartAddressToFormAddress(*fs.Cart.Deliveries[0].DeliveryInfo.DeliveryLocation.Address, &formData.ShippingAddress)
			}
		}

		formData.PersonalData.DateOfBirth = fs.Cart.Purchaser.PersonalDetails.DateOfBirth
		formData.PersonalData.PassportNumber = fs.Cart.Purchaser.PersonalDetails.PassportNumber
		formData.PersonalData.PassportCountry = fs.Cart.Purchaser.PersonalDetails.PassportCountry
		if fs.Cart.Purchaser.Address != nil {
			fs.mapCartAddressToFormAddress(*fs.Cart.Purchaser.Address, &formData.PersonalData.Address)
		}
	}
	return formData
}

func (fs *CheckoutFormService) mapCustomerAddressToFormAddress(address customerDomain.Address, targetAddress *AddressFormData) {
	if targetAddress.Email == "" {
		targetAddress.Email = address.Email
	}
	if targetAddress.Firstname == "" {
		targetAddress.Firstname = address.Firstname
	}
	if targetAddress.Lastname == "" {
		targetAddress.Lastname = address.Lastname
	}
	if targetAddress.CountryCode == "" {
		targetAddress.CountryCode = address.CountryCode
	}
	if targetAddress.PhoneNumber == "" {
		targetAddress.PhoneNumber = address.Telephone
	}

	if targetAddress.Street == "" && targetAddress.City == "" {
		targetAddress.Street = address.Street
		targetAddress.StreetNr = address.StreetNr
		targetAddress.City = address.City
	}
}

func (fs *CheckoutFormService) mapCartAddressToFormAddress(address cart.Address, targetAddress *AddressFormData) {
	if address.Firstname != "" {
		targetAddress.Firstname = address.Firstname
	}

	if address.Lastname != "" {
		targetAddress.Lastname = address.Lastname
	}

	if address.Email != "" {
		targetAddress.Email = address.Email
	}

	if address.Street != "" {
		targetAddress.Street = address.Street
	}

	if address.StreetNr != "" {
		targetAddress.StreetNr = address.StreetNr
	}

	if address.Salutation != "" {
		targetAddress.Title = address.Salutation
	}

	if address.City != "" {
		targetAddress.City = address.City
	}

	if address.Telephone != "" {
		targetAddress.PhoneNumber = address.Telephone
	}

	if address.CountryCode != "" {
		targetAddress.CountryCode = address.CountryCode
	}

	if address.Company != "" {
		targetAddress.Company = address.Company
	}
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

func (fs CheckoutFormService) GetAdditionFormFields(formData CheckoutFormData) map[string]string {
	additionalFormData := make(map[string]string)

	if fs.AdditionalFormValues != nil {
		for _, key := range fs.AdditionalFormValues {
			if key == "lp_membership_id" && formData.LoyaltyMemberShipNumber != "" {
				additionalFormData[key.(string)] = formData.LoyaltyMemberShipNumber
			}
		}
	}

	return additionalFormData
}

//ValidateFormData - from FormService interface
func (fs *CheckoutFormService) ValidateFormData(data interface{}) (formDomain.ValidationInfo, error) {
	formData, ok := data.(CheckoutFormData)
	if !ok {
		return formDomain.ValidationInfo{}, errors.New("Cannot convert to AddressFormData")
	}
	validate := validator.New()
	validationInfo := application.ValidationErrorsToValidationInfo(validate.Struct(formData))

	if fs.PersonalData_DateOfBirthRequired {
		if formData.PersonalData.DateOfBirth == "" {
			validationInfo.AddFieldError("personalData.dateOfBirth", "formerror_dateOfBirth_required", "date of birth is required")
		} else if !formDomain.ValidateDate(formData.PersonalData.DateOfBirth) {
			validationInfo.AddFieldError("personalData.dateOfBirth", "formerror_dateOfBirth_formaterror", "date of birth has wrong format required: yyyy-mm-dd")
		} else if fs.PersonalData_MinAge > 0 {
			if !formDomain.ValidateAge(formData.PersonalData.DateOfBirth, int(fs.PersonalData_MinAge)) {
				validationInfo.AddFieldError("personalData.dateOfBirth", "formerror_dateOfBirth_tooyoung", "you need to be at least "+strconv.Itoa(int(fs.PersonalData_MinAge)))
			}
		}
	}

	if fs.BillingAddress_PhoneNumberRequired {
		if formData.BillingAddress.PhoneNumber == "" {
			validationInfo.AddFieldError("billingAddress.phoneNumber", "formerror_phoneNumber_required", "phone number is required")
		}
	}
	if fs.PersonalData_PassportCountryRequired {
		if formData.PersonalData.PassportCountry == "" {
			validationInfo.AddFieldError("personalData.passportCountry", "formerror_passportCountry_required", "passport infos are required")
		}
	}
	if fs.PersonalData_PassportNumberRequired {
		if formData.PersonalData.PassportNumber == "" {
			validationInfo.AddFieldError("personalData.passportNumber", "formerror_passportNumber_required", "passport infos are required")
		}
	}

	return validationInfo, nil
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

func MapPerson(data CheckoutFormData) *cart.Person {
	if data.PersonalData.IsEmpty() {
		return nil
	}
	address := mapAddress(data.PersonalData.Address)
	person := cart.Person{
		PersonalDetails: cart.PersonalDetails{
			PassportNumber:  data.PersonalData.PassportNumber,
			PassportCountry: data.PersonalData.PassportCountry,
			DateOfBirth:     data.PersonalData.DateOfBirth,
		},
		Address: address,
	}
	return &person
}

func mapAddress(addressData AddressFormData) *cart.Address {

	lines := make([]string, 2)
	lines[0] = addressData.AddressLine1
	lines[1] = addressData.AddressLine2

	address := cart.Address{
		Vat:                    addressData.Vat,
		Firstname:              addressData.Firstname,
		Lastname:               addressData.Lastname,
		MiddleName:             addressData.MiddleName,
		Title:                  addressData.Title,
		Salutation:             addressData.Salutation,
		Street:                 addressData.Street,
		StreetNr:               addressData.StreetNr,
		AdditionalAddressLines: lines,
		Company:                addressData.Company,
		PostCode:               addressData.PostCode,
		City:                   addressData.City,
		State:                  addressData.State,
		RegionCode:             addressData.RegionCode,
		Country:                addressData.Country,
		CountryCode:            addressData.CountryCode,
		Email:                  addressData.Email,
		Telephone:              addressData.PhoneCountryCode + addressData.PhoneAreaCode + addressData.PhoneNumber,
	}
	return &address
}

func (p PersonalData) IsEmpty() bool {
	if p.PassportNumber == "" && p.PassportCountry == "" && p.DateOfBirth == "" {
		return true
	}
	return false
}
