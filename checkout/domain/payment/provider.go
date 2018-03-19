package payment

import (
	"go.aoe.com/flamingo/framework/web"
)

type (
	PaymentMethod struct {
		Title string
		IsExternalPayment bool
		ExternalRedirectUri string
	}

	PaymentProvider interface {
		GetPaymentMethods() []PaymentMethod
		RedirectExternalPayment (PaymentMethod) web.Response
	}
)
