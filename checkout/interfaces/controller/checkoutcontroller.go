package controller

import (
	"encoding/gob"

	application3 "go.aoe.com/flamingo/core/auth/application"
	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	application2 "go.aoe.com/flamingo/core/checkout/application"
	"go.aoe.com/flamingo/core/checkout/interfaces/controller/formDto"
	formApplicationService "go.aoe.com/flamingo/core/form/application"
	"go.aoe.com/flamingo/core/form/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// CheckoutViewData represents the checkout view data
	CheckoutViewData struct {
		DecoratedCart        cart.DecoratedCart
		Form                 domain.Form
		CartValidationResult cart.CartValidationResult
		ErrorMessage         string
		HasFormError         bool
	}

	// SuccessViewData represents the success view data
	SuccessViewData struct {
		OrderId string
		Email   string
	}

	// CheckoutController represents the checkout controller with its injectsions
	CheckoutController struct {
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		ApplicationCartService  application.CartService     `inject:""`
		PaymentService          application2.PaymentService `inject:""`
		UserService             application3.UserService    `inject:""`
		Router                  *router.Router              `inject:""`
		CheckoutFormService     domain.FormService          `inject:""`
		Logger                  flamingo.Logger             `inject:""`
		SourcingEngine          application2.SourcingEngine `inject:""`
		OrderService            application2.OrderService   `inject:""`
	}
)

func init() {
	gob.Register(SuccessViewData{})
}

// StartAction handles the checkout start action
func (cc *CheckoutController) StartAction(ctx web.Context) web.Response {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		cc.Logger.Errorf("cart.checkoutcontroller.viewaction: Error %v", e)
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
		HasFormError:  false,
	})
}

// SubmitUserCheckoutAction handles the user order submit
// TODO: implement this
func (cc *CheckoutController) SubmitUserCheckoutAction(ctx web.Context) web.Response {
	return cc.SubmitGuestCheckoutAction(ctx)
}

// SubmitGuestCheckoutAction handles the guest order submit
func (cc *CheckoutController) SubmitGuestCheckoutAction(ctx web.Context) web.Response {

	//Guard Clause if Cart cannout be fetched
	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		cc.Logger.Errorf("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	if cc.CheckoutFormService == nil {
		cc.Logger.Error("cart.checkoutcontroller.viewaction: Error CheckoutFormService not present!", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	form, e := formApplicationService.ProcessFormRequest(ctx, cc.CheckoutFormService)
	// return on error (template need to handle error display)
	if e != nil {
		return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
			DecoratedCart:        decoratedCart,
			CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
			Form:                 form,
		})
	}

	//Guard Clause if Cart is empty
	if decoratedCart.Cart.ItemCount() == 0 {
		return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
			DecoratedCart:        decoratedCart,
			CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
			Form:                 form,
		})
	}

	if form.IsValidAndSubmitted() {
		if checkoutFormData, ok := form.Data.(formDto.CheckoutFormData); ok {
			orderID, err := cc.placeOrder(ctx, checkoutFormData, decoratedCart)
			if err != nil {
				return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
					DecoratedCart:        decoratedCart,
					CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
					HasFormError:         true,
					Form:                 form,
					ErrorMessage:         err.Error(),
				})
			}
			shippingEmail := checkoutFormData.ShippingAddress.Email
			if shippingEmail == "" {
				shippingEmail = checkoutFormData.BillingAddress.Email
			}
			return cc.Redirect("checkout.success", nil).With("checkout.success.data", SuccessViewData{
				OrderId: orderID,
				Email:   shippingEmail,
			})
		}
	}

	return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
		DecoratedCart:        decoratedCart,
		CartValidationResult: cc.ApplicationCartService.ValidateCart(ctx, decoratedCart),
		Form:                 form,
	})
}

// SuccessAction handles the order success action
func (cc *CheckoutController) SuccessAction(ctx web.Context) web.Response {
	flashes := ctx.Session().Flashes("checkout.success.data")
	if len(flashes) > 0 {
		return cc.Render(ctx, "checkout/success", flashes[0].(SuccessViewData))
	}

	return cc.Render(ctx, "checkout/expired", nil)
}

func (cc *CheckoutController) placeOrder(ctx web.Context, checkoutFormData formDto.CheckoutFormData, decoratedCart cart.DecoratedCart) (string, error) {

	billingAddress, shippingAddress := formDto.MapAddresses(checkoutFormData)
	_ = shippingAddress
	return cc.OrderService.PlaceOrder(ctx, decoratedCart, "ispu", "ispu", billingAddress, billingAddress)
}
