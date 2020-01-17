package states

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// PlaceOrder state
	PlaceOrder struct {
		paymentGateway interfaces.WebCartPaymentGateway
	}
)

var _ process.State = PlaceOrder{}

func init() {
	gob.Register(PlaceOrder{})
}

// Inject dependencies
func (po *PlaceOrder) Inject(
	paymentGateway interfaces.WebCartPaymentGateway,
) *PlaceOrder {
	po.paymentGateway = paymentGateway

	return po
}

// Name get state name
func (PlaceOrder) Name() string {
	return "CreatePayment"
}

// Run the state operations
func (po PlaceOrder) Run(ctx context.Context, p *process.Process) process.RunResult {
	return process.RunResult{}
}

// Rollback the state operations
func (po PlaceOrder) Rollback(data process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (po PlaceOrder) IsFinal() bool {
	return false
}
