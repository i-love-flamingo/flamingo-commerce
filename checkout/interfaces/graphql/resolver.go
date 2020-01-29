package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// CommerceCheckoutQueryResolver resolves graphql checkout queries
type CommerceCheckoutQueryResolver struct {
	placeOrderHandler *placeorder.Handler
}

// Inject dependencies
func (r *CommerceCheckoutQueryResolver) Inject(
	placeOrderHandler *placeorder.Handler,
	logger flamingo.Logger) {
	r.placeOrderHandler = placeOrderHandler
}

// CommerceCheckoutActivePlaceOrder checks if there is an order in unfinished state
func (r *CommerceCheckoutQueryResolver) CommerceCheckoutActivePlaceOrder(ctx context.Context) (bool, error) {
	return r.placeOrderHandler.HasUnfinishedProcess(ctx)
}
