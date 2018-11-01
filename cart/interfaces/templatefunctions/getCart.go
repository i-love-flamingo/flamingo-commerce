package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/session"
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

func (tf *GetCart) Func(ctx context.Context) interface{} {
	return func() cartDomain.Cart {
		session, _ := session.FromContext(ctx)
		cart, e := tf.ApplicationCartReceiverService.ViewCart(ctx, session.G())
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

func (tf *GetDecoratedCart) Func(ctx context.Context) interface{} {
	return func() cartDomain.DecoratedCart {
		session, _ := session.FromContext(ctx)
		cart, e := tf.ApplicationCartReceiverService.ViewDecoratedCart(ctx, session.G())
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
