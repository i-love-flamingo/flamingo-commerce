package graphql

import (
	"context"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// CommerceCheckoutQueryResolver resolves graphql checkout mutations
type CommerceCheckoutQueryResolver struct {
	placeorderHandler *placeorder.Handler
	orderService      *application.OrderService
	cartService       *cartApplication.CartService
	logger            flamingo.Logger
}

// Inject dependencies
func (r *CommerceCheckoutQueryResolver) Inject(
	placeorderHandler *placeorder.Handler,
	orderService *application.OrderService,
	cartService *cartApplication.CartService,
	logger flamingo.Logger) {
	r.placeorderHandler = placeorderHandler
	r.orderService = orderService
	r.cartService = cartService
	r.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "graphql")

}

//CommerceCheckoutActivePlaceOrder checks if there is an order in unfinished state
func (r *CommerceCheckoutQueryResolver) CommerceCheckoutActivePlaceOrder(ctx context.Context) (bool, error) {
	return r.placeorderHandler.HasUnfinishedProcess(ctx)
}
