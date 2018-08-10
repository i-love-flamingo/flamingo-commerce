package templatefunctions

import (
	"flamingo.me/flamingo-commerce/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (
	// GetCart is exported as a template function
	GetCart struct {
		ApplicationCartReceiverService *application.CartReceiverService `inject:""`
		Logger                         flamingo.Logger                  `inject:""`
	}
	// GetDecoratedCart is exported as a template function
	GetDecoratedCart struct {
		ApplicationCartReceiverService *application.CartReceiverService `inject:""`
		Logger                         flamingo.Logger                  `inject:""`
	}
)

// Name alias for use in template
func (tf GetCart) Name() string {
	return "getCart"
}

func (tf GetCart) Func(ctx web.Context) interface{} {
	return func() cartDomain.Cart {
		cart, e := tf.ApplicationCartReceiverService.ViewCart(ctx, ctx.Session())
		if e != nil {
			tf.Logger.Error("Error: cart.interfaces.templatefunc %v", e)
		}
		if cart == nil {
			cart = &cartDomain.Cart{}
		}
		//dereference (should lower risk of undesired sideeffects if missused in template)
		return *cart
	}
}

// Name alias for use in template
func (tf GetDecoratedCart) Name() string {
	return "getDecoratedCart"
}

func (tf GetDecoratedCart) Func(ctx web.Context) interface{} {
	return func() cartDomain.DecoratedCart {
		cart, e := tf.ApplicationCartReceiverService.ViewDecoratedCart(ctx, ctx.Session())
		if e != nil {
			tf.Logger.Error("Error: cart.interfaces.templatefunc %v", e)
		}
		if cart == nil {
			cart = &cartDomain.DecoratedCart{}
		}
		//dereference (should lower risk of undesired sideeffects if missused in template)
		return *cart
	}
}
