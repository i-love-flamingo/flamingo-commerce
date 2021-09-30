package states

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// PrepareCart state
	PrepareCart struct {
		cartService *application.CartService
	}
)

var _ process.State = PrepareCart{}

// Inject dependencies
func (v *PrepareCart) Inject(
	cartService *application.CartService,
) *PrepareCart {
	v.cartService = cartService

	return v
}

// Name get state name
func (PrepareCart) Name() string {
	return "PrepareCart"
}

// Run the state operations
func (v PrepareCart) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/PrepareCart/Run")
	defer span.End()

	c, err := v.cartService.ForceReserveOrderIDAndSave(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	if c.GrandTotal.IsZero() {
		p.UpdateState(ValidateCart{}.Name(), nil)
		p.UpdateCart(*c)
		return process.RunResult{}
	}

	if c.PaymentSelection == nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: cart.ErrPaymentSelectionNotSet.Error()},
		}
	}

	paymentSelection, err := c.PaymentSelection.GenerateNewIdempotencyKey()
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	err = v.cartService.UpdatePaymentSelection(ctx, web.SessionFromContext(ctx), paymentSelection)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateState(ValidateCart{}.Name(), nil)
	p.UpdateCart(*c)
	return process.RunResult{}
}

// Rollback the state operations
func (v PrepareCart) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (v PrepareCart) IsFinal() bool {
	return false
}
