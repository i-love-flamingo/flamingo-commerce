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
func (gcs *DefaultCustomerCartService) Inject(
	behaviour *DefaultCartBehaviour,
	logger flamingo.Logger,
) {
	gcs.defaultBehaviour = behaviour
	gcs.logger = logger
}

// GetCart gets a customer cart from the in memory customer cart service
func (gcs *DefaultCustomerCartService) GetCart(ctx context.Context, auth domain.Auth, cartID string) (*cart.Cart, error) {
	customersCartId, err := gcs.authToId(auth, cartID)
	foundCart, err := gcs.defaultBehaviour.GetCart(ctx, customersCartId)
	if err == nil {
		return foundCart, err
	}
	if err == cart.ErrCartNotFound {
		return gcs.defaultBehaviour.NewCart(ctx, customersCartId)
	}
	return nil, err
}

func (cob *DefaultCustomerCartService) authToId(auth domain.Auth, cartID string) (string, error) {
	if auth.IDToken == nil {
		return "", errors.New("Idtoken not given")
	}
	if auth.IDToken.Subject == "" {
		return "", errors.New("Idtoken empty")
	}
	return fmt.Sprintf("%v-%v", auth.IDToken.Subject, cartID), nil
}

// GetModifyBehaviour gets the cart order behaviour of the service
func (gcs *DefaultCustomerCartService) GetModifyBehaviour(context.Context, domain.Auth) (cart.ModifyBehaviour, error) {
	return gcs.defaultBehaviour, nil
}

// RestoreCart restores a previously used cart
func (gcs *DefaultCustomerCartService) RestoreCart(ctx context.Context, auth domain.Auth, cart cart.Cart) (*cart.Cart, error) {
	gcs.logger.Warn("RestoreCart depricated")
	return &cart, nil
}
