package forms

import (
	"context"
	"errors"

	"flamingo.me/form/domain"

	"flamingo.me/form/application"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	customerApplication "flamingo.me/flamingo-commerce/v3/customer/application"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//BillingAddressForm - the form for billing address
	BillingAddressForm AddressForm

	// BillingAddressFormService implements Form(Data)Provider interface of form package
	BillingAddressFormService struct {
		customerApplicationService *customerApplication.Service
		userService                *authApplication.UserService
	}

	// BillingAddressFormController - the (mini) MVC
	BillingAddressFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService
		userService                    *authApplication.UserService
		logger                         flamingo.Logger
		customerApplicationService     *customerApplication.Service
		formHandlerFactory             application.FormHandlerFactory
		billingAddressFormProvider     *BillingAddressFormService
	}
)

// Inject - dependencies
func (p *BillingAddressFormService) Inject(customerApplicationService *customerApplication.Service, userService *authApplication.UserService) {
	p.customerApplicationService = customerApplicationService
	p.userService = userService
}

// GetFormData from data provider
func (p *BillingAddressFormService) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
	session := web.SessionFromContext(ctx)
	billingAddressForm := AddressForm{}
	if p.userService.IsLoggedIn(ctx, session) {
		customer, err := p.customerApplicationService.GetForAuthenticatedUser(ctx, session)
		if err != nil {
			billingAddress := customer.GetDefaultBillingAddress()
			if billingAddress != nil {
				billingAddressForm.LoadFromCustomerAddress(*billingAddress)
			}
		}
	}
	return BillingAddressForm(billingAddressForm), nil
}

func (c *BillingAddressFormController) Inject(responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	userService *authApplication.UserService,
	logger flamingo.Logger,
	customerApplicationService *customerApplication.Service,
	formHandlerFactory application.FormHandlerFactory,
	billingAddressFormProvider *BillingAddressFormService) {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService
	c.userService = userService
	c.customerApplicationService = customerApplicationService
	c.formHandlerFactory = formHandlerFactory
	c.logger = logger
	c.billingAddressFormProvider = billingAddressFormProvider
}

func (c *BillingAddressFormController) getFormHandler() (domain.FormHandler, error) {
	builder := c.formHandlerFactory.GetFormHandlerBuilder()
	err := builder.SetNamedFormService("commerce.cart.billingFormService")
	if err != nil {
		return nil, err
	}
	return builder.Build(), nil
}

//GetUnsubmittedForm - returns unsubmitted
func (c *BillingAddressFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*domain.Form, error) {
	formHandler, err := c.getFormHandler()
	if err != nil {
		return nil, err
	}
	return formHandler.HandleUnsubmittedForm(ctx, r)
}

//HandleFormAction - return the form or error. If the form was submitted the action is performed
func (c *BillingAddressFormController) HandleFormAction(ctx context.Context, r *web.Request) (form *domain.Form, actionSuccessFull bool, err error) {
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
	billingAddressForm, ok := form.Data.(BillingAddressForm)
	if !ok {
		return form, false, errors.New("cannot convert to AddressForm ")
	}
	if !form.IsValidAndSubmitted() {
		return form, false, nil
	}
	addressForm := AddressForm(billingAddressForm)
	billingAddress := addressForm.MapToDomainAddress()

	//update Billing
	err = c.applicationCartService.UpdateBillingAddress(ctx, session, &billingAddress)
	if err != nil {
		c.logger.Error("BillingAddressFormController UpdateBillingAddress Error %v", err)
		return form, false, err
	}
	return form, true, nil
}
