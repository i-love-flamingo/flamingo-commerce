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
		cartReceiverService *application.CartReceiverService
		logger              flamingo.Logger
	}
	// GetDecoratedCart is exported as a template function
	GetDecoratedCart struct {
		cartReceiverService *application.CartReceiverService
		logger              flamingo.Logger
	}
)

// Inject dependencies
func (tf *GetCart) Inject(
	applicationCartReceiverService *application.CartReceiverService,
	logger flamingo.Logger,

) {
	tf.cartReceiverService = applicationCartReceiverService
	tf.logger = logger
}

// Func defines the GetCart template function
func (tf *GetCart) Func(ctx context.Context) interface{} {
	return func() cartDomain.Cart {
		session, _ := session.FromContext(ctx)
		cart, e := tf.cartReceiverService.ViewCart(ctx, session.G())
		if e != nil {
			tf.logger.Error("Error: cart.interfaces.templatefunc %v", e)
		}
		if cart == nil {
			cart = &cartDomain.Cart{}
		}
		//dereference (should lower risk of undesired sideeffects if missused in template)
		return *cart
	}
}

// Inject dependencies
func (tf *GetDecoratedCart) Inject(
	cartReceiverService *application.CartReceiverService,
	logger flamingo.Logger,
) {
	tf.cartReceiverService = cartReceiverService
	tf.logger = logger
}

// Func defines the GetDecoratedCrt deplate function
func (tf *GetDecoratedCart) Func(ctx context.Context) interface{} {
	return func() cartDomain.DecoratedCart {
		session, _ := session.FromContext(ctx)
		cart, e := tf.cartReceiverService.ViewDecoratedCart(ctx, session.G())
		if e != nil {
			tf.logger.Error("Error: cart.interfaces.templatefunc %v", e)
		}
		if cart == nil {
			cart = &cartDomain.DecoratedCart{}
		}
		//dereference (should lower risk of undesired sideeffects if missused in template)
		return *cart
	}
}
