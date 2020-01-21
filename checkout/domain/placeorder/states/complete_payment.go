package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// CompletePayment state
	CompletePayment struct {
		paymentGateway interfaces.WebCartPaymentGateway
	}
)

var _ process.State = CompletePayment{}

// Inject dependencies
func (c *CompletePayment) Inject(
	paymentGateway interfaces.WebCartPaymentGateway,
) *CompletePayment {
	c.paymentGateway = paymentGateway

	return c
}

// Name get state name
func (CompletePayment) Name() string {
	return "CompletePayment"
}

// Run the state operations
func (c CompletePayment) Run(ctx context.Context, p *process.Process) process.RunResult {
	cart := p.Context().Cart

	payment, err := c.paymentGateway.OrderPaymentFromFlow(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	err = c.paymentGateway.ConfirmResult(ctx, &cart, payment)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateState(ValidatePayment{}.Name())
	return process.RunResult{}
}

// Rollback the state operations
func (c CompletePayment) Rollback(_ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (c CompletePayment) IsFinal() bool {
	return false
}
