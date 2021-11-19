package placeorder

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

const (
	// ValidatePaymentErrorNoActionURL used for errors when the needed URL is missing from the ActionData struct
	ValidatePaymentErrorNoActionURL = "no url set for action"
	// ValidatePaymentErrorNoWalletDetails used for errors when the needed URL is missing from the ActionData struct
	ValidatePaymentErrorNoWalletDetails = "no wallet details set for action"
	// ValidatePaymentErrorNoActionDisplayData used for errors when the needed DisplayData/HTML is missing from the ActionData struct
	ValidatePaymentErrorNoActionDisplayData = "no display data / html set for action"
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
				return process.RunResult{
					Failed: process.PaymentErrorOccurredReason{Error: ValidatePaymentErrorNoActionURL},
				}
			}
			p.UpdateState(states.PostRedirect{}.Name(), states.NewPostRedirectStateData(flowStatus.ActionData.URL, formFields))
		case paymentDomain.PaymentFlowActionShowWalletPayment:
			if flowStatus.ActionData.WalletDetails == nil {
				return process.RunResult{
					Failed: process.PaymentErrorOccurredReason{Error: ValidatePaymentErrorNoWalletDetails},
				}
			}
			p.UpdateState(states.ShowWalletPayment{}.Name(), states.NewShowWalletPaymentStateData(states.ShowWalletPaymentData(*flowStatus.ActionData.WalletDetails)))
		case paymentDomain.PaymentFlowActionRedirect:
			if flowStatus.ActionData.URL == nil {
				return process.RunResult{
					Failed: process.PaymentErrorOccurredReason{Error: ValidatePaymentErrorNoActionURL},
				}
			}
			p.UpdateState(states.Redirect{}.Name(), states.NewRedirectStateData(flowStatus.ActionData.URL))
		case paymentDomain.PaymentFlowActionShowHTML:
			if flowStatus.ActionData.DisplayData == "" {
				return process.RunResult{
					Failed: process.PaymentErrorOccurredReason{Error: ValidatePaymentErrorNoActionDisplayData},
				}
			}
			p.UpdateState(states.ShowHTML{}.Name(), states.NewShowHTMLStateData(flowStatus.ActionData.DisplayData))
		case paymentDomain.PaymentFlowActionShowIframe:
			if flowStatus.ActionData.URL == nil {
				return process.RunResult{
					Failed: process.PaymentErrorOccurredReason{Error: ValidatePaymentErrorNoActionURL},
				}
			}
			p.UpdateState(states.ShowIframe{}.Name(), states.NewShowIframeStateData(flowStatus.ActionData.URL))
		case paymentDomain.PaymentFlowActionTriggerClientSDK:
			if flowStatus.ActionData.URL == nil {
				return process.RunResult{
					Failed: process.PaymentErrorOccurredReason{Error: ValidatePaymentErrorNoActionURL},
				}
			}
			p.UpdateState(states.TriggerClientSDK{}.Name(), states.NewTriggerClientSDKStateData(flowStatus.ActionData.URL, flowStatus.ActionData.DisplayData))
		default:
			return process.RunResult{
				Failed: process.PaymentErrorOccurredReason{Error: fmt.Sprintf("Payment action not supported: %q", flowStatus.Action)},
			}
		}
	case paymentDomain.PaymentFlowStatusApproved:
		// payment is done but needs confirmation
		p.UpdateState(states.CompletePayment{}.Name(), nil)
	case paymentDomain.PaymentFlowStatusCompleted:
		// payment is done and confirmed, place order if not already placed
		p.UpdateState(states.Success{}.Name(), nil)
	case paymentDomain.PaymentFlowStatusFailed, paymentDomain.PaymentFlowStatusCancelled:
		err := ""
		if flowStatus.Error != nil {
			err = flowStatus.Error.Error()
		}
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err},
		}
	case paymentDomain.PaymentFlowStatusAborted:
		return process.RunResult{
			Failed: process.PaymentCanceledByCustomerReason{},
		}
	case paymentDomain.PaymentFlowWaitingForCustomer:
		// payment pending, waiting for customer doing async stuff like finishing is payment in mobile app
		p.UpdateState(states.WaitForCustomer{}.Name(), nil)
	default:
		// unknown payment flow status
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: fmt.Sprintf("Payment status not supported: %q", flowStatus.Status)},
		}
	}

	return process.RunResult{}
}
