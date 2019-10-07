package graphql

import (
	"context"
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartModel "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutMutationResolver struct {
	cartService    *cartApplication.CartService
	graphQLService *Service
}

// Inject dependencies
func (r *CommerceCheckoutMutationResolver) Inject(
	cartService *cartApplication.CartService,
	graphQLService *Service) {
	r.cartService = cartService
	r.graphQLService = graphQLService
}

// CommerceUpdateBillingAddress
func (r *CommerceCheckoutMutationResolver) CommerceUpdateBillingAddress(ctx context.Context, address *cartModel.Address) (*cartModel.Address, error) {
	session := web.SessionFromContext(ctx)
	err := r.cartService.UpdateBillingAddress(ctx, session, address)
	if err != nil {
		return nil, err
	}

	cart, err := r.cartService.GetCartReceiverService().ViewDecoratedCart(ctx, session)
	if err != nil {
		return nil, err
	}
	return cart.Cart.BillingAddress, nil
}
