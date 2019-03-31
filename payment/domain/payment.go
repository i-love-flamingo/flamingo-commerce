package domain

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Method contains information about a general payment method
	Method struct {
		//A speaking title
		Title string
		//A unique Code
		Code string
	}

	// WebPaymentGateway is an interface offering (online) payment service - most probably against a external payment gateway API
	WebPaymentGateway interface {

		// Methods returns the PaymentGateway available Payment Methods
		Methods() []Method

		// StartFlow - starts the processing of an asyncron Payment Flow for the cart
		// @param correlationId - is used later to fetch the result of this processing request
		// @param returnUrl - this is the desired end of the async payment flow.
		// @return the web.Result need to be executed(returned) by the call to give control to the Flow
		StartFlow(ctx context.Context, cart *cart.Cart, selectedMethod Method, correlationId string, returnUrl *url.URL) (web.Result, error)

		// GetFlowResult - will be used to fetch the result of the payment flow
		// @param correlationId - is used to fetch the result of a processing request started by this correlationId
		GetFlowResult(ctx context.Context, cart *cart.Cart, correlationId string) (*cart.Payment, error)
	}

	// Error should be used by PaymentGateway to indicate that payment failed (so that the customer can see a speaking message)
	Error struct {
		ErrorMessage string
		ErrorCode    string
	}
)

const (
	// PaymentErrorCodeCancelled cancelled
	PaymentErrorCodeCancelled = "payment_cancelled"
	// PaymentErrorCodeAuthorizeFailed error
	PaymentErrorCodeAuthorizeFailed = "authorization_failed"
	// PaymentErrorCodeCaptureFailed error
	PaymentErrorCodeCaptureFailed = "capture_failed"
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
