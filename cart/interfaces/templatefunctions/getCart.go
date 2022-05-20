package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
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
	tf.logger = logger.WithField(flamingo.LogKeyModule, "cart").WithField(flamingo.LogKeyCategory, "getCart")
}

// Func defines the GetCart template function
func (tf *GetCart) Func(ctx context.Context) interface{} {
	return func() cartDomain.Cart {
		session := web.SessionFromContext(ctx)
		cart, e := tf.cartReceiverService.ViewCart(ctx, session)
		if e != nil {
			tf.logger.WithContext(ctx).Error("Error: cart.interfaces.templatefunc %v", e)
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

// Func defines the GetDecoratedCart template function
func (tf *GetDecoratedCart) Func(ctx context.Context) interface{} {
	return func() decorator.DecoratedCart {
		session := web.SessionFromContext(ctx)
		cart, e := tf.cartReceiverService.ViewDecoratedCart(ctx, session)
		if e != nil {
			if e == application.ErrTemporaryCartService {
				tf.logger.WithContext(ctx).Warn("Error: cart.interfaces.templatefunc %v", e)
			} else {
				tf.logger.WithContext(ctx).Error("Error: cart.interfaces.templatefunc %v", e)
			}
		}
		if cart == nil {
			cart = &decorator.DecoratedCart{}
		}
		//dereference (should lower risk of undesired sideeffects if missused in template)
		return *cart
	}
}
