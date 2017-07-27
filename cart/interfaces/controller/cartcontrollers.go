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
		Cart domain.Cart
		Test string
	}

	// CartController for carts
	CartViewController struct {
		responder.RenderAware  `inject:""`
		ApplicationCartService application.CartService `inject:""`
	}
)

// Get the Cart View ( / cart)
func (cc *CartViewController) Get(ctx web.Context) web.Response {

	cart, e := cc.ApplicationCartService.GetCart(ctx)
	if e != nil {
		fmt.Println(e)
		return cc.Render(ctx, "pages/checkout/carterror", nil)
	}
	return cc.Render(ctx, "pages/checkout/cart", CartViewData{
		Cart: cart,
		Test: "ddddd",
	})

}
