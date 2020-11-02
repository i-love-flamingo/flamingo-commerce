package states

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/application"
	"flamingo.me/flamingo-commerce/v3/payment/domain"

	"go.opencensus.io/trace"
)

type (
	// ShowWalletPayment state
	ShowWalletPayment struct {
		paymentService *application.PaymentService
		validator      process.PaymentValidatorFunc
	}

	// ShowWalletPaymentData holds details regarding the wallet payment
	ShowWalletPaymentData domain.WalletDetails
)

func init() {
	gob.Register(ShowWalletPaymentData{})
}

var _ process.State = ShowWalletPayment{}

// NewShowWalletPaymentStateData creates new StateData with (persisted) Data required for this state
func NewShowWalletPaymentStateData(walletDetails ShowWalletPaymentData) process.StateData {
	return process.StateData(walletDetails)
}

// Inject dependencies
func (pr *ShowWalletPayment) Inject(
	paymentService *application.PaymentService,
	validator process.PaymentValidatorFunc,
) *ShowWalletPayment {
	pr.paymentService = paymentService
	pr.validator = validator

	return pr
}

// Name get state name
func (ShowWalletPayment) Name() string {
	return "ShowWalletPayment"
}

// Run the state operations
func (pr ShowWalletPayment) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/ShowWalletPayment/Run")
	defer span.End()

	return pr.validator(ctx, p, pr.paymentService)
}

// Rollback the state operations
func (pr ShowWalletPayment) Rollback(ctx context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (pr ShowWalletPayment) IsFinal() bool {
	return false
}
