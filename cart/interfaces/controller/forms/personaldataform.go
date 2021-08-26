package forms

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/application"
	"flamingo.me/form/domain"
	"github.com/go-playground/form/v4"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	customerApplication "flamingo.me/flamingo-commerce/v3/customer/application"
)

type (
	// PersonalDataForm - interface for the form DTO
	PersonalDataForm interface {
		MapPerson() *cart.Person
		MapAdditionalData() *cart.AdditionalData
	}

	// DefaultPersonalDataFormProvider for creating form instances
	DefaultPersonalDataFormProvider func() *DefaultPersonalDataForm

	// DefaultPersonalDataForm - the standard form dto for Person data (that implements PersonalDataForm)
	DefaultPersonalDataForm struct {
		DateOfBirth             string      `form:"dateOfBirth"`
		PassportCountry         string      `form:"passportCountry"`
		PassportNumber          string      `form:"passportNumber"`
		Address                 AddressForm `form:"address" validate:"-"`
		additionalFormData      map[string]string
		additionalFormFieldsCfg config.Slice
	}

	// DefaultPersonalDataFormService implements Form(Data)Provider interface of form package
	DefaultPersonalDataFormService struct {
		applicationCartReceiverService  *cartApplication.CartReceiverService
		defaultPersonalDataFormProvider DefaultPersonalDataFormProvider
		customerApplicationService      *customerApplication.Service
		additionalFormFieldsCfg         config.Slice
		dateOfBirthRequired             bool
		minAge                          int
		passportCountryRequired         bool
		passportNumberRequired          bool
	}

	// PersonalDataFormController - the (mini) MVC for handling Personal Data (Purchaser)
	PersonalDataFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService

		logger flamingo.Logger

		formHandlerFactory application.FormHandlerFactory
	}
)

// Inject - dependencies
func (p *DefaultPersonalDataFormService) Inject(
	applicationCartReceiverService *cartApplication.CartReceiverService,
	defaultPersonalDataFormProvider DefaultPersonalDataFormProvider,
	customerApplicationService *customerApplication.Service,
	cfg *struct {
		AdditionalFormValues    config.Slice `inject:"config:commerce.cart.personalDataForm.additionalFormFields,optional"`
		DateOfBirthRequired     bool         `inject:"config:commerce.cart.personalDataForm.dateOfBirthRequired,optional"`
		MinAge                  float64      `inject:"config:commerce.cart.personalDataForm.minAge,optional"`
		PassportCountryRequired bool         `inject:"config:commerce.cart.personalDataForm.passportCountryRequired,optional"`
		PassportNumberRequired  bool         `inject:"config:commerce.cart.personalDataForm.passportNumberRequired,optional"`
	},
) *DefaultPersonalDataFormService {
	p.applicationCartReceiverService = applicationCartReceiverService
	p.defaultPersonalDataFormProvider = defaultPersonalDataFormProvider
	p.customerApplicationService = customerApplicationService
	if cfg != nil {
		p.additionalFormFieldsCfg = cfg.AdditionalFormValues
		p.dateOfBirthRequired = cfg.DateOfBirthRequired
		p.minAge = int(cfg.MinAge)
		p.passportCountryRequired = cfg.PassportCountryRequired
		p.passportNumberRequired = cfg.PassportNumberRequired
	}

	return p
}

// GetFormData from data provider
func (p *DefaultPersonalDataFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
	formData := p.defaultPersonalDataFormProvider()

	cart, err := p.applicationCartReceiverService.ViewCart(ctx, req.Session())
	if err != nil {
		return *formData, nil
	}

	customer, err := p.customerApplicationService.GetForIdentity(ctx, req)
	if err == nil {
		personalData := customer.GetPersonalData()
		formData.DateOfBirth = personalData.Birthday.Format("2006-01-02")
	}

	if cart.Purchaser != nil {
		formData.DateOfBirth = cart.Purchaser.PersonalDetails.DateOfBirth
		formData.PassportCountry = cart.Purchaser.PersonalDetails.PassportCountry
		formData.PassportNumber = cart.Purchaser.PersonalDetails.PassportNumber
		if cart.Purchaser.Address != nil {
			formData.Address.LoadFromCartAddress(*cart.Purchaser.Address)
		}
	}

	if p.additionalFormFieldsCfg != nil {
		for _, key := range p.additionalFormFieldsCfg {
			formData.additionalFormData[key.(string)] = cart.AdditionalData.CustomAttributes[key.(string)]
		}
	}

	return *formData, nil
}

// Decode fills the form data from the web request
func (p *DefaultPersonalDataFormService) Decode(_ context.Context, _ *web.Request, values url.Values, _ interface{}) (interface{}, error) {
	decoder := form.NewDecoder()
	personalDataForm := p.defaultPersonalDataFormProvider()
	err := decoder.Decode(personalDataForm, values)
	if err != nil {
		return nil, err
	}

	if p.additionalFormFieldsCfg != nil {
		for _, key := range p.additionalFormFieldsCfg {
			personalDataForm.additionalFormData[key.(string)] = values.Get(key.(string))
		}
	}

	return personalDataForm, nil
}

// Validate form data
func (p *DefaultPersonalDataFormService) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
	personalDataForm, ok := formData.(*DefaultPersonalDataForm)
	if !ok {
		return nil, errors.New("no DefaultPersonalDataForm given")
	}
	validationInfo := domain.ValidationInfo{}
	validationInfo = validatorProvider.Validate(ctx, req, personalDataForm)

	if p.dateOfBirthRequired && personalDataForm.DateOfBirth == "" {
		validationInfo.AddFieldError("personalData.dateOfBirth", "formerror_dateOfBirth_required", "date of birth is required")
	}

	if personalDataForm.DateOfBirth != "" && !validateDOB(personalDataForm.DateOfBirth) {
		validationInfo.AddFieldError("personalData.dateOfBirth", "formerror_dateOfBirth_formaterror", "date of birth has wrong format required: yyyy-mm-dd")
	}

	if p.minAge > 0 {
		if !validateMinimumAge(personalDataForm.DateOfBirth, p.minAge) {
			validationInfo.AddFieldError("personalData.dateOfBirth", "formerror_dateOfBirth_tooyoung", "you need to be at least "+strconv.Itoa(p.minAge))
		}
	}

	if p.passportCountryRequired && personalDataForm.PassportCountry == "" {
		validationInfo.AddFieldError("personalData.passportCountry", "formerror_passportCountry_required", "passport infos are required")
	}

	if p.passportNumberRequired && personalDataForm.PassportNumber == "" {
		validationInfo.AddFieldError("personalData.passportNumber", "formerror_passportNumber_required", "passport infos are required")
	}

	return &validationInfo, nil
}

func validateDOB(value string) bool {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}

	return date.Before(time.Now().Add(24 * time.Hour))
}

func validateMinimumAge(value string, minimumAge int) bool {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}

	now := time.Now()
	required := time.Date(
		now.Year()-minimumAge,
		now.Month(),
		now.Day(),
		0,
		0,
		0,
		0,
		date.Location(),
	)

	return date.Add(-time.Minute).Before(required)
}

// Inject dependencies
func (c *PersonalDataFormController) Inject(
	responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	logger flamingo.Logger,
	formHandlerFactory application.FormHandlerFactory,
) *PersonalDataFormController {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService

	c.formHandlerFactory = formHandlerFactory
	c.logger = logger

	return c
}

// GetUnsubmittedForm returns a Unsubmitted form - using the registered FormService
func (c *PersonalDataFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*domain.Form, error) {
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, err
	}
	return formHandler.HandleUnsubmittedForm(ctx, r)
}

// HandleFormAction handles post of personal data and updates cart
func (c *PersonalDataFormController) HandleFormAction(ctx context.Context, r *web.Request) (form *domain.Form, actionSuccessFull bool, err error) {
	session := web.SessionFromContext(ctx)
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, false, err
	}
	// ##  Handle the submitted form (validation etc)
	form, err = formHandler.HandleSubmittedForm(ctx, r)
	if err != nil {
		return nil, false, err
	}
	personalDataForm, ok := form.Data.(PersonalDataForm)
	if !ok {
		return form, false, errors.New("cannot convert to PersonalDataForm ")
	}
	if !form.IsValidAndSubmitted() {
		return form, false, nil
	}

	// UpdatePurchaser (and add additional data)
	err = c.applicationCartService.UpdatePurchaser(ctx, session, personalDataForm.MapPerson(), personalDataForm.MapAdditionalData())
	if err != nil {
		c.logger.WithContext(ctx).Error("PersonalDataFormController UpdatePurchaser Error %v", err)
		return form, false, err
	}
	return form, true, nil
}

func (c *PersonalDataFormController) getFormHandler() (domain.FormHandler, error) {
	builder := c.formHandlerFactory.GetFormHandlerBuilder()
	err := builder.SetNamedFormService("commerce.cart.personaldataFormService")
	if err != nil {
		return nil, err
	}
	return builder.Build(), nil
}

// Inject dependencies
func (p *DefaultPersonalDataForm) Inject(
	cfg *struct {
		AdditionalFormValues config.Slice `inject:"config:commerce.cart.personalDataForm.additionalFormFields,optional"`
	},
) *DefaultPersonalDataForm {
	if cfg != nil {
		p.additionalFormFieldsCfg = cfg.AdditionalFormValues
	}

	p.additionalFormData = make(map[string]string)

	return p
}

// MapPerson maps the checkout form data to the cart.Person domain struct
func (p *DefaultPersonalDataForm) MapPerson() *cart.Person {
	person := cart.Person{
		PersonalDetails: cart.PersonalDetails{
			PassportNumber:  p.PassportNumber,
			PassportCountry: p.PassportCountry,
			DateOfBirth:     p.DateOfBirth,
		},
	}
	return &person
}

// AdditionalData returns the additional form data
func (p DefaultPersonalDataForm) AdditionalData(key string) string {
	return p.additionalFormData[key]
}

// MapAdditionalData - mapper
func (p DefaultPersonalDataForm) MapAdditionalData() *cart.AdditionalData {
	return &cart.AdditionalData{
		CustomAttributes: p.additionalFormData,
	}
}
