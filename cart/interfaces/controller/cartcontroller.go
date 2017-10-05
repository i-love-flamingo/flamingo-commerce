package controller

import (
	"log"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// ViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart cart.DecoratedCart
		Items         []cart.DecoratedCartItem
	}

	// CartController for carts
	CartViewController struct {
		responder.RenderAware  `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}
)

// ViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) ViewAction(ctx web.Context) web.Response {

	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		log.Printf("cart.cartcontroller.viewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart: decoratedCart,
		Items:         decoratedCart.Cartitems,
	})

}

// AddAndViewAction the DecoratedCart View ( / cart)
func (cc *CartViewController) AddAndViewAction(ctx web.Context) web.Response {
	addRequest := AddRequestFromRequestContext(ctx)
	e := cc.ApplicationCartService.AddProduct(ctx, addRequest)
	if e != nil {
		log.Printf("cart.cartcontroller.addandviewaction: Error %v", e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}
	return cc.ViewAction(ctx)

}
