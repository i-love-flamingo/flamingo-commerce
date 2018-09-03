package infrastructure

import (
	"context"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
)

type (
	InMemoryCustomerCartService struct {
		InMemoryCartOrderBehaviour *InMemoryCartOrderBehaviour `inject:""`
	}
)

var (
	_ cart.CustomerCartService = (*InMemoryCustomerCartService)(nil)
)

func (gcs *InMemoryCustomerCartService) GetCart(ctx context.Context, auth cart.Auth, cartId string) (*cart.Cart, error) {
	cart, err := gcs.InMemoryCartOrderBehaviour.GetCart(ctx, cartId)
	return cart, err
}

func (gcs *InMemoryCustomerCartService) GetCartOrderBehaviour(context.Context, cart.Auth) (cart.CartBehaviour, error) {
	return gcs.InMemoryCartOrderBehaviour, nil
}
