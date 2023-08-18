package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/payment/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name WebCartPaymentGateway --case snake

type (
	// WebCartPaymentGatewayProvider defines the map of providers for payment providers
	WebCartPaymentGatewayProvider func() map[string]WebCartPaymentGateway

	// WebCartPaymentGateway is an interface offering (online) payment service - most probably against a external payment gateway API
	WebCartPaymentGateway interface {
		// Methods returns the PaymentGateway available Payment Methods
		Methods() []domain.Method

		// StartFlow returns the data for a new flow
		StartFlow(ctx context.Context, cart *cart.Cart, correlationID string, returnURL *url.URL) (*domain.FlowResult, error)

		// FlowStatus returns the status of a previously started flow (see StartFlow())
		FlowStatus(ctx context.Context, cart *cart.Cart, correlationID string) (*domain.FlowStatus, error)

		// ConfirmResult used to finally confirm the result
		ConfirmResult(ctx context.Context, cart *cart.Cart, cartPayment *placeorder.Payment) error

		// OrderPaymentFromFlow generates a place order payment for a previously created flow
		OrderPaymentFromFlow(ctx context.Context, cart *cart.Cart, correlationID string) (*placeorder.Payment, error)

		// CancelOrderPayment cancels the place order payment
		CancelOrderPayment(ctx context.Context, cartPayment *placeorder.Payment) error
	}
)
