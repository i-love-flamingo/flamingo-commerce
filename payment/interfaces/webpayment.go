package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/payment/domain"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/web"
)

type (

	// WebCartPaymentGatewayProvider defines the map of providers for payment providers
	WebCartPaymentGatewayProvider func() map[string]WebCartPaymentGateway

	// WebCartPaymentGateway is an interface offering (online) payment service - most probably against a external payment gateway API
	WebCartPaymentGateway interface {

		// Methods returns the PaymentGateway available Payment Methods
		Methods() []domain.Method

		// StartFlow - starts the processing of an asyncron Payment Flow for the cart
		// @param correlationID - is used later to fetch the result of this processing request
		// @param returnUrl - this is the desired end of the async payment flow.
		// @return the web.Result need to be executed(returned) by the call to give control to the Flow
		StartFlow(ctx context.Context, cart *cart.Cart, correlationID string, returnUrl *url.URL) (web.Result, error)

		// GetFlowResult - will be used to fetch the result of the payment flow
		// @param correlationID - is used to fetch the result of a processing request started by this correlationId
		GetFlowResult(ctx context.Context, cart *cart.Cart, correlationID string) (*cart.Payment, error)

		// ConfirmResult - used to finally confirm the result
		ConfirmResult(ctx context.Context, cart *cart.Cart, cartPayment *cart.Payment) error
	}
)
