package states

import (
	"context"
	"encoding/gob"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
)

type (
	// CreatePayment state
	CreatePayment struct {
		paymentGateway interfaces.WebCartPaymentGateway
	}

	// CreatePaymentRollbackData needed for rollback
	CreatePaymentRollbackData struct {
		Payment *placeorder.Payment
	}
)

var _ process.State = CreatePayment{}

func init() {
	gob.Register(CreatePayment{})
}

// Inject dependencies
func (c *CreatePayment) Inject(
	paymentGateway interfaces.WebCartPaymentGateway,
) *CreatePayment {
	c.paymentGateway = paymentGateway

	return c
}

// Name get state name
func (CreatePayment) Name() string {
	return "CreatePayment"
}

// Run the state operations
func (c CreatePayment) Run(ctx context.Context, p *process.Process) process.RunResult {
	cart := p.Context().Cart
	flowResult, err := c.paymentGateway.StartFlow(ctx, &cart, p.Context().UUID, p.Context().ReturnURL)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}
	payment, err := c.paymentGateway.OrderPaymentFromFlow(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}
	result := process.RunResult{
		RollbackData: CreatePaymentRollbackData{Payment: payment},
	}

	if flowResult.EarlyPlaceOrder {
		p.UpdateState(PlaceOrder{})
		return result
	}

	p.UpdateState(ValidatePayment{})
	return result
}

// Rollback the state operations
func (c CreatePayment) Rollback(data process.RollbackData) error {
	rollbackData, ok := data.(CreatePaymentRollbackData)
	if !ok {
		return fmt.Errorf("rollback data not of expected type 'CreatePaymentRollbackData', but %T", rollbackData)
	}
	err := c.paymentGateway.CancelOrderPayment(context.Background(), rollbackData.Payment)
	if err != nil {
		return err
	}

	return nil
}

// IsFinal if state is a final state
func (c CreatePayment) IsFinal() bool {
	return false
}
