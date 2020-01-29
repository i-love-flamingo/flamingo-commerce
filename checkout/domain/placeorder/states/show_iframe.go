package states

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
)

type (
	// ShowIframe state
	ShowIframe struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}
)

var _ process.State = ShowIframe{}

// NewShowIframeStateData creates new state data for this state
func NewShowIframeStateData(url url.URL) process.StateData {
	return process.StateData(url)
}

// Inject dependencies
func (si *ShowIframe) Inject(
	paymentService *application.PaymentService,
	validator process.PaymentValidatorFunc,
) *ShowIframe {
	si.paymentService = paymentService
	si.validator = validator

	return si
}

// Name get state name
func (ShowIframe) Name() string {
	return "ShowIframe"
}

// Run the state operations
func (si ShowIframe) Run(ctx context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	return si.validator(ctx, p, si.paymentService)
}

// Rollback the state operations
func (si ShowIframe) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (si ShowIframe) IsFinal() bool {
	return false
}
