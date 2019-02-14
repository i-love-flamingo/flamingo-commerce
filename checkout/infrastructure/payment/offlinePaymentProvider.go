package payment

import (
	"context"
	"net/url"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/payment"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	OfflinePaymentProvider struct {
		Enabled bool `inject:"config:checkout.enableOfflinePaymentProvider,optional"`
	}
)

var (
	_ payment.PaymentProvider = &OfflinePaymentProvider{}
)

func (pa *OfflinePaymentProvider) GetCode() string {
	return "offlinepayment"
}

// GetPaymentMethods returns the Payment Providers available Payment Methods
func (pa *OfflinePaymentProvider) GetPaymentMethods() []payment.PaymentMethod {
	var result []payment.PaymentMethod
	result = append(result, payment.PaymentMethod{
		Title:             "Cash on delivery",
		Code:              "offlinepayment_cashondelivery",
		IsExternalPayment: false,
	})
	return result
}

// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
func (pa *OfflinePaymentProvider) RedirectExternalPayment(ctx context.Context, r *web.Request, currentCart *cartDomain.Cart, method *payment.PaymentMethod, returnUrl *url.URL) (web.Result, error) {
	return nil, errors.New("No Redirect")
}

func (pa *OfflinePaymentProvider) IsActive() bool {
	return pa.Enabled
}

func (pa *OfflinePaymentProvider) ProcessPayment(ctx context.Context, r *web.Request, currentCart *cartDomain.Cart, method *payment.PaymentMethod, _ map[string]string) (*cartDomain.CartPayment, error) {
	paymentInfo := cartDomain.PaymentInfo{
		Method:   method.Code,
		Provider: pa.GetCode(),
		Status:   cartDomain.PAYMENT_STATUS_OPEN,
	}

	var assignments []cartDomain.CartPaymentAssignment
	for _, itemReference := range currentCart.GetItemCartReferences() {
		assignments = append(assignments, cartDomain.CartPaymentAssignment{
			ItemCartReference: itemReference,
			PaymentInfo:       &paymentInfo,
		})
	}
	var paymentInfos []*cartDomain.PaymentInfo
	paymentInfos = append(paymentInfos, &paymentInfo)

	cartPayment := cartDomain.CartPayment{
		PaymentInfos: paymentInfos,
		Assignments:  assignments,
	}
	return &cartPayment, nil
}
