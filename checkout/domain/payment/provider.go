package payment

import (
	"context"

	"go.aoe.com/flamingo/framework/web"
)

type (
	PaymentMethod struct {
		Title string
		IsExternalPayment bool
		ExternalRedirectUri string
	}

	PaymentProvider interface {
		// GetPaymentMethods returns the Payment Providers available Payment Methods
		GetPaymentMethods() []PaymentMethod
		// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
		RedirectExternalPayment (context.Context, PaymentMethod) (web.Response, error)
		IsActive() bool
	}
)
