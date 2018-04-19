package controller

import (
	"encoding/gob"
	"errors"
	"net/url"
	"strings"

	authApplication "go.aoe.com/flamingo/core/auth/application"
	canonicalApp "go.aoe.com/flamingo/core/canonicalUrl/application"
	cartApplication "go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/checkout/application"
	paymentDomain "go.aoe.com/flamingo/core/checkout/domain/payment"
	"go.aoe.com/flamingo/core/checkout/interfaces/controller/formDto"
	customerApplication "go.aoe.com/flamingo/core/customer/application"
	formApplicationService "go.aoe.com/flamingo/core/form/application"
	formDomain "go.aoe.com/flamingo/core/form/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	PaymentProviderProvider func() map[string]paymentDomain.PaymentProvider

	// CheckoutViewData represents the checkout view data
	CheckoutViewData struct {
		DecoratedCart        cart.DecoratedCart
		Form                 formDomain.Form
		CartValidationResult cart.CartValidationResult
		ErrorMessage         string
		HasSubmitError       bool
		PaymentProviders     map[string]paymentDomain.PaymentProvider
	}

	// SuccessViewData represents the success view data
	SuccessViewData struct {
		OrderId              string
		Email                string
		PlacedDecoratedItems []cart.DecoratedCartItem
		CartTotals           cart.CartTotals
	}

	// PlaceOrderFlashData represents the data passed to the success page - they need to be "glob"able
	PlaceOrderFlashData struct {
		OrderId string
		Email   string
		//Encodeable cart data to pass
		PlacedItems []cart.Item
		CartTotals  cart.CartTotals
	}

	// CheckoutController represents the checkout controller with its injectsions
	CheckoutController struct {
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		Router                  *router.Router `inject:""`

		CheckoutFormService  *formDto.CheckoutFormService `inject:""`
		OrderService         *application.OrderService    `inject:""`
		PaymentService       *application.PaymentService  `inject:""`
		DecoratedCartFactory *cart.DecoratedCartFactory   `inject:""`

		SkipStartAction bool `inject:"config:checkout.skipStartAction,optional"`

		ApplicationCartService         *cartApplication.CartService         `inject:""`
		ApplicationCartReceiverService *cartApplication.CartReceiverService `inject:""`

		UserService *authApplication.UserService `inject:""`

		Logger flamingo.Logger `inject:""`

		CustomerApplicationService *customerApplication.Service `inject:""`
		PaymentProvider            PaymentProviderProvider      `inject:""`

		CanonicalService *canonicalApp.Service `inject:""`
	}
)

func init() {
	gob.Register(PlaceOrderFlashData{})
}

// StartAction handles the checkout start action
func (cc *CheckoutController) StartAction(ctx web.Context) web.Response {
	//Guard Clause if Cart cannout be fetched

	decoratedCart, e := cc.ApplicationCartReceiverService.ViewDecoratedCart(ctx)
	if e != nil {
		cc.Logger.WithField("category", "checkout").Errorf("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/startcheckout", CheckoutViewData{
			DecoratedCart: *decoratedCart,
		})
	}

	if cc.UserService.IsLoggedIn(ctx) {
		return cc.Redirect("checkout.user", nil)
	}

	if cc.SkipStartAction {
		return cc.Redirect("checkout.guest", nil)
	}

	return cc.Render(ctx, "checkout/startcheckout", CheckoutViewData{
		DecoratedCart: *decoratedCart,
	})
}

func (cc *CheckoutController) hasAvailablePaymentProvider() bool {
	return len(cc.getPaymentProviders()) > 0
}

func (cc *CheckoutController) getPayment(ctx web.Context, paymentProviderCode string, paymentMethodCode string) (paymentDomain.PaymentProvider, *paymentDomain.PaymentMethod, error) {
	providers := cc.getPaymentProviders()

	provider := providers[paymentProviderCode]

	if provider == nil {
		return nil, nil, errors.New("Payment provider " + paymentProviderCode + " not found")
	}

	paymentMethods := provider.GetPaymentMethods()

	var paymentMethod *paymentDomain.PaymentMethod
	for _, method := range paymentMethods {
		if method.Code == paymentMethodCode {
			paymentMethod = &method
			break
		}
	}

	if paymentMethod == nil {
		return nil, nil, errors.New("payment method not found")
	}
	return provider, paymentMethod, nil
}

// SubmitUserCheckoutAction handles the user order submit
func (cc *CheckoutController) SubmitUserCheckoutAction(ctx web.Context) web.Response {
	//Guard
	if !cc.UserService.IsLoggedIn(ctx) {
		return cc.Redirect("checkout.start", nil)
	}
	customer, err := cc.CustomerApplicationService.GetForAuthenticatedUser(ctx)
	if err == nil {
		//give the customer to the form service - so that it can prepopulate default values
		cc.CheckoutFormService.Customer = customer
	}

	return cc.submitOrderForm(ctx, cc.CheckoutFormService, "checkout/usercheckout")
}

// SubmitGuestCheckoutAction handles the guest order submit
func (cc *CheckoutController) SubmitGuestCheckoutAction(ctx web.Context) web.Response {
	cc.CheckoutFormService.Customer = nil
	if cc.UserService.IsLoggedIn(ctx) {
		return cc.Redirect("checkout.user", nil)
	}
	return cc.submitOrderForm(ctx, cc.CheckoutFormService, "checkout/guestcheckout")
}

// ProcessPaymentAction functions as a return/notification URL for Payment Providers
func (cc *CheckoutController) ProcessPaymentAction(ctx web.Context) web.Response {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.ApplicationCartReceiverService.ViewDecoratedCart(ctx)
	if e != nil {
		cc.Logger.WithField("category", "checkout").Errorf("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	providercode := ctx.MustParam1("providercode")
	methodcode := ctx.MustParam1("methodcode")

	//In case we need to show an error we need to know the template and the formdata:
	template := "checkout/guestcheckout"
	if cc.UserService.IsLoggedIn(ctx) {
		template = "checkout/usercheckout"
	}
	checkoutForm := formDomain.Form{} //TODO - get from session or from cart...
	email := "todo"

	provider, paymentMethod, err := cc.getPayment(ctx, providercode, methodcode)

	cartPayment, err := provider.ProcessPayment(ctx, &decoratedCart.Cart, paymentMethod, nil)
	if err != nil {
		return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, checkoutForm, err)
	}

	response, err := cc.placeOrder(ctx, *cartPayment, email, *decoratedCart)
	if err != nil {
		return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, checkoutForm, err)
	}
	return response
}

// SuccessAction handles the order success action
func (cc *CheckoutController) SuccessAction(ctx web.Context) web.Response {
	flashes := ctx.Session().Flashes("checkout.success.data")
	if len(flashes) > 0 {
		if placeOrderFlashData, ok := flashes[0].(PlaceOrderFlashData); ok {
			viewData := SuccessViewData{
				CartTotals:           placeOrderFlashData.CartTotals,
				Email:                placeOrderFlashData.Email,
				OrderId:              placeOrderFlashData.OrderId,
				PlacedDecoratedItems: cc.DecoratedCartFactory.CreateDecorateCartItems(ctx, placeOrderFlashData.PlacedItems),
			}
			return cc.Render(ctx, "checkout/success", viewData)
		}
	}

	return cc.Render(ctx, "checkout/expired", nil)
}

func (cc *CheckoutController) getPaymentReturnUrl(PaymentProvider string, PaymentMethod string) *url.URL {
	baseUrl := cc.CanonicalService.BaseUrl
	paymentUrl := cc.Router.URL("checkout.processpayment", router.P{"providercode": PaymentProvider, "methodcode": PaymentMethod})

	rawUrl := strings.TrimRight(baseUrl, "/") + paymentUrl.String()

	urlResult, _ := url.Parse(rawUrl)

	return urlResult
}

func (cc *CheckoutController) submitOrderForm(ctx web.Context, formservice *formDto.CheckoutFormService, template string) web.Response {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.ApplicationCartReceiverService.ViewDecoratedCart(ctx)
	if e != nil {
		cc.Logger.WithField("category", "checkout").Errorf("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if formservice == nil {
		cc.Logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error CheckoutFormService not present!")
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if !cc.hasAvailablePaymentProvider() {
		cc.Logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error No Payment set")
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	form, e := formApplicationService.ProcessFormRequest(ctx, formservice)
	// return on error (template need to handle error display)
	if e != nil {
		return cc.Render(ctx, template, CheckoutViewData{
			DecoratedCart:        *decoratedCart,
			CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
			Form:                 form,
			PaymentProviders:     cc.getPaymentProviders(),
		})
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, template, CheckoutViewData{
			DecoratedCart:        *decoratedCart,
			CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
			Form:                 form,
			PaymentProviders:     cc.getPaymentProviders(),
		})
	}

	if form.IsValidAndSubmitted() {

		if checkoutFormData, ok := form.Data.(formDto.CheckoutFormData); ok {

			billingAddress, shippingAddress := formDto.MapAddresses(checkoutFormData)
			person := formDto.MapPerson(checkoutFormData)

			err := cc.OrderService.CurrentCartSaveInfos(ctx, billingAddress, shippingAddress, person)
			if err != nil {
				return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, form, err)
			}
			//procces Payment:
			paymentProvider, paymentMethod, err := cc.getPayment(ctx, checkoutFormData.SelectedPaymentProvider, checkoutFormData.SelectedPaymentProviderMethod)
			if err != nil {
				return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, form, err)
			}
			//Payment Method requests an redirect - execute it
			if paymentMethod.IsExternalPayment {
				returnUrl := cc.getPaymentReturnUrl(paymentProvider.GetCode(), paymentMethod.Code)
				hostedPaymentPageResponse, err := paymentProvider.RedirectExternalPayment(ctx, &decoratedCart.Cart, paymentMethod, returnUrl)
				if err != nil {
					return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, form, err)
				}
				return hostedPaymentPageResponse
			}

			//Paymentmethod that need no external Redirect - can be processed directly
			cartPayment, err := paymentProvider.ProcessPayment(ctx, &decoratedCart.Cart, paymentMethod, nil)
			if err != nil {
				return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, form, err)
			}
			shippingEmail := checkoutFormData.ShippingAddress.Email
			if shippingEmail == "" {
				shippingEmail = checkoutFormData.BillingAddress.Email
			}
			response, err := cc.placeOrder(ctx, *cartPayment, shippingEmail, *decoratedCart)
			if err != nil {
				return cc.placeOrderErrorResponse(ctx, template, *decoratedCart, form, err)
			}
			return response

		} else {
			cc.Logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error cannot type convert to CheckoutFormData!")
			return cc.Render(ctx, "checkout/carterror", nil)
		}
	} else {
		if form.IsSubmitted && form.HasGeneralErrors() {
			cc.Logger.WithField("category", "checkout").Warnf("CheckoutForm has general error: %#v", form.ValidationInfo.GeneralErrors)
		}
	}

	cc.Logger.Debugf("paymentProviders %#v", cc.getPaymentProviders())
	//Default: Form not submitted yet or submitted with validation errors:
	return cc.Render(ctx, template, CheckoutViewData{
		DecoratedCart:        *decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
		Form:                 form,
		PaymentProviders:     cc.getPaymentProviders(),
	})
}

func (cc *CheckoutController) placeOrderErrorResponse(ctx web.Context, template string, decoratedCart cart.DecoratedCart, form formDomain.Form, err error) web.Response {
	cc.Logger.Warnf("Place Order Error: %s", err.Error())
	return cc.Render(ctx, template, CheckoutViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, &decoratedCart),
		HasSubmitError:       true,
		Form:                 form,
		ErrorMessage:         err.Error(),
		PaymentProviders:     cc.getPaymentProviders(),
	})

}

func (cc *CheckoutController) placeOrder(ctx web.Context, cartPayment cart.CartPayment, email string, decoratedCart cart.DecoratedCart) (web.Response, error) {
	orderID, err := cc.OrderService.CurrentCartPlaceOrder(ctx, cartPayment)
	if err != nil {
		return nil, err
	}

	return cc.Redirect("checkout.success", nil).With("checkout.success.data", PlaceOrderFlashData{
		OrderId:     orderID,
		Email:       email,
		PlacedItems: decoratedCart.Cart.Cartitems,
		CartTotals:  decoratedCart.Cart.CartTotals,
	}), nil

}
func (cc *CheckoutController) getPaymentProviders() map[string]paymentDomain.PaymentProvider {
	result := make(map[string]paymentDomain.PaymentProvider)

	paymentProviders := cc.PaymentProvider()

	if paymentProviders != nil {
		for name, paymentProvider := range cc.PaymentProvider() {
			if paymentProvider.IsActive() {
				result[name] = paymentProvider
			}
		}
	}

	return result
}
