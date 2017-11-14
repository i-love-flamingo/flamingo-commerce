package controller

import (
	"log"

	"go.aoe.com/flamingo/core/breadcrumbs"
	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"

	checkoutApplication "go.aoe.com/flamingo/core/checkout/application"
	"go.aoe.com/flamingo/core/magento/infrastructure/cartservice"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// CheckoutViewData is used for cart views/templates
	CheckoutViewData struct {
		DecoratedCart cart.DecoratedCart
		ErrorMessage  string
		HasError      bool
	}

	// SuccessViewData is used for cart views/templates
	SuccessViewData struct {
		OrderId      string
		ErrorMessage string
		HasError     bool
	}

	// CheckoutController for carts
	CheckoutController struct {
		responder.RenderAware          `inject:""`
		ApplicationCartService         application.CartService                    `inject:""`
		Router                         *router.Router                             `inject:""`
		ContextService                 checkoutApplication.ContextService         `inject:""`
		MagentoGuestCartServiceAdapter cartservice.MagentoGuestCartServiceAdapter `inject:""`
	}
)

func (cc *CheckoutController) StartAction(ctx web.Context) web.Response {

	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	breadcrumbs.Add(ctx, breadcrumbs.Crumb{
		Title: "Shopping Bag",
		Url:   cc.Router.URL("cart.view", nil).String(),
	})
	breadcrumbs.Add(ctx, breadcrumbs.Crumb{
		Title: "Reserve & Collect",
		Url:   cc.Router.URL("checkout.start", nil).String(),
	})

	return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
		DecoratedCart: decoratedCart,
		HasError:      false,
	})
}

func (cc *CheckoutController) SubmitAction(ctx web.Context) web.Response {

	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.checkoutcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	breadcrumbs.Add(ctx, breadcrumbs.Crumb{
		Title: "Shopping Bag",
		Url:   cc.Router.URL("cart.view", nil).String(),
	})
	breadcrumbs.Add(ctx, breadcrumbs.Crumb{
		Title: "Reserve & Collect",
		Url:   cc.Router.URL("checkout.start", nil).String(),
	})

	billingAddress := cart.Address{
		Lastname:    "Pötzinger",
		Firstname:   "Daniel",
		Email:       "poetzinger@aoe.com",
		City:        "Wiesbaden",
		PostCode:    "65183",
		CountryCode: "DE",
		Street:      "Luisenstraße",
		StreetNr:    "6",
	}
	shippingAddress := billingAddress
	err := cc.MagentoGuestCartServiceAdapter.SetShippingInformation(ctx, decoratedCart.Cart.ID, &shippingAddress, &billingAddress, "flatrate", "flatrate")
	if err != nil {
		cc.Render(ctx, "checkout/checkout", CheckoutViewData{
			DecoratedCart: decoratedCart,
			ErrorMessage:  err.Error(),
			HasError:      true,
		})
	}
	payment := cart.Payment{
		Method: "checkmo",
	}

	err, orderid := cc.MagentoGuestCartServiceAdapter.PlaceOrder(ctx, decoratedCart.Cart.ID, &payment)
	if err != nil {
		cc.Render(ctx, "checkout/checkout", CheckoutViewData{
			DecoratedCart: decoratedCart,
			ErrorMessage:  err.Error(),
			HasError:      true,
		})
	}
	return cc.Render(ctx, "checkout/success", SuccessViewData{
		OrderId: orderid,
	})
}
