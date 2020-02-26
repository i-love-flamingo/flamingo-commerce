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
	// PostRedirect state
	PostRedirect struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}

	// PostRedirectData holds details regarding the redirect
	PostRedirectData struct {
		FormFields map[string]FormField
		URL        *url.URL
	}

	// FormField represents a form field to be displayed to the user
	FormField struct {
		Value []string
	}
)

func init() {
	gob.Register(PostRedirectData{})
}

var _ process.State = PostRedirect{}

// NewPostRedirectStateData creates new StateData with (persisted) Data required for this state
func NewPostRedirectStateData(url *url.URL, formParameter map[string]FormField) process.StateData {
	return process.StateData(PostRedirectData{
		FormFields: formParameter,
		URL:        url,
	})
}

// Inject dependencies
func (pr *PostRedirect) Inject(
	paymentService *application.PaymentService,
	validator process.PaymentValidatorFunc,
) *PostRedirect {
	pr.paymentService = paymentService
	pr.validator = validator

	return pr
}

// Name get state name
func (PostRedirect) Name() string {
	return "PostRedirect"
}

// Run the state operations
func (pr PostRedirect) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/PostRedirect/Run")
	defer span.End()

	return pr.validator(ctx, p, pr.paymentService)
}

// Rollback the state operations
func (pr PostRedirect) Rollback(ctx context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (pr PostRedirect) IsFinal() bool {
	return false
}
