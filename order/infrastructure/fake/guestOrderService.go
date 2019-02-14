package fake

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/order/domain"
	"flamingo.me/flamingo-commerce/v3/order/infrastructure/inmemory"
)

type (
	// GuestOrders is the fake guest order api service
	GuestOrders struct{}
)

var (
	_ domain.GuestOrderService = (*GuestOrders)(nil)
)

// GetBehaviour returns the behaviour of the fake servide
func (g *GuestOrders) GetBehaviour(context.Context) (domain.Behaviour, error) {
	return new(inmemory.Behaviour), nil
}

