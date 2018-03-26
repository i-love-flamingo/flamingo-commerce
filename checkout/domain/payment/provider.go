package payment

import (
	"go.aoe.com/flamingo/framework/web"
)

type (
	PaymentMethod struct {
		Title string
		Code string
		IsExternalPayment bool
		ExternalRedirectUri string
	}

	PaymentProvider interface {
		// GetPaymentMethods returns the Payment Providers available Payment Methods
		GetPaymentMethods() []PaymentMethod
		// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
		RedirectExternalPayment (web.Context, PaymentMethod) (web.Response, error)
		IsActive() bool
	}
)
