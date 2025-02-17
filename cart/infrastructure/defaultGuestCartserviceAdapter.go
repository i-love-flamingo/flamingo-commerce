package infrastructure

import (
	"context"
	"math/rand"
	"strconv"

	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
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
	ctx, span := trace.StartSpan(ctx, "cart/DefaultGuestCartService/GetCart")
	defer span.End()

	cart, err := gcs.defaultBehaviour.GetCart(ctx, cartID)
	return cart, err
}

// GetNewCart gets a new cart from the in memory guest cart service
func (gcs *DefaultGuestCartService) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	ctx, span := trace.StartSpan(ctx, "cart/DefaultGuestCartService/GetNewCart")
	defer span.End()

	return gcs.defaultBehaviour.StoreNewCart(ctx, &cart.Cart{ID: strconv.Itoa(rand.Int())})
}

// GetModifyBehaviour returns the cart order behaviour of the service
func (gcs *DefaultGuestCartService) GetModifyBehaviour(ctx context.Context) (cart.ModifyBehaviour, error) {
	_, span := trace.StartSpan(ctx, "cart/DefaultGuestCartService/GetModifyBehaviour")
	defer span.End()

	return gcs.defaultBehaviour, nil
}

// RestoreCart restores a previously used guest cart
// Deprecated: (deprecated in the interface)
func (gcs *DefaultGuestCartService) RestoreCart(ctx context.Context, cart cart.Cart) (*cart.Cart, error) {
	_, span := trace.StartSpan(ctx, "cart/DefaultGuestCartService/RestoreCart")
	defer span.End()

	// RestoreCart restores a previously used cart
	gcs.logger.Warn("RestoreCart deprecated")
	return &cart, nil
}
