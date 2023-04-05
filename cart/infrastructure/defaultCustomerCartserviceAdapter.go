package infrastructure

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	//DefaultCustomerCartService defines the in memory customer cart service
	DefaultCustomerCartService struct {
		defaultBehaviour *DefaultCartBehaviour
		logger           flamingo.Logger
	}
)

var (
	_ cart.CustomerCartService = (*DefaultCustomerCartService)(nil)
)

// Inject dependencies
func (cs *DefaultCustomerCartService) Inject(
	behaviour *DefaultCartBehaviour,
	logger flamingo.Logger,
) {
	cs.defaultBehaviour = behaviour
	cs.logger = logger
}

// GetCart gets a customer cart from the in memory customer cart service
func (cs *DefaultCustomerCartService) GetCart(ctx context.Context, identity auth.Identity, _ string) (*cart.Cart, error) {
	id := identity.Subject()

	foundCart, err := cs.defaultBehaviour.GetCart(ctx, id)
	if err == nil {
		return foundCart, nil
	}

	if errors.Is(err, cart.ErrCartNotFound) {
		cart := &cart.Cart{ID: id}
		cart.BelongsToAuthenticatedUser = true
		cart.AuthenticatedUserID = id

		return cs.defaultBehaviour.StoreNewCart(ctx, cart)
	}

	return nil, err
}

// GetModifyBehaviour gets the cart order behaviour of the service
func (cs *DefaultCustomerCartService) GetModifyBehaviour(context.Context, auth.Identity) (cart.ModifyBehaviour, error) {
	return cs.defaultBehaviour, nil
}

// RestoreCart restores a previously used cart
// Deprecated: (deprecated in the interface)
func (cs *DefaultCustomerCartService) RestoreCart(ctx context.Context, auth auth.Identity, cart cart.Cart) (*cart.Cart, error) {
	cs.logger.Warn("RestoreCart depricated")
	return &cart, nil
}
