package email

import (
	"context"
	"errors"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	authDomain "flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// PlaceOrderServiceAdapter provides an implementation of the Service as email adapter
	//  TODO - this example adapter need to be implemented
	PlaceOrderServiceAdapter struct {
		emailAddress string
		logger       flamingo.Logger
	}
)

var (
	_ placeorder.Service = new(PlaceOrderServiceAdapter)
)

// Inject dependencies
func (e *PlaceOrderServiceAdapter) Inject(logger flamingo.Logger, config *struct {
	EmailAddress string `inject:"config:commerce.cart.emailAdapter.emailAddress"`
}) {
	e.emailAddress = config.EmailAddress
	e.logger = logger.WithField("module", "cart").WithField("category", "emailAdapter")
}

// PlaceGuestCart places a guest cart as order email
func (e *PlaceOrderServiceAdapter) PlaceGuestCart(ctx context.Context, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	if payment == nil && cart.GrandTotal().IsPositive() {
		return nil, errors.New("No valid Payment given")
	}
	if cart.GrandTotal().IsPositive() {
		totalPrice, err := payment.TotalValue()
		if err != nil {
			return nil, err
		}
		if !totalPrice.Equal(cart.GrandTotal()) {
			return nil, errors.New("Payment Total does not match with Grandtotal")
		}
	}

	var placedOrders placeorder.PlacedOrderInfos
	placedOrders = append(placedOrders, placeorder.PlacedOrderInfo{
		OrderNumber: cart.ID,
	})

	return placedOrders, nil
}

// PlaceCustomerCart places a customer cart as order email
func (e *PlaceOrderServiceAdapter) PlaceCustomerCart(ctx context.Context, auth authDomain.Auth, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	var placedOrders placeorder.PlacedOrderInfos
	placedOrders = append(placedOrders, placeorder.PlacedOrderInfo{
		OrderNumber: cart.ID,
	})

	return placedOrders, nil
}

// ReserveOrderID returns the reserved order id
func (e *PlaceOrderServiceAdapter) ReserveOrderID(ctx context.Context, cart *cartDomain.Cart) (string, error) {
	return cart.ID, nil
}

// CancelGuestOrder cancels a guest order
func (e *PlaceOrderServiceAdapter) CancelGuestOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos) error {
	// since we don't actual place orders we just return nil here
	return nil
}

// CancelCustomerOrder cancels a customer order
func (e *PlaceOrderServiceAdapter) CancelCustomerOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos, auth authDomain.Auth) error {
	// since we don't actual place orders we just return nil here
	return nil
}
