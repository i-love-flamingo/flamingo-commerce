package controller

import (
	"encoding/gob"

	authApplication "go.aoe.com/flamingo/core/auth/application"
	cartApplication "go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/checkout/application"
	paymentDomain "go.aoe.com/flamingo/core/checkout/domain/payment"
	"go.aoe.com/flamingo/core/checkout/interfaces/controller/formDto"
	customerApplication "go.aoe.com/flamingo/core/customer/application"
	formApplicationService "go.aoe.com/flamingo/core/form/application"
	"go.aoe.com/flamingo/core/form/domain"
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
		Form                 domain.Form
		CartValidationResult cart.CartValidationResult
		ErrorMessage         string
		HasSubmitError       bool
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
		Router *router.Router   `inject:""`

		CheckoutFormService  *formDto.CheckoutFormService `inject:""`
		OrderService         *application.OrderService    `inject:""`
		PaymentService       *application.PaymentService  `inject:""`
		DecoratedCartFactory *cart.DecoratedCartFactory   `inject:""`

		ApplicationCartService *cartApplication.CartService `inject:""`

		UserService *authApplication.UserService `inject:""`

		Logger flamingo.Logger `inject:""`

		CustomerApplicationService *customerApplication.Service `inject:""`
		PaymentProvider            PaymentProviderProvider      `inject:""`
	}
)

func init() {
	gob.Register(PlaceOrderFlashData{})
}

// StartAction handles the checkout start action
func (cc *CheckoutController) StartAction(ctx web.Context) web.Response {
	methods := cc.getPaymentMethods()
	paymark := methods["paymark"]


	paymentMethod := paymark.GetPaymentMethods()[0]
	test, error :=paymark.RedirectExternalPayment(ctx, paymentMethod)

	cc.Logger.Info("hallo: %v error: %v", test, error)

	return test


	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		cc.Logger.WithField("category", "checkout").Errorf("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if cc.UserService.IsLoggedIn(ctx) {
		return cc.Redirect("checkout.user", nil)
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/startcheckout", CheckoutViewData{
			DecoratedCart: decoratedCart,
		})
	}

	return cc.Render(ctx, "checkout/startcheckout", CheckoutViewData{
		DecoratedCart: decoratedCart,
	})
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

func (cc *CheckoutController) submitOrderForm(ctx web.Context, formservice *formDto.CheckoutFormService, template string) web.Response {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		cc.Logger.WithField("category", "checkout").Errorf("cart.checkoutcontroller.submitaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if formservice == nil {
		cc.Logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error CheckoutFormService not present!")
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	form, e := formApplicationService.ProcessFormRequest(ctx, formservice)
	// return on error (template need to handle error display)
	if e != nil {
		return cc.Render(ctx, template, CheckoutViewData{
			DecoratedCart:        decoratedCart,
			CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
			Form:                 form,
		})
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, template, CheckoutViewData{
			DecoratedCart:        decoratedCart,
			CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
			Form:                 form,
		})
	}

	if form.IsValidAndSubmitted() {
		if checkoutFormData, ok := form.Data.(formDto.CheckoutFormData); ok {
			orderID, err := cc.placeOrder(ctx, checkoutFormData, decoratedCart)
			if err != nil {
				//Place Order Error
				return cc.Render(ctx, template, CheckoutViewData{
					DecoratedCart:        decoratedCart,
					CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
					HasSubmitError:       true,
					Form:                 form,
					ErrorMessage:         err.Error(),
				})
			}
			shippingEmail := checkoutFormData.ShippingAddress.Email
			if shippingEmail == "" {
				shippingEmail = checkoutFormData.BillingAddress.Email
			}
			return cc.Redirect("checkout.success", nil).With("checkout.success.data", PlaceOrderFlashData{
				OrderId:     orderID,
				Email:       shippingEmail,
				PlacedItems: decoratedCart.Cart.Cartitems,
				CartTotals:  decoratedCart.Cart.GetCartTotals(),
			})
		} else {
			cc.Logger.WithField("category", "checkout").Error("cart.checkoutcontroller.submitaction: Error cannot type convert to CheckoutFormData!")
			return cc.Render(ctx, "checkout/carterror", nil)
		}
	} else {
		if form.IsSubmitted && form.HasGeneralErrors() {
			cc.Logger.WithField("category", "checkout").Warnf("CheckoutForm has general error: %#v", form.ValidationInfo.GeneralErrors)
		}
	}
	//Default: Form not submitted yet or submitted with validation errors:
	return cc.Render(ctx, template, CheckoutViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
		Form:                 form,
	})
}

func (cc *CheckoutController) placeOrder(ctx web.Context, checkoutFormData formDto.CheckoutFormData, decoratedCart cart.DecoratedCart) (string, error) {
	billingAddress, shippingAddress := formDto.MapAddresses(checkoutFormData)
	return cc.OrderService.PlaceOrder(ctx, decoratedCart, "ispu", "ispu", billingAddress, shippingAddress)
}

func (cc *CheckoutController) getPaymentMethods() map[string]paymentDomain.PaymentProvider {
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
