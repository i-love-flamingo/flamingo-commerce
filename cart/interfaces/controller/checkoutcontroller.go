package controller

import (
	"log"

	"go.aoe.com/flamingo/core/breadcrumbs"
	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"

	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// ViewData is used for cart views/templates
	CheckoutViewData struct {
		DecoratedCart cart.DecoratedCart
	}

	// CheckoutController for carts
	CheckoutController struct {
		responder.RenderAware  `inject:""`
		ApplicationCartService application.CartService `inject:""`
		Router                 *router.Router          `inject:""`
	}
)

// TODO - all this should trigger stateflow
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
		Url:   cc.Router.URL("cart.view", nil).String(),
	})

	return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
		DecoratedCart: decoratedCart,
	})
}
