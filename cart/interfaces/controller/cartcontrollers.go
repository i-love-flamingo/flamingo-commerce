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
		Items         []domain.DecoratedCartItem
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
	fmt.Printf("%+v", decoratedCart.Cartitems)
	return cc.Render(ctx, "checkout/cart", CartViewData{
		DecoratedCart: decoratedCart,
		Items:         decoratedCart.Cartitems,
	})

}
