package placeorder

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

// PaymentValidator to decide over the next state
func PaymentValidator(ctx context.Context, p *process.Process, paymentService *application.PaymentService) process.RunResult {
	cart := p.Context().Cart
	gateway, err := paymentService.PaymentGatewayByCart(p.Context().Cart)
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
			formFields := make(map[string]states.FormField, len(flowStatus.ActionData.FormParameter))
			for k, v := range flowStatus.ActionData.FormParameter {
				formFields[k] = states.FormField{
					Value: v.Value,
				}
			}
			if flowStatus.ActionData.URL == nil {
				p.Failed(ctx, process.ErrorOccurredReason{Error: "no redirect url set for action"})
			}
			p.UpdateState(states.PostRedirect{}.Name(), states.NewPostRedirectStateData(*flowStatus.ActionData.URL, formFields))
		case paymentDomain.PaymentFlowActionRedirect:
			if flowStatus.ActionData.URL == nil {
				p.Failed(ctx, process.ErrorOccurredReason{Error: "no redirect url set for action"})
			}
			p.UpdateState(states.Redirect{}.Name(), states.NewRedirectStateData(*flowStatus.ActionData.URL))
		case paymentDomain.PaymentFlowActionShowHTML:
			p.UpdateState(states.ShowHTML{}.Name(), states.NewShowHTMLStateData(flowStatus.ActionData.DisplayData))
		case paymentDomain.PaymentFlowActionShowIFrame:
			if flowStatus.ActionData.URL == nil {
				p.Failed(ctx, process.ErrorOccurredReason{Error: "no redirect url set for action"})
			}
			p.UpdateState(states.ShowIframe{}.Name(), states.NewShowIframeStateData(*flowStatus.ActionData.URL))
		default:
			p.Failed(ctx, process.PaymentErrorOccurredReason{
				Error: fmt.Sprintf("Payment action not supported: %q", flowStatus.Action),
			})
		}
	case paymentDomain.PaymentFlowStatusApproved:
		// payment is done but needs confirmation
		p.UpdateState(states.CompletePayment{}.Name(), nil)
	case paymentDomain.PaymentFlowStatusCompleted:
		// payment is done and confirmed, place order if not already placed
		p.UpdateState(states.Success{}.Name(), nil)
	case paymentDomain.PaymentFlowStatusAborted, paymentDomain.PaymentFlowStatusFailed, paymentDomain.PaymentFlowStatusCancelled:
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: flowStatus.Status},
		}
	case paymentDomain.PaymentFlowWaitingForCustomer:
		// payment pending, waiting for customer doing async stuff
		// todo: add new state representing customer wait
		p.UpdateState(states.Wait{}.Name(), nil)
	default:
		// unknown payment flowStatus, let frontend handle it
		p.UpdateState(states.Failed{}.Name(), nil)
	}

	return process.RunResult{}
}
