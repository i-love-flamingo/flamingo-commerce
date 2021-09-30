package infrastructure

import (
	"context"
	"math/rand"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// DefaultGuestCartService defines the in memory guest cart service
	DefaultGuestCartService struct {
		defaultBehaviour *DefaultCartBehaviour
		logger           flamingo.Logger
	}
)

var (
	_ cart.GuestCartService = (*DefaultGuestCartService)(nil)
)

// Inject dependencies
func (gcs *DefaultGuestCartService) Inject(
	InMemoryCartOrderBehaviour *DefaultCartBehaviour,
	logger flamingo.Logger,
) {
	gcs.defaultBehaviour = InMemoryCartOrderBehaviour
	gcs.logger = logger
}

// GetCart fetches a cart from the in memory guest cart service
func (gcs *DefaultGuestCartService) GetCart(ctx context.Context, cartID string) (*cart.Cart, error) {
	cart, err := gcs.defaultBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetNewCart gets a new cart from the in memory guest cart service
func (gcs *DefaultGuestCartService) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	return gcs.defaultBehaviour.StoreNewCart(ctx, &cart.Cart{ID: strconv.Itoa(rand.Int())})
}

// GetModifyBehaviour returns the cart order behaviour of the service
func (gcs *DefaultGuestCartService) GetModifyBehaviour(context.Context) (cart.ModifyBehaviour, error) {
	return gcs.defaultBehaviour, nil
}

// RestoreCart restores a previously used guest cart
// Deprecated: (deprecated in the interface)
func (gcs *DefaultGuestCartService) RestoreCart(ctx context.Context, cart cart.Cart) (*cart.Cart, error) {
	// RestoreCart restores a previously used cart
	gcs.logger.Warn("RestoreCart deprecated")
	return &cart, nil
}
