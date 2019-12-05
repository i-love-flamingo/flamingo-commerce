package forms

import (
	"context"
	"errors"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"

	"flamingo.me/form/domain"

	"flamingo.me/form/application"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//SimplePaymentForm - the form for simple select of payment gateway
	SimplePaymentForm struct {
		Gateway string `form:"gateway"  validate:"required"`
		Method  string `form:"method"  validate:"required"`
	}

	// SimplePaymentFormService implements Form(Data)Provider interface of form package
	SimplePaymentFormService struct {
		userService                    *authApplication.UserService
		applicationCartReceiverService *cartApplication.CartReceiverService
	}

	// SimplePaymentFormController - the (mini) MVC
	SimplePaymentFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService

		logger flamingo.Logger

		formHandlerFactory       application.FormHandlerFactory
		simplePaymentFormService *SimplePaymentFormService
	}
)

// Inject - dependencies
func (p *SimplePaymentFormService) Inject(
	userService *authApplication.UserService,
	applicationCartReceiverService *cartApplication.CartReceiverService) {
	p.userService = userService
	p.applicationCartReceiverService = applicationCartReceiverService
}

// GetFormData from data provider
func (p *SimplePaymentFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
	session := web.SessionFromContext(ctx)

	cart, err := p.applicationCartReceiverService.ViewCart(ctx, session)

	if err != nil {
		return SimplePaymentForm{}, nil
	}

	if cart.PaymentSelection != nil {
		return SimplePaymentForm{
			Gateway: cart.PaymentSelection.Gateway(),
			Method:  cart.PaymentSelection.MethodByType(priceDomain.ChargeTypeMain),
		}, nil
	}

	return SimplePaymentForm{}, nil
}

//Inject - Inject
func (c *SimplePaymentFormController) Inject(responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	logger flamingo.Logger,
	formHandlerFactory application.FormHandlerFactory,
	simplePaymentFormService *SimplePaymentFormService) {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService

	c.formHandlerFactory = formHandlerFactory
	c.logger = logger.WithField(flamingo.LogKeyModule, "cart").WithField(flamingo.LogKeyCategory, "simplepaymentform")
	c.simplePaymentFormService = simplePaymentFormService
}

func (c *SimplePaymentFormController) getFormHandler() (domain.FormHandler, error) {
	builder := c.formHandlerFactory.GetFormHandlerBuilder()
	err := builder.SetFormService(c.simplePaymentFormService)
	if err != nil {
		return nil, err
	}
	return builder.Build(), nil
}

//GetUnsubmittedForm - returns unsubmitted
func (c *SimplePaymentFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*domain.Form, error) {
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, err
	}
	return formHandler.HandleUnsubmittedForm(ctx, r)
}

//HandleFormAction - return the form or error. If the form was submitted the action is performed
func (c *SimplePaymentFormController) HandleFormAction(ctx context.Context, r *web.Request) (form *domain.Form, actionSuccessFull bool, err error) {
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
	simplePaymentForm, ok := form.Data.(SimplePaymentForm)
	if !ok {
		return form, false, errors.New("cannot convert to SimplePaymentForm ")
	}
	if !form.IsValidAndSubmitted() {
		return form, false, nil
	}
	currentCart, err := c.applicationCartReceiverService.ViewCart(ctx, r.Session())
	if err != nil {
		return nil, false, err
	}
	paymentSelection := simplePaymentForm.MapToPaymentSelection(currentCart)

	//update cart
	err = c.applicationCartService.UpdatePaymentSelection(ctx, session, paymentSelection)
	if err != nil {
		c.logger.WithContext(ctx).Error("SimplePaymentFormController UpdatePaymentSelection Error %v", err)
		return form, false, err
	}
	return form, true, nil
}

//MapToPaymentSelection - mapper from form values to domain
func (f *SimplePaymentForm) MapToPaymentSelection(currentCart *cartDomain.Cart) cartDomain.PaymentSelection {
	chargeTypeToPaymentMethod := map[string]string{
		priceDomain.ChargeTypeMain: f.Method,
	}
	selection, _ := cartDomain.NewDefaultPaymentSelection(f.Gateway, chargeTypeToPaymentMethod, *currentCart)
	return selection
}
