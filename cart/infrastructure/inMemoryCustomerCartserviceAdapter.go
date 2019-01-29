package infrastructure

import (
	"context"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/core/auth/domain"
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

// GetBehaviour gets the cart order behaviour of the service
func (gcs *InMemoryCustomerCartService) GetBehaviour(context.Context, domain.Auth) (cart.Behaviour, error) {
	return gcs.inMemoryBehaviour, nil
}
