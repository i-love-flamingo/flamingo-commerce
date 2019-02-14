package inmemory

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/order/domain"
)

type (
	// GuestOrderService defines the in memory guest order service
	GuestOrderService struct {
		behaviour *Behaviour
	}
)

var (
	_ domain.GuestOrderService = (*GuestOrderService)(nil)
)

// Inject dependencies
func (os *GuestOrderService) Inject(
	Behaviour *Behaviour,
) {
	os.behaviour = Behaviour
}

// GetBehaviour gets the in memory guest order service behaviour
func (os *GuestOrderService) GetBehaviour(ctx context.Context) (domain.Behaviour, error) {
	return os.behaviour, nil
}
