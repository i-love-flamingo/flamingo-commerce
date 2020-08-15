package infrastructure

import (
	"context"
	"errors"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/oauth/domain"
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
func (cs *DefaultCustomerCartService) GetCart(ctx context.Context, auth domain.Auth, cartID string) (*cart.Cart, error) {
	customersCartID, err := cs.authToID(auth, cartID)
	if err == nil {
		return nil, err
	}
	foundCart, err := cs.defaultBehaviour.GetCart(ctx, customersCartID)
	if err == nil {
		return foundCart, err
	}
	if err == cart.ErrCartNotFound {
		cart := &cart.Cart{ID: customersCartID}
		cart.BelongsToAuthenticatedUser = true
		cart.AuthenticatedUserID = auth.IDToken.Subject
		return cs.defaultBehaviour.StoreNewCart(ctx, cart)
	}
	return nil, err
}

func (cs *DefaultCustomerCartService) authToID(auth domain.Auth, cartID string) (string, error) {
	if auth.IDToken == nil {
		return "", errors.New("Idtoken not given")
	}
	if auth.IDToken.Subject == "" {
		return "", errors.New("Idtoken empty")
	}
	return fmt.Sprintf("%v-%v", auth.IDToken.Subject, cartID), nil
}

// GetModifyBehaviour gets the cart order behaviour of the service
func (cs *DefaultCustomerCartService) GetModifyBehaviour(context.Context, domain.Auth) (cart.ModifyBehaviour, error) {
	return cs.defaultBehaviour, nil
}

// RestoreCart restores a previously used cart
// Deprecated: (deprecated in the interface)
func (cs *DefaultCustomerCartService) RestoreCart(ctx context.Context, auth domain.Auth, cart cart.Cart) (*cart.Cart, error) {
	cs.logger.Warn("RestoreCart depricated")
	return &cart, nil
}
