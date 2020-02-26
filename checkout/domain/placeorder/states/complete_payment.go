package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"go.opencensus.io/trace"
)

type (
	// CompletePayment state
	CompletePayment struct {
		paymentService *application.PaymentService
	}
)

var _ process.State = CompletePayment{}

// Inject dependencies
func (c *CompletePayment) Inject(
	paymentService *application.PaymentService,
) *CompletePayment {
	c.paymentService = paymentService

	return c
}

// Name get state name
func (CompletePayment) Name() string {
	return "CompletePayment"
}

// Run the state operations
func (c CompletePayment) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/CompletePayment/Run")
	defer span.End()

	cart := p.Context().Cart
	paymentGateway, err := c.paymentService.PaymentGatewayByCart(cart)
	if err != nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
		}
	}

	payment, err := paymentGateway.OrderPaymentFromFlow(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
		}
	}

	err = paymentGateway.ConfirmResult(ctx, &cart, payment)
	if err != nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateState(ValidatePayment{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (c CompletePayment) Rollback(ctx context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (c CompletePayment) IsFinal() bool {
	return false
}
