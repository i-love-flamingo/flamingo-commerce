package infrastructure

import (
	"context"
	"math/rand"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// InMemoryGuestCartService defines the in memory guest cart service
	InMemoryGuestCartService struct {
		inMemoryBehaviour *InMemoryBehaviour
	}
)

var (
	_ cart.GuestCartService = (*InMemoryGuestCartService)(nil)
)

// Inject dependencies
func (gcs *InMemoryGuestCartService) Inject(
	InMemoryCartOrderBehaviour *InMemoryBehaviour,
) {
	gcs.inMemoryBehaviour = InMemoryCartOrderBehaviour
}

// GetCart fetches a cart from the in memory guest cart service
func (gcs *InMemoryGuestCartService) GetCart(ctx context.Context, cartID string) (*cart.Cart, error) {
	cart, err := gcs.inMemoryBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetNewCart gets a new cart from the in memory guest cart service
func (gcs *InMemoryGuestCartService) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	guestCart := &cart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}

	err := gcs.inMemoryBehaviour.StoreCart(guestCart)
	return guestCart, err
}

// GetModifyBehaviour returns the cart order behaviour of the service
func (gcs *InMemoryGuestCartService) GetModifyBehaviour(context.Context) (cart.ModifyBehaviour, error) {
	return gcs.inMemoryBehaviour, nil
}

// RestoreCart restores a previously used guest cart
func (gcs *InMemoryGuestCartService) RestoreCart(ctx context.Context, cart cart.Cart) (*cart.Cart, error) {
	guestCart := cart
	guestCart.ID = strconv.Itoa(rand.Int())

	err := gcs.inMemoryBehaviour.StoreCart(&guestCart)
	return &guestCart, err
}
