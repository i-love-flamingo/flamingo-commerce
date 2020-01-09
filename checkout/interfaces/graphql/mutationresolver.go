package graphql

import (
	"context"
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutMutationResolver struct {
	placeorderHandler *placeorder.Handler
	cartService       *cartApplication.CartService
	logger            flamingo.Logger
}

// Inject dependencies
func (r *CommerceCheckoutMutationResolver) Inject(
	placeorderHandler *placeorder.Handler,
	cartService *cartApplication.CartService,
	logger flamingo.Logger) {
	r.placeorderHandler = placeorderHandler
	r.cartService = cartService
	r.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "graphql")

}

//CommerceCheckoutStartPlaceOrder starts a new process (if not running)
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutStartPlaceOrder(ctx context.Context, returnURLRaw string) (bool, error) {
	session := web.SessionFromContext(ctx)
	cart, err := r.cartService.GetCartReceiverService().ViewCart(ctx, session)
	if err != nil {
		return false, err
	}
	startPlaceOrderCommand := placeorder.StartPlaceOrderCommand{Cart: *cart}
	_, err = r.placeorderHandler.StartPlaceOrder(ctx, startPlaceOrderCommand)
	if err == placeorder.ErrAnotherPlaceOrderProcessRunning {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

//CommerceCheckoutCancelPlaceOrder starts a new process (if not running)
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutCancelPlaceOrder(ctx context.Context) (bool, error) {
	return true, nil
}
