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
	// TriggerClientSDK state
	TriggerClientSDK struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}

	// TriggerClientSDKData holds the data which must be sent to the client via SDK
	TriggerClientSDKData struct {
		URL  *url.URL
		Data string
	}
)

var _ process.State = TriggerClientSDK{}

func init() {
	gob.Register(TriggerClientSDKData{})
}

// NewTriggerClientSDKStateData creates data required for this state
func NewTriggerClientSDKStateData(url *url.URL, data string) process.StateData {
	return process.StateData(TriggerClientSDKData{
		URL:  url,
		Data: data,
	})
}

// Inject dependencies
func (r *TriggerClientSDK) Inject(
	paymentService *application.PaymentService,
	validator process.PaymentValidatorFunc,
) *TriggerClientSDK {
	r.paymentService = paymentService
	r.validator = validator

	return r
}

// Name get state name
func (TriggerClientSDK) Name() string {
	return "TriggerClientSDK"
}

// Run the state operations
func (r TriggerClientSDK) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/TriggerClientSDK/Run")
	defer span.End()

	return r.validator(ctx, p, r.paymentService)
}

// Rollback the state operations
func (r TriggerClientSDK) Rollback(_ context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (r TriggerClientSDK) IsFinal() bool {
	return false
}
