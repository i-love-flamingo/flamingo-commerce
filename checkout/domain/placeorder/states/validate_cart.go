package states

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// ValidateCart state
	ValidateCart struct {
		cartService *application.CartService
	}
)

var _ process.State = ValidateCart{}

// Inject dependencies
func (v *ValidateCart) Inject(
	cartService *application.CartService,
) *ValidateCart {
	v.cartService = cartService

	return v
}

// Name get state name
func (ValidateCart) Name() string {
	return "ValidateCart"
}

// Run the state operations
func (v ValidateCart) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/ValidateCart/Run")
	defer span.End()

	result, err := v.cartService.ValidateCurrentCart(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	if !result.IsValid() {
		return process.RunResult{
			Failed: process.CartValidationErrorReason{
				ValidationResult: result,
			},
		}
	}

	if p.Context().Cart.GrandTotal.IsZero() {
		p.UpdateState(CompleteCart{}.Name(), nil)
		return process.RunResult{}
	}

	p.UpdateState(ValidatePaymentSelection{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (v ValidateCart) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (v ValidateCart) IsFinal() bool {
	return false
}
