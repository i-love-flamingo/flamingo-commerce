package payment

import (
	"net/url"

	"github.com/pkg/errors"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/checkout/domain/payment"
	"go.aoe.com/flamingo/framework/web"
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
func (pa *OfflinePaymentProvider) RedirectExternalPayment(ctx web.Context, currentCart *cartDomain.Cart, method *payment.PaymentMethod, returnUrl *url.URL) (web.Response, error) {
	return nil, errors.New("No Redirect")
}

func (pa *OfflinePaymentProvider) IsActive() bool {
	return pa.Enabled
}

func (pa *OfflinePaymentProvider) ProcessPayment(ctx web.Context, currentCart *cartDomain.Cart, method *payment.PaymentMethod, _ map[string]string) (*cartDomain.CartPayment, error) {
	paymentInfo := cartDomain.PaymentInfo{
		Method:   method.Code,
		Provider: pa.GetCode(),
		Status:   cartDomain.PAYMENT_STATUS_OPEN,
	}

	idAssignments := make(map[string]*cartDomain.PaymentInfo)
	for _, itemId := range currentCart.GetItemIds() {
		idAssignments[itemId] = &paymentInfo
	}
	cartPayment := cartDomain.CartPayment{
		PaymentInfos:     []cartDomain.PaymentInfo{paymentInfo},
		ItemIDAssignment: idAssignments,
	}

	return &cartPayment, nil
}
