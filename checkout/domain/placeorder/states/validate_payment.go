package states

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// ValidatePayment state
	ValidatePayment struct {
		paymentGateway interfaces.WebCartPaymentGateway
	}
)

var _ process.State = ValidatePayment{}

func init() {
	gob.Register(ValidatePayment{})
}

// Name get state name
func (ValidatePayment) Name() string {
	return "ValidatePayment"
}

// Run the state operations
func (v ValidatePayment) Run(ctx context.Context, p *process.Process) process.RunResult {
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
