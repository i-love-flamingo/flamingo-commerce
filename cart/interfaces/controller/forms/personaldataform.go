package forms

import (
	"context"
	"errors"

	"flamingo.me/form/application"
	"flamingo.me/form/domain"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//PersonalDataForm - interface for the form DTO
	PersonalDataForm interface {
		MapPerson() (*cart.Person)
		MapAdditionalData() (*cart.AdditionalData)
	}

	//DefaultPersonalDataForm - the standard form dto for Person data (that implements PersonalDataForm)
	DefaultPersonalDataForm struct {
		DateOfBirth     string      `form:"dateOfBirth"`
		PassportCountry string      `form:"passportCountry"`
		PassportNumber  string      `form:"passportNumber"`
		Address         AddressForm `form:"address" validate:"-"`
	}

	// DefaultPersonalDataFormService implements Form(Data)Provider interface of form package
	DefaultPersonalDataFormService struct {
		applicationCartReceiverService *cartApplication.CartReceiverService
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
	applicationCartReceiverService *cartApplication.CartReceiverService) {
	p.applicationCartReceiverService = applicationCartReceiverService
}

// GetFormData from data provider
func (p *DefaultPersonalDataFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
	cart, err := p.applicationCartReceiverService.ViewCart(ctx, req.Session())
	if err == nil {
		if cart.Purchaser != nil {
			formData := DefaultPersonalDataForm{
				DateOfBirth:     cart.Purchaser.PersonalDetails.DateOfBirth,
				PassportCountry: cart.Purchaser.PersonalDetails.PassportCountry,
				PassportNumber:  cart.Purchaser.PersonalDetails.PassportNumber,
			}
			if cart.Purchaser.Address != nil {
				formData.Address.LoadFromCartAddress(*cart.Purchaser.Address)
			}
			return formData, nil
		}
	}
	return DefaultPersonalDataForm{}, nil
}

//Inject - Inject
func (c *PersonalDataFormController) Inject(responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	logger flamingo.Logger,
	formHandlerFactory application.FormHandlerFactory) {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService

	c.formHandlerFactory = formHandlerFactory
	c.logger = logger
}

//GetUnsubmittedForm - returns a Unsubmitted form - using the registered FormService
func (c *PersonalDataFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*domain.Form, error) {
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, err
	}
	return formHandler.HandleUnsubmittedForm(ctx, r)
}

//HandleFormAction - handles post of personal data and updates cart
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

	//UpdatePurchaser (and add addittional data)
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

// MapPerson maps the checkout form data to the cart.Person domain struct
func (p DefaultPersonalDataForm) MapPerson() (*cart.Person) {
	person := cart.Person{
		PersonalDetails: cart.PersonalDetails{
			PassportNumber:  p.PassportNumber,
			PassportCountry: p.PassportCountry,
			DateOfBirth:     p.DateOfBirth,
		},
	}
	return &person
}


// MapAdditionalData - mapper
func (p DefaultPersonalDataForm) MapAdditionalData() (*cart.AdditionalData) {
	return nil
}
