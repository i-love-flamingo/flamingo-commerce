package states

import (
	"context"
	"encoding/gob"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"go.opencensus.io/trace"
)

type (
	// ShowIframe state
	ShowIframe struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}
)

var _ process.State = ShowIframe{}

func init() {
	gob.Register(&url.URL{})
}

// NewShowIframeStateData creates new state data for this state
func NewShowIframeStateData(url *url.URL) process.StateData {
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
func (si ShowIframe) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/ShowIframe/Run")
	defer span.End()

	return si.validator(ctx, p, si.paymentService)
}

// Rollback the state operations
func (si ShowIframe) Rollback(ctx context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (si ShowIframe) IsFinal() bool {
	return false
}
