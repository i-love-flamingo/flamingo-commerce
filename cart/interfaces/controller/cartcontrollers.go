package controller

import (
	"flamingo/core/cart/application"
	"flamingo/core/cart/domain"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"fmt"
)

type (
	// ViewData is used for cart views/templates
	CartViewData struct {
		DecoratedCart domain.DecoratedCart
	}

	// CartController for carts
	CartViewController struct {
		responder.RenderAware  `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}
)

// Get the DecoratedCart View ( / cart)
func (cc *CartViewController) Get(ctx web.Context) web.Response {

	cart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		fmt.Println(e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}
	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart: *cart,
	})

}
