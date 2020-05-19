package infrastructure

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	//DefaultCustomerCartService defines the in memory customer cart service
	DefaultCustomerCartService struct {
		defaultBehaviour *DefaultCartBehaviour
	}
)

var (
	_ cart.CustomerCartService = (*DefaultCustomerCartService)(nil)
)

// Inject dependencies
func (gcs *DefaultCustomerCartService) Inject(
	InMemoryBehaviour *DefaultCartBehaviour,
) {
	gcs.defaultBehaviour = InMemoryBehaviour
}

// GetCart gets a customer cart from the in memory customer cart service
func (gcs *DefaultCustomerCartService) GetCart(ctx context.Context, auth domain.Auth, cartID string) (*cart.Cart, error) {
	cart, err := gcs.defaultBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetModifyBehaviour gets the cart order behaviour of the service
func (gcs *DefaultCustomerCartService) GetModifyBehaviour(context.Context, domain.Auth) (cart.ModifyBehaviour, error) {
	return gcs.defaultBehaviour, nil
}

// RestoreCart restores a previously used cart
func (gcs *DefaultCustomerCartService) RestoreCart(ctx context.Context, auth domain.Auth, cart cart.Cart) (*cart.Cart, error) {
	customerCart := cart

	err := gcs.defaultBehaviour.storeCart(&customerCart)
	return &customerCart, err
}
