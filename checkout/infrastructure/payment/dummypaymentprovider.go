package payment

import (
	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/checkout/domain/payment"
	"go.aoe.com/flamingo/framework/web"
)

type (
	DummyPaymentProvider struct{}
)

func (pa *DummyPaymentProvider) GetCode() string {
	return "dummy"
}

// GetPaymentMethods returns the Payment Providers available Payment Methods
func (pa *DummyPaymentProvider) GetPaymentMethods() []payment.PaymentMethod {
	var result []payment.PaymentMethod

	return result
}

// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
func (pa *DummyPaymentProvider) RedirectExternalPayment(ctx web.Context, method *payment.PaymentMethod) (web.Response, error) {
	return nil, errors.New("Only a Dummy Adapter")
}

func (pa *DummyPaymentProvider) IsActive() bool {
	return false
}
