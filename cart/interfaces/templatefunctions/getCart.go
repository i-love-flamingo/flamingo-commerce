package templatefunctions

import (

	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	// GetProduct is exported as a template function
	GetCart struct {
		ApplicationCartService *application.CartService `inject:""`
		Logger flamingo.Logger `inject:""`
	}
)

// Name alias for use in template
func (tf GetCart) Name() string {
	return "getCart"
}

func (tf GetCart) Func(ctx web.Context) interface{} {
	return func() cartDomain.Cart {
		cart, e := tf.ApplicationCartService.GetCart(ctx)
		if e != nil {
			tf.Logger.Error("Error: cart.interfaces.templatefunc %v", e)
		}
		return cart
	}
}
