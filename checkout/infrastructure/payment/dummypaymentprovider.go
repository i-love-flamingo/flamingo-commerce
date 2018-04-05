package payment

import (
	"net/url"

	"github.com/pkg/errors"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
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
func (pa *DummyPaymentProvider) RedirectExternalPayment(ctx web.Context, method *payment.PaymentMethod, returnUrl *url.URL) (web.Response, error) {
	return nil, errors.New("Only a Dummy Adapter")
}

func (pa *DummyPaymentProvider) IsActive() bool {
	return false
}

func (pa *DummyPaymentProvider) ProcessPayment(ctx web.Context, method *payment.PaymentMethod, _ map[string]string) (*cartDomain.CartPayment, error) {
	return nil, nil
}
