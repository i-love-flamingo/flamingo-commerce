package infrastructure

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	//InMemoryCustomerCartService defines the in memory customer cart service
	InMemoryCustomerCartService struct {
		inMemoryBehaviour *InMemoryBehaviour
	}
)

var (
	_ cart.CustomerCartService = (*InMemoryCustomerCartService)(nil)
)

// Inject dependencies
func (gcs *InMemoryCustomerCartService) Inject(
	InMemoryBehaviour *InMemoryBehaviour,
) {
	gcs.inMemoryBehaviour = InMemoryBehaviour
}

// GetCart gets a customer cart from the in memory customer cart service
func (gcs *InMemoryCustomerCartService) GetCart(ctx context.Context, auth domain.Auth, cartID string) (*cart.Cart, error) {
	cart, err := gcs.inMemoryBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetModifyBehaviour gets the cart order behaviour of the service
func (gcs *InMemoryCustomerCartService) GetModifyBehaviour(context.Context, domain.Auth) (cart.ModifyBehaviour, error) {
	return gcs.inMemoryBehaviour, nil
}

// RestoreCart restores a previously used cart
func (gcs *InMemoryCustomerCartService) RestoreCart(ctx context.Context, auth domain.Auth, cart cart.Cart) (*cart.Cart, error) {
	customerCart := cart

	err := gcs.inMemoryBehaviour.StoreCart(&customerCart)
	return &customerCart, err
}
