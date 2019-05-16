package forms

import (
	"context"

	"flamingo.me/form/application"
	"flamingo.me/form/domain"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//PersonalDataForm - the form for Person data
	PersonalDataForm struct {
		DateOfBirth     string      `form:"dateOfBirth"`
		PassportCountry string      `form:"passportCountry"`
		PassportNumber  string      `form:"passportNumber"`
		Address         AddressForm `form:"address" validate:"-"`
	}

	// PersonalDataFormService implements Form(Data)Provider interface of form package
	PersonalDataFormService struct {
	}

	// PersonalDataFormController - the (mini) MVC
	PersonalDataFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService

		logger flamingo.Logger

		formHandlerFactory application.FormHandlerFactory
	}
)

// GetFormData from data provider
func (p *PersonalDataFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
	return PersonalDataForm{}, nil
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

func (p *PersonalDataFormController) HandleFormAction(ctx context.Context, r *web.Request) (*domain.Form, bool, error) {
	return &domain.Form{}, true, nil // todo
}

func (p *PersonalDataFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*domain.Form, error) {
	return &domain.Form{}, nil // todo
}

// MapPerson maps the checkout form data to the cart.Person domain struct
func (p *PersonalDataForm) MapPerson() *cart.Person {
	person := cart.Person{
		PersonalDetails: cart.PersonalDetails{
			PassportNumber:  p.PassportNumber,
			PassportCountry: p.PassportCountry,
			DateOfBirth:     p.DateOfBirth,
		},
	}
	return &person
}
