package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
)

type (
	// ShowHTML state
	ShowHTML struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}
)

var _ process.State = ShowHTML{}

// NewShowHTMLStateData creates new StateData required for this ShowHTML state
func NewShowHTMLStateData(html string) process.StateData {
	return process.StateData(html)
}

// Inject dependencies
func (sh *ShowHTML) Inject(
	paymentService *application.PaymentService,
	validator process.PaymentValidatorFunc,
) *ShowHTML {
	sh.paymentService = paymentService
	sh.validator = validator
	return sh
}

// Name get state name
func (ShowHTML) Name() string {
	return "ShowHTML"
}

// Run the state operations
func (sh ShowHTML) Run(ctx context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	return sh.validator(ctx, p, sh.paymentService)
}

// Rollback the state operations
func (sh ShowHTML) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (sh ShowHTML) IsFinal() bool {
	return false
}
