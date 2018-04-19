package payment

import (
	"net/url"

	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
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
		RedirectExternalPayment(web.Context, *cartDomain.Cart, *PaymentMethod, *url.URL) (web.Response, error)
		// ProcessPayment, map is for form Data, payment Data, etc - whatever the Payment Method requires
		ProcessPayment(web.Context, *cartDomain.Cart, *PaymentMethod, map[string]string) (*cartDomain.CartPayment, error)
		IsActive() bool
	}
)
