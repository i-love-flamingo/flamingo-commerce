package graphql

import (
	"context"
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutMutationResolver struct {
	orderService         *application.OrderService
	decoratedCartFactory *decorator.DecoratedCartFactory
	cartService          *cartApplication.CartService
	logger               flamingo.Logger
}

// Inject dependencies
func (r *CommerceCheckoutMutationResolver) Inject(
	orderService *application.OrderService,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	cartService *cartApplication.CartService,
	logger flamingo.Logger) {
	r.orderService = orderService
	r.decoratedCartFactory = decoratedCartFactory
	r.cartService = cartService
	r.logger = logger.WithField(flamingo.LogKeyModule, "om3oms").WithField(flamingo.LogKeyCategory, "graphql")

}

//CommerceCheckoutStartPlaceOrder starts a new process (if not running)
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutStartPlaceOrder(ctx context.Context, returnURLRaw string) (bool, error) {
	return true, nil
}

//CommerceCheckoutStartPlaceOrder starts a new process (if not running)
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutCancelPlaceOrder(ctx context.Context) (bool, error) {
	return true, nil
}
