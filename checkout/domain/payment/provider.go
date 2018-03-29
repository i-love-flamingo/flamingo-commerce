package payment

import (
	"net/url"

	"go.aoe.com/flamingo/framework/web"
)

type (
	PaymentMethod struct {
		Title               string
		Code                string
		IsExternalPayment   bool
		ExternalRedirectUri string
	}

	PaymentProvider interface {
		GetCode() string
		// GetPaymentMethods returns the Payment Providers available Payment Methods
		GetPaymentMethods() []PaymentMethod
		// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
		RedirectExternalPayment(web.Context, *PaymentMethod, *url.URL) (web.Response, error)

		ProcessPayment(web.Context, *PaymentMethod) (bool, error, interface {})
		IsActive() bool
	}
)
