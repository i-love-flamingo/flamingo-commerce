package states

import (
	"context"

	"go.opencensus.io/trace"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// ValidatePaymentSelection state
	ValidatePaymentSelection struct {
		validator            validation.PaymentSelectionValidator
		cartDecoratorFactory *decorator.DecoratedCartFactory
	}
)

var _ process.State = ValidatePaymentSelection{}

// Inject dependencies
func (v *ValidatePaymentSelection) Inject(
	cartDecoratorFactory *decorator.DecoratedCartFactory,
	opts *struct {
		Validator validation.PaymentSelectionValidator `inject:",optional"`
	},
) *ValidatePaymentSelection {
	v.cartDecoratorFactory = cartDecoratorFactory
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
	_, span := trace.StartSpan(ctx, "placeorder/state/ValidatePaymentSelection/Run")
	defer span.End()

	paymentSelection := p.Context().Cart.PaymentSelection
	if paymentSelection == nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{
				Error: cart.ErrPaymentSelectionNotSet.Error(),
			},
		}
	}

	if v.validator != nil {
		decoratedCart := v.cartDecoratorFactory.Create(ctx, p.Context().Cart)
		err := v.validator.Validate(ctx, decoratedCart, paymentSelection)
		if err != nil {
			return process.RunResult{
				Failed: process.PaymentErrorOccurredReason{
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
