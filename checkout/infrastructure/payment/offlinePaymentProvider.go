package payment

import (
	"context"
	"net/url"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/payment"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

// OfflinePaymentProvider provides an offline payment integration
type OfflinePaymentProvider struct {
	Enabled bool `inject:"config:checkout.enableOfflinePaymentProvider,optional"`
}

var _ payment.Provider = &OfflinePaymentProvider{}

// GetCode for payment
func (pa *OfflinePaymentProvider) GetCode() string {
	return "offlinepayment"
}

// GetPaymentMethods returns the Payment Providers available Payment Methods
func (pa *OfflinePaymentProvider) GetPaymentMethods() []payment.Method {
	return []payment.Method{{
		Title:             "Cash on delivery",
		Code:              "offlinepayment_cashondelivery",
		IsExternalPayment: false,
	}}
}

// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
func (pa *OfflinePaymentProvider) RedirectExternalPayment(ctx context.Context, r *web.Request, currentCart *cartDomain.Cart, method *payment.Method, returnURL *url.URL) (web.Result, error) {
	return nil, errors.New("No Redirect")
}

// IsActive check
func (pa *OfflinePaymentProvider) IsActive() bool {
	return pa.Enabled
}

// ProcessPayment creates an offline payment
func (pa *OfflinePaymentProvider) ProcessPayment(ctx context.Context, r *web.Request, currentCart *cartDomain.Cart, method *payment.Method, _ map[string]string) (*cartDomain.Payment, error) {
	paymentInfo := cartDomain.PaymentInfo{
		Method:   method.Code,
		Provider: pa.GetCode(),
		Status:   cartDomain.PaymentStatusOpen,
	}

	var assignments []cartDomain.PaymentAssignment
	for _, itemReference := range currentCart.GetItemCartReferences() {
		assignments = append(assignments, cartDomain.PaymentAssignment{
			ItemCartReference: itemReference,
			PaymentInfo:       &paymentInfo,
		})
	}
	var paymentInfos []*cartDomain.PaymentInfo
	paymentInfos = append(paymentInfos, &paymentInfo)

	cartPayment := cartDomain.Payment{
		PaymentInfos: paymentInfos,
		Assignments:  assignments,
	}
	return &cartPayment, nil
}
