package domain

import (
	"context"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/core/auth/domain"
)

type (
	// GuestOrderService interface
	GuestOrderService interface {
		GetBehaviour(context.Context) (Behaviour, error)
	}

	// CustomerOrderService interface
	CustomerOrderService interface {
		Get(context.Context, domain.Auth) ([]*Order, error)
		GetByID(context.Context, domain.Auth, string) (*Order, error)

		GetBehaviour(context.Context, domain.Auth) (Behaviour, error)
	}

	// Behaviour is a Port that can be implemented by other packages to provide cart order actions
	Behaviour interface {
		PlaceOrder(ctx context.Context, cart *cart.Cart, payment *cart.CartPayment) (PlacedOrderInfos, error)
	}
)
