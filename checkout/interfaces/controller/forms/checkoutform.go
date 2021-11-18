package forms

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form/application"
	"flamingo.me/form/domain"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartInterfaceForms "flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
)

type (
	// CheckoutFormComposite - a complete form (composite) for collecting all checkout data
	CheckoutFormComposite struct {
		// BillingAddressForm - the processed Form object for the BillingAddressForm
		// incoming form values are expected with namespace "billingAddress"
		BillingAddressForm *domain.Form
		// DeliveryForms - the processed Form object for the DeliveryForms
		// incoming form values are expected with namespace "deliveries.###DELIVERYCODE###"
		DeliveryForms map[string]*domain.Form
		// SimplePaymentForm - the processed Form object for the SimplePaymentForm
		// incoming form values are expected with namespace "payment"
		SimplePaymentForm *domain.Form
		// PersonalDataForm - the processed Form object for personal data
		PersonalDataForm *domain.Form
	}

	// checkoutFormBuilder - private builder for a form with CheckoutForm Data
	checkoutFormBuilder struct {
		checkoutForm *CheckoutFormComposite
	}

	// CheckoutFormController - the (mini) MVC for the complete checkout form
	CheckoutFormController struct {
		responder                      *web.Responder
		applicationCartService         *cartApplication.CartService
		applicationCartReceiverService *cartApplication.CartReceiverService
		logger                         flamingo.Logger
		formHandlerFactory             application.FormHandlerFactory
		billingAddressFormController   *cartInterfaceForms.BillingAddressFormController
		deliveryFormController         *cartInterfaceForms.DeliveryFormController
		simplePaymentFormController    *cartInterfaceForms.SimplePaymentFormController
		personalDataFormController     *cartInterfaceForms.PersonalDataFormController
		useDeliveryForms               bool
		usePersonalDataForm            bool
	}
)

// Inject dependencies
func (c *CheckoutFormController) Inject(responder *web.Responder,
	applicationCartService *cartApplication.CartService,
	applicationCartReceiverService *cartApplication.CartReceiverService,
	logger flamingo.Logger,
	formHandlerFactory application.FormHandlerFactory,
	billingAddressFormController *cartInterfaceForms.BillingAddressFormController,
	deliveryFormController *cartInterfaceForms.DeliveryFormController,
	simplePaymentFormController *cartInterfaceForms.SimplePaymentFormController,
	personalDataFormController *cartInterfaceForms.PersonalDataFormController,
	config *struct {
		UseDeliveryForms    bool `inject:"config:commerce.checkout.useDeliveryForms"`
		UsePersonalDataForm bool `inject:"config:commerce.checkout.usePersonalDataForm"`
	},
) {
	c.responder = responder
	c.applicationCartReceiverService = applicationCartReceiverService
	c.applicationCartService = applicationCartService
	c.formHandlerFactory = formHandlerFactory
	c.logger = logger
	c.billingAddressFormController = billingAddressFormController
	c.deliveryFormController = deliveryFormController
	c.simplePaymentFormController = simplePaymentFormController
	c.personalDataFormController = personalDataFormController
	if config != nil {
		c.useDeliveryForms = config.UseDeliveryForms
		c.usePersonalDataForm = config.UsePersonalDataForm
	}
}

// GetUnsubmittedForm returns the form
func (c *CheckoutFormController) GetUnsubmittedForm(ctx context.Context, r *web.Request) (*CheckoutFormComposite, error) {
	checkoutFormBuilder := newCheckoutFormBuilder()

	// Add the billing form:
	billingForm, err := c.billingAddressFormController.GetUnsubmittedForm(ctx, r)
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}
	err = checkoutFormBuilder.addBillingForm(billingForm)
	if err != nil {
		return checkoutFormBuilder.getForm(), err
	}

	if c.useDeliveryForms {
		// Add a Delivery Form for every delivery:
		cart, err := c.applicationCartReceiverService.ViewCart(ctx, r.Session())
		if err != nil {
			return checkoutFormBuilder.getForm(), err
		}
		for _, delivery := range cart.Deliveries {
			if !delivery.HasItems() {
				continue
			}
			r.Params["deliveryCode"] = delivery.DeliveryInfo.Code
			deliveryForm, err := c.deliveryFormController.GetUnsubmittedForm(ctx, r)
			if err != nil {
				return checkoutFormBuilder.getForm(), err
			}
			err = checkoutFormBuilder.addDeliveryForm(delivery.DeliveryInfo.Code, deliveryForm)
			if err != nil {
				return checkoutFormBuilder.getForm(), err
			}
		}
	}

	if c.usePersonalDataForm {
		// 3. Personal Data
		personalDataForm, err := c.personalDataFormController.GetUnsubmittedForm(ctx, r)
		if err != nil {
			return checkoutFormBuilder.getForm(), err
		}
		err = checkoutFormBuilder.addPersonalDataForm(personalDataForm)
		if err != nil {
			return checkoutFormBuilder.getForm(), err
		}
	}

	// 4. Add the simplePaymentForm
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

// newRequestWithResolvedNamespace creates a new request with only the namespaced form values
func newRequestWithResolvedNamespace(namespace string, r *web.Request) *web.Request {
	newRequest := web.CreateRequest(r.Request(), r.Session())
	newURLValues := make(url.Values)
	for k, values := range newRequest.Request().Form {
		if !strings.HasPrefix(k, namespace+".") {
			continue
		}
		strippedKey := strings.Replace(k, namespace+".", "", 1)
		newURLValues[strippedKey] = values
	}
	newRequest.Request().Form = newURLValues
	return newRequest
}

// HandleFormAction handles the submitted form request
func (c *CheckoutFormController) HandleFormAction(ctx context.Context, r *web.Request) (*CheckoutFormComposite, bool, error) {
	checkoutFormBuilder := newCheckoutFormBuilder()
	overallSuccess := true
	err := r.Request().ParseForm()
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}

	cart, err := c.applicationCartReceiverService.ViewCart(ctx, r.Session())
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}

	// 1, #### Process and add Billing Form Controller result
	billingForm, success, err := c.billingAddressFormController.HandleFormAction(ctx, newRequestWithResolvedNamespace("billingAddress", r))
	overallSuccess = overallSuccess && success
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}
	err = checkoutFormBuilder.addBillingForm(billingForm)
	if err != nil {
		return checkoutFormBuilder.getForm(), false, err
	}

	if c.useDeliveryForms {
		// 2. #### Process ALL the delivery forms:
		// Add a Delivery Form for every delivery:
		for _, delivery := range cart.Deliveries {
			if !delivery.HasItems() {
				continue
			}
			deliveryFormNamespace := "deliveries." + delivery.DeliveryInfo.Code
			// Add the billing form:
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
	}

	if c.usePersonalDataForm {
		// 3. ### Add the personalDataForm
		personalDataForm, success, err := c.personalDataFormController.HandleFormAction(ctx, newRequestWithResolvedNamespace("personalData", r))
		overallSuccess = overallSuccess && success
		if err != nil {
			return checkoutFormBuilder.getForm(), false, err
		}
		err = checkoutFormBuilder.addPersonalDataForm(personalDataForm)
		if err != nil {
			return checkoutFormBuilder.getForm(), false, err
		}
	}

	if !cart.GrandTotal.IsZero() {
		// 4. ### Add the simplePaymentForm if payment is required.
		simplePaymentForm, success, err := c.simplePaymentFormController.HandleFormAction(ctx, newRequestWithResolvedNamespace("payment", r))
		overallSuccess = overallSuccess && success
		if err != nil {
			return checkoutFormBuilder.getForm(), false, err
		}
		err = checkoutFormBuilder.addSimplePaymentForm(simplePaymentForm)
		if err != nil {
			return checkoutFormBuilder.getForm(), false, err
		}
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

func (b *checkoutFormBuilder) addPersonalDataForm(personalDataForm *domain.Form) error {
	b.checkoutForm.PersonalDataForm = personalDataForm
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

// HasAnyGeneralErrors checks if any from the included forms has general errors
func (c *CheckoutFormComposite) HasAnyGeneralErrors() bool {
	errs := c.GetAllGeneralErrors()
	return len(errs) > 0
}

// GetAllGeneralErrors from the included forms
func (c *CheckoutFormComposite) GetAllGeneralErrors() []domain.Error {
	var generalErrors []domain.Error
	generalErrors = append(generalErrors, c.BillingAddressForm.GetGeneralErrors()...)
	for _, form := range c.DeliveryForms {
		generalErrors = append(generalErrors, form.GetGeneralErrors()...)
	}
	return generalErrors
}
