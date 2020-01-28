package states

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

type (
	// ValidatePayment state
	ValidatePayment struct {
		paymentService *application.PaymentService
	}
)

var _ process.State = ValidatePayment{}

// Inject dependencies
func (v *ValidatePayment) Inject(
	paymentService *application.PaymentService,
) *ValidatePayment {
	v.paymentService = paymentService

	return v
}

// Name get state name
func (ValidatePayment) Name() string {
	return "ValidatePayment"
}

// Run the state operations
func (v ValidatePayment) Run(ctx context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	cart := p.Context().Cart
	gateway, err := v.paymentService.PaymentGatewayByCart(p.Context().Cart)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	flowStatus, err := gateway.FlowStatus(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	switch flowStatus.Status {
	case paymentDomain.PaymentFlowStatusUnapproved:
		switch flowStatus.Action {
		case paymentDomain.PaymentFlowActionPostRedirect:
			formFields := make(map[string]FormField, len(flowStatus.ActionData.FormParameter))
			for k, v := range flowStatus.ActionData.FormParameter {
				formFields[k] = FormField{
					Value: v.Value,
				}
			}
			if flowStatus.ActionData.URL == nil {
				p.Failed(ctx, process.ErrorOccurredReason{Error: "no redirect url set for action"})
			}
			p.UpdateState(PostRedirect{}.Name(), NewPostRedirectStateData(*flowStatus.ActionData.URL, formFields))
		case paymentDomain.PaymentFlowActionRedirect:
			if flowStatus.ActionData.URL == nil {
				p.Failed(ctx, process.ErrorOccurredReason{Error: "no redirect url set for action"})
			}
			p.UpdateState(Redirect{}.Name(), NewRedirectStateData(*flowStatus.ActionData.URL))
		case paymentDomain.PaymentFlowActionShowHTML:
			p.UpdateState(ShowHTML{}.Name(), NewShowHTMLStateData(flowStatus.ActionData.DisplayData))
		case paymentDomain.PaymentFlowActionShowIFrame:
			if flowStatus.ActionData.URL == nil {
				p.Failed(ctx, process.ErrorOccurredReason{Error: "no redirect url set for action"})
			}
			p.UpdateState(ShowIframe{}.Name(), NewShowIframeStateData(*flowStatus.ActionData.URL))
		default:
			p.Failed(ctx, process.PaymentErrorOccurredReason{
				Error: fmt.Sprintf("Payment action not supported: %q", flowStatus.Action),
			})
		}
	case paymentDomain.PaymentFlowStatusApproved:
		// payment is done but needs confirmation
		p.UpdateState(CompletePayment{}.Name(), nil)
	case paymentDomain.PaymentFlowStatusCompleted:
		// payment is done and confirmed, place order if not already placed
		p.UpdateState(Success{}.Name(), nil)
	case paymentDomain.PaymentFlowStatusAborted, paymentDomain.PaymentFlowStatusFailed, paymentDomain.PaymentFlowStatusCancelled:
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: flowStatus.Status},
		}
	case paymentDomain.PaymentFlowWaitingForCustomer:
		// payment pending, waiting for customer doing async stuff
		p.UpdateState(Wait{}.Name(), nil)
	default:
		// unknown payment flowStatus, let frontend handle it
		p.UpdateState(Failed{}.Name(), nil)
	}

	return process.RunResult{}
}

// Rollback the state operations
func (v ValidatePayment) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (v ValidatePayment) IsFinal() bool {
	return false
}
