package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// CloseCart state
	CloseCart struct {
		cartService application.CartService
	}
)

var _ process.State = CloseCart{}

// Inject dependencies
func (c *CloseCart) Inject(
	service application.CartService,
) *CloseCart {
	c.cartService = service

	return c
}

// Name get state name
func (CloseCart) Name() string {
	return "CloseCart"
}

// Run the state operations
func (c CloseCart) Run(_ context.Context, p *process.Process) process.RunResult {
	p.UpdateState(CreatePayment{}.Name())

	return process.RunResult{}
}

// Rollback the state operations
func (c CloseCart) Rollback(process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (c CloseCart) IsFinal() bool {
	return false
}
