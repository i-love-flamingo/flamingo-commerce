package states

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// ValidatePayment state
	ValidatePayment struct {
		paymentGateway  interfaces.WebCartPaymentGateway
		EarlyPlaceOrder bool
	}
)

var _ process.State = ValidatePayment{}

func init() {
	gob.Register(ValidatePayment{})
}

// Inject dependencies
func (v *ValidatePayment) Inject(
	paymentGateway interfaces.WebCartPaymentGateway,
) *ValidatePayment {
	v.paymentGateway = paymentGateway

	return v
}

// Name get state name
func (ValidatePayment) Name() string {
	return "ValidatePayment"
}

// Run the state operations
func (v ValidatePayment) Run(ctx context.Context, p *process.Process) process.RunResult {
	cart := p.Context().Cart
	flowStatus, err := v.paymentGateway.FlowStatus(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	switch flowStatus.Status {
	case paymentDomain.PaymentFlowStatusUnapproved:
		// payment just started, frontend needs to do actions
		p.UpdateState(Wait{})
	case paymentDomain.PaymentFlowStatusApproved:
		// payment is done but needs confirmation
		p.UpdateState(CompletePayment{})
	case paymentDomain.PaymentFlowStatusCompleted:
		// payment is done and confirmed, place order if not already placed
		if v.EarlyPlaceOrder {
			p.UpdateState(Success{})
		} else {
			p.UpdateState(PlaceOrder{})
		}
	case paymentDomain.PaymentFlowStatusAborted, paymentDomain.PaymentFlowStatusFailed, paymentDomain.PaymentFlowStatusCancelled:
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: flowStatus.Status},
		}
	case paymentDomain.PaymentFlowWaitingForCustomer:
		// payment pending, waiting for customer doing async stuff
		p.UpdateState(Wait{})
	default:
		// unknown payment flowStatus, let frontend handle it
		p.UpdateState(Wait{})
	}

	return process.RunResult{}
}

// Rollback the state operations
func (v ValidatePayment) Rollback(data process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (v ValidatePayment) IsFinal() bool {
	return false
}
