package states

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// CompleteCart state
	CompleteCart struct {
		cartService         *application.CartService
		cartReceiverService *application.CartReceiverService
	}
	// CompleteCartRollbackData for later rollbacks
	CompleteCartRollbackData struct {
		completedCart *cart.Cart
	}
)

var _ process.State = CompleteCart{}

// Inject dependencies
func (c *CompleteCart) Inject(
	cartService *application.CartService,
	cartReceiverService *application.CartReceiverService,
) *CompleteCart {
	c.cartService = cartService
	c.cartReceiverService = cartReceiverService

	return c
}

// Name get state name
func (CompleteCart) Name() string {
	return "CompleteCart"
}

// Run the state operations
func (c CompleteCart) Run(ctx context.Context, p *process.Process) process.RunResult {
	behaviour, err := c.cartReceiverService.ModifyBehaviour(ctx)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	_, ok := behaviour.(cart.CompleteBehaviour)
	if !ok {
		// cart does not support completing, proceed to place order
		p.UpdateState(PlaceOrder{}.Name())
		return process.RunResult{}
	}

	completedCart, err := c.cartService.CompleteCurrentCart(ctx)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateState(PlaceOrder{}.Name())
	return process.RunResult{
		RollbackData: &CompleteCartRollbackData{
			completedCart: completedCart,
		},
	}
}

// Rollback the state operations
func (c CompleteCart) Rollback(data process.RollbackData) error {
	rollbackData, ok := data.(CompleteCartRollbackData)
	if !ok {
		return fmt.Errorf("rollback data not of expected type 'CompleteCartRollbackData', but %T", rollbackData)
	}

	// todo: context is missing but needed here to restore cart cache..
	_, err := c.cartService.RestoreCart(context.Background(), rollbackData.completedCart)
	if err != nil {
		return err
	}

	return nil
}

// IsFinal if state is a final state
func (c CompleteCart) IsFinal() bool {
	return false
}
