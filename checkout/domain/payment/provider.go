package payment

import (
	"context"
	"net/url"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Method contains information about a general payment method
	Method struct {
		Title               string
		Code                string
		IsExternalPayment   bool
		ExternalRedirectURI string
	}

	// Provider returns the payment methods
	Provider interface {
		GetCode() string
		// GetPaymentMethods returns the Payment Providers available Payment Methods
		GetPaymentMethods() []Method
		// RedirectExternalPayment starts a Redirect to an external Payment Page (if applicable)
		RedirectExternalPayment(context.Context, *web.Request, *cartDomain.Cart, *Method, *url.URL) (web.Result, error)
		// ProcessPayment, map is for form Data, payment Data, etc - whatever the Payment Method requires
		ProcessPayment(context.Context, *web.Request, *cartDomain.Cart, *Method, map[string]string) (*cartDomain.Payment, error)
		IsActive() bool
	}

	// Error should be used by PaymentProviders to indicate that payment failed (so that the customer can see a speaking message)
	Error struct {
		ErrorMessage string
		ErrorCode    string
	}
)

const (
	// PaymentCancelled cancelled
	PaymentCancelled = "payment_cancelled"
	// PaymentAuthorizeFailed error
	PaymentAuthorizeFailed = "authorization_failed"
	// PaymentCaptureFailed error
	PaymentCaptureFailed = "capture_failed"
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
