package states

import (
	"context"
	"encoding/gob"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"go.opencensus.io/trace"
)

type (
	// CreatePayment state
	CreatePayment struct {
		paymentService *application.PaymentService
	}

	// CreatePaymentRollbackData needed for rollback
	CreatePaymentRollbackData struct {
		PaymentID          string
		Gateway            string
		RawTransactionData interface{}
	}
)

var _ process.State = CreatePayment{}

func init() {
	gob.Register(CreatePaymentRollbackData{})
}

// Inject dependencies
func (c *CreatePayment) Inject(
	paymentService *application.PaymentService,
) *CreatePayment {
	c.paymentService = paymentService

	return c
}

// Name get state name
func (CreatePayment) Name() string {
	return "CreatePayment"
}

// Run the state operations
func (c CreatePayment) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/CreatePayment/Run")
	defer span.End()

	cart := p.Context().Cart
	paymentGateway, err := c.paymentService.PaymentGatewayByCart(cart)
	if err != nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
		}
	}

	_, err = paymentGateway.StartFlow(ctx, &cart, p.Context().UUID, p.Context().ReturnURL)
	if err != nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
		}
	}

	payment, err := paymentGateway.OrderPaymentFromFlow(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateState(CompleteCart{}.Name(), nil)
	return process.RunResult{
		RollbackData: CreatePaymentRollbackData{
			PaymentID:          payment.PaymentID,
			Gateway:            payment.Gateway,
			RawTransactionData: payment.RawTransactionData,
		},
	}
}

// Rollback the state operations
func (c CreatePayment) Rollback(ctx context.Context, data process.RollbackData) error {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/CreatePayment/Rollback")
	defer span.End()

	rollbackData, ok := data.(CreatePaymentRollbackData)
	if !ok {
		return fmt.Errorf("rollback data not of expected type 'CreatePaymentRollbackData', but %T", rollbackData)
	}

	paymentGateway, err := c.paymentService.PaymentGateway(rollbackData.Gateway)
	if err != nil {
		return err
	}

	err = paymentGateway.CancelOrderPayment(
		ctx,
		&placeorder.Payment{
			Gateway:            rollbackData.Gateway,
			PaymentID:          rollbackData.PaymentID,
			RawTransactionData: rollbackData.RawTransactionData,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// IsFinal if state is a final state
func (c CreatePayment) IsFinal() bool {
	return false
}
