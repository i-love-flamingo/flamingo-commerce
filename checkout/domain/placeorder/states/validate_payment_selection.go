package states

import (
	"context"

	"go.opencensus.io/trace"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// PaymentSelectionValidator decides if the PaymentSelection is valid
	PaymentSelectionValidator func(selection cart.PaymentSelection) error

	// ValidatePaymentSelection state
	ValidatePaymentSelection struct {
		validator PaymentSelectionValidator
	}
)

var _ process.State = ValidatePaymentSelection{}

// Inject dependencies
func (v *ValidatePaymentSelection) Inject(
	opts *struct {
		Validator PaymentSelectionValidator `inject:",optional"`
	},
) *ValidatePaymentSelection {
	if opts != nil {
		v.validator = opts.Validator
	}

	return v
}

// Name get state name
func (ValidatePaymentSelection) Name() string {
	return "ValidatePaymentSelection"
}

// Run the state operations
func (v ValidatePaymentSelection) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/ValidatePaymentSelection/Run")
	defer span.End()

	paymentSelection := p.Context().Cart.PaymentSelection
	if paymentSelection == nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{
				Error: "no payment selection on cart",
			},
		}
	}

	if v.validator != nil {
		err := v.validator(paymentSelection)
		if err != nil {
			return process.RunResult{
				Failed: process.ErrorOccurredReason{
					Error: err.Error(),
				},
			}
		}
	}

	p.UpdateState(CreatePayment{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (v ValidatePaymentSelection) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (v ValidatePaymentSelection) IsFinal() bool {
	return false
}
