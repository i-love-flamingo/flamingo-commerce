package forms

import (
	"context"
	"errors"
	"net/url"
	"strings"

	cartInterfaceForms "flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"

	"flamingo.me/form/domain"

	"flamingo.me/form/application"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	customerApplication "flamingo.me/flamingo-commerce/v3/customer/application"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//CheckoutFormComposite - a complete form (composite) for collecting all checkout data
	CheckoutFormComposite struct {
		//BillingAddressForm - the processed Form object for the BillingAddressForm
		// incoming form values are expected with namespace "billingAddress"
		BillingAddressForm *domain.Form
		//DeliveryForms - the processed Form object for the DeliveryForms
		// incoming form values are expected with namespace "deliveries.###DELIVERYCODE###"
		DeliveryForms map[string]*domain.Form
		//SimplePaymentForm - the processed Form object for the SimplePaymentForm
		// incoming form values are expected with namespace "payment"
		SimplePaymentForm *domain.Form
	}

	//checkoutFormBuilder - private builder for a form with CheckoutForm Data
	checkoutFormBuilder struct {
		checkoutForm *CheckoutFormComposite
	}

	// CheckoutFormController - the (mini) MVC for the complete checkout form
	CheckoutFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService
		userService                    *authApplication.UserService
		logger                         flamingo.Logger
		customerApplicationService     *customerApplication.Service
		formHandlerFactory             application.FormHandlerFactory
		billingAddressFormController   *cartInterfaceForms.BillingAddressFormController
		deliveryFormController         *cartInterfaceForms.DeliveryFormController
		simplePaymentFormController    *cartInterfaceForms.SimplePaymentFormController
	}
)

func (c *CheckoutFormController) Inject(responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	userService *authApplication.UserService,
	logger flamingo.Logger,
	customerApplicationService *customerApplication.Service,
	formHandlerFactory application.FormHandlerFactory,
	billingAddressFormController *cartInterfaceForms.BillingAddressFormController,
	deliveryFormController *cartInterfaceForms.DeliveryFormController,
	simplePaymentFormController *cartInterfaceForms.SimplePaymentFormController) {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService
	c.userService = userService
	c.customerApplicationService = customerApplicationService
	c.formHandlerFactory = formHandlerFactory
	c.logger = logger
	c.billingAddressFormController = billingAddressFormController
	c.deliveryFormController = deliveryFormController
	c.simplePaymentFormController = simplePaymentFormController
}

//GetUnsubmittedForm - Action that returns
func (c *CheckoutFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*CheckoutFormComposite, error) {
	checkoutFormBuilder := newCheckoutFormBuilder()

	//Add the billing form:
	billingForm, err := c.billingAddressFormController.GetUnsubmittedForm(ctx, r)
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}
	err = checkoutFormBuilder.addBillingForm(billingForm)
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}

	//Add a Delivery Form for every delivery:
	cart, err := c.applicationCartReceiverService.ViewCart(ctx, r.Session())
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}
	for _, delivery := range cart.Deliveries {
		deliveryForm, err := c.deliveryFormController.GetUnsubmittedForm(ctx, r)
		if err != nil {
			return checkoutFormBuilder.getForm(), err
		}
		err = checkoutFormBuilder.addDeliveryForm(delivery.DeliveryInfo.Code, deliveryForm)
		if err != nil {
			return checkoutFormBuilder.getForm(), err
		}
	}

	//3. Add the simplePaymentForm
	simplePaymentForm, err := c.simplePaymentFormController.GetUnsubmittedForm(ctx, r)
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}
	err = checkoutFormBuilder.addSimplePaymentForm(simplePaymentForm)
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}

	return checkoutFormBuilder.getForm(), nil

}

//newRequestWithResolvedNamespace - creates a new request with only the namespaced form values
func newRequestWithResolvedNamespace(namespace string, r *web.Request) *web.Request {
	newRequest := web.CreateRequest(r.Request(), r.Session())
	newUrlValues := make(url.Values)
	for k, values := range newRequest.Request().Form {
		if !strings.HasPrefix(k, namespace+".") {
			continue
		}
		strippedKey := strings.Replace(k, namespace+".", "", 1)
		newUrlValues[strippedKey] = values
	}
	newRequest.Request().Form = newUrlValues
	return newRequest
}

//HandleFormAction - Action that returns
func (c *CheckoutFormController) HandleFormAction(ctx context.Context, r *web.Request) (*CheckoutFormComposite, bool, error) {
	//session := web.SessionFromContext(ctx)
	checkoutFormBuilder := newCheckoutFormBuilder()
	overallSuccess := true
	err := r.Request().ParseForm()
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}

	//1, #### Process and add Billing Form Controller result
	billingForm, success, err := c.billingAddressFormController.HandleFormAction(ctx, newRequestWithResolvedNamespace("billingAddress", r))
	overallSuccess = overallSuccess && success
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}
	err = checkoutFormBuilder.addBillingForm(billingForm)
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}

	//2. #### Process ALL the delivery forms:
	//Add a Delivery Form for every delivery:
	cart, err := c.applicationCartReceiverService.ViewCart(ctx, r.Session())
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}
	for _, delivery := range cart.Deliveries {
		deliveryFormNamespace := "deliveries." + delivery.DeliveryInfo.Code
		//Add the billing form:
		deliverySubRequest := newRequestWithResolvedNamespace(deliveryFormNamespace, r)
		deliverySubRequest.Params["deliveryCode"] = delivery.DeliveryInfo.Code
		deliveryForm, success, err := c.deliveryFormController.HandleFormAction(ctx, deliverySubRequest)
		overallSuccess = overallSuccess && success
		if err != nil {
			return checkoutFormBuilder.getForm(), false, err
		}
		err = checkoutFormBuilder.addDeliveryForm(delivery.DeliveryInfo.Code, deliveryForm)
		if err != nil {
			return checkoutFormBuilder.getForm(), false, err
		}
	}

	//3. ### Add the simplePaymentForm
	simplePaymentForm, success, err := c.simplePaymentFormController.HandleFormAction(ctx, newRequestWithResolvedNamespace("payment", r))
	overallSuccess = overallSuccess && success
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}
	err = checkoutFormBuilder.addSimplePaymentForm(simplePaymentForm)
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}

	return checkoutFormBuilder.getForm(), overallSuccess, nil

}

func newCheckoutFormBuilder() *checkoutFormBuilder {
	b := &checkoutFormBuilder{
		checkoutForm: &CheckoutFormComposite{
			DeliveryForms: make(map[string]*domain.Form),
		},
	}
	return b
}

func (b *checkoutFormBuilder) getForm() *CheckoutFormComposite {
	return b.checkoutForm
}

func (b *checkoutFormBuilder) addDeliveryForm(deliveryCode string, deliveryForm *domain.Form) error {
	_, ok := deliveryForm.Data.(cartInterfaceForms.DeliveryForm)
	if !ok {
		return errors.New("no deliveryFormData?")
	}
	b.checkoutForm.DeliveryForms[deliveryCode] = deliveryForm
	return nil
}

func (b *checkoutFormBuilder) addBillingForm(billingForm *domain.Form) error {
	_, ok := billingForm.Data.(cartInterfaceForms.BillingAddressForm)
	if !ok {
		return errors.New("no billingFormData?")
	}
	b.checkoutForm.BillingAddressForm = billingForm
	return nil
}

func (b *checkoutFormBuilder) addSimplePaymentForm(simplePaymentForm *domain.Form) error {
	_, ok := simplePaymentForm.Data.(cartInterfaceForms.SimplePaymentForm)
	if !ok {
		return errors.New("no SimplePaymentForm?")
	}
	b.checkoutForm.SimplePaymentForm = simplePaymentForm
	return nil
}

//HasAnyGeneralErrors - true if any from the included forms has general errors
func (c *CheckoutFormComposite) HasAnyGeneralErrors() bool {
	errors := c.GetAllGeneralErrors()
	return len(errors) > 0
}

//GetAllGeneralErrors - gets all general errors from the included forms
func (c *CheckoutFormComposite) GetAllGeneralErrors() []domain.Error {
	var generalErrors []domain.Error
	generalErrors = append(generalErrors, c.BillingAddressForm.GetGeneralErrors()...)
	for _, form := range c.DeliveryForms {
		generalErrors = append(generalErrors, form.GetGeneralErrors()...)
	}
	return generalErrors
}
