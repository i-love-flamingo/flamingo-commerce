package controller

import (
	"fmt"

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

// Get the DecoratedCart View ( / cart)
func (cc *CartViewController) Get(ctx web.Context) web.Response {

	decoratedCart, e := cc.ApplicationCartService.GetDecoratedCart(ctx)
	if e != nil {
		fmt.Println(e)
		return cc.Render(ctx, "checkout/carterror", nil)
	}

	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart: decoratedCart,
		Items:         decoratedCart.Cartitems,
	})

}
