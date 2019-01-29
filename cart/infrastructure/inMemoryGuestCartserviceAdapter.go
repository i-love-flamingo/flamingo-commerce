package infrastructure

import (
	"context"
	"math/rand"
	"strconv"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
)

type (
	// InMemoryGuestCartService defines the in memory guest cart service
	InMemoryGuestCartService struct {
		inMemoryCartOrderBehaviour *InMemoryBehaviour
	}
)

var (
	_ cart.GuestCartService = (*InMemoryGuestCartService)(nil)
)

// Inject dependencies
func (gcs *InMemoryGuestCartService) Inject(
	InMemoryCartOrderBehaviour *InMemoryBehaviour,
) {
	gcs.inMemoryCartOrderBehaviour = InMemoryCartOrderBehaviour
}

// GetCart fetches a cart from the in memory guest cart service
func (gcs *InMemoryGuestCartService) GetCart(ctx context.Context, cartID string) (*cart.Cart, error) {
	cart, err := gcs.inMemoryCartOrderBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetNewCart gets a new cart from the in memory guest cart service
func (gcs *InMemoryGuestCartService) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	guestCart := &cart.Cart{
		ID: strconv.Itoa(rand.Int()),
	}

	error := gcs.inMemoryCartOrderBehaviour.StoreCart(guestCart)
	return guestCart, error
}

// GetBehaviour returns the cart order behaviour of the service
func (gcs *InMemoryGuestCartService) GetBehaviour(context.Context) (cart.Behaviour, error) {
	return gcs.inMemoryCartOrderBehaviour, nil
}
