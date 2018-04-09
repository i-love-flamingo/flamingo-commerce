package infrastructure

import (
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"context"
	"strconv"
	"math/rand"
)

type (
	InMemoryGuestCartService struct {
		InMemoryCartOrderBehaviour *InMemoryCartOrderBehaviour `inject:""`
	}
)

var (
	_ cart.GuestCartService = (*InMemoryGuestCartService)(nil)
)

func (gcs *InMemoryGuestCartService) GetCart(ctx context.Context, cartId string) (*cart.Cart, error) {
	cart, err := gcs.InMemoryCartOrderBehaviour.GetCart(ctx, cartId)
	return cart, err
}

func (gcs *InMemoryGuestCartService) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	guestCart := cart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}

	error := gcs.InMemoryCartOrderBehaviour.CartStorage.StoreCart(guestCart)
	return &guestCart, error
}

func (gcs *InMemoryGuestCartService) GetCartOrderBehaviour(context.Context) (cart.CartBehaviour, error) {
	return gcs.InMemoryCartOrderBehaviour, nil
}
