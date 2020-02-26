package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
)

type (
	// WaitForCustomer state
	WaitForCustomer struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}
)

var _ process.State = WaitForCustomer{}

// Inject dependencies
func (wc *WaitForCustomer) Inject(
	paymentService *application.PaymentService,
	validator process.PaymentValidatorFunc,
) *WaitForCustomer {
	wc.paymentService = paymentService
	wc.validator = validator

	return wc
}

// Name get state name
func (WaitForCustomer) Name() string {
	return "WaitForCustomer"
}

// Run the state operations
func (wc WaitForCustomer) Run(ctx context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	return wc.validator(ctx, p, wc.paymentService)
}

// Rollback the state operations
func (wc WaitForCustomer) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (wc WaitForCustomer) IsFinal() bool {
	return false
}
