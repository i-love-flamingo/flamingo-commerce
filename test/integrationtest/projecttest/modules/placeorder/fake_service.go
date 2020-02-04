package placeorder

import (
	"context"
	"errors"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	authDomain "flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	// FakeAdapter provides fake place order adapter
	FakeAdapter struct {
		placedOrders map[string]placeorder.PlacedOrderInfos
	}
)

var (
	_ placeorder.Service = new(FakeAdapter)
)

// Inject dependencies
func (f *FakeAdapter) Inject() *FakeAdapter {
	f.placedOrders = make(map[string]placeorder.PlacedOrderInfos)
	return f
}

// PlaceGuestCart places a guest cart order
func (f *FakeAdapter) PlaceGuestCart(ctx context.Context, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return f.placeCart(cart)
}

// PlaceCustomerCart places a customer cart
func (f *FakeAdapter) PlaceCustomerCart(ctx context.Context, auth authDomain.Auth, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return f.placeCart(cart)
}

func (f *FakeAdapter) placeCart(cart *cartDomain.Cart) (placeorder.PlacedOrderInfos, error) {
	reservedID := cart.AdditionalData.ReservedOrderID

	_, found := f.placedOrders[reservedID]

	if found {
		return nil, errors.New("Another order with #" + reservedID + " already placed")
	}

	var placedOrders placeorder.PlacedOrderInfos
	placedOrders = append(placedOrders, placeorder.PlacedOrderInfo{
		OrderNumber: reservedID,
	})

	f.placedOrders[reservedID] = placedOrders

	return placedOrders, nil
}

// ReserveOrderID returns the reserved order id
func (f *FakeAdapter) ReserveOrderID(ctx context.Context, cart *cartDomain.Cart) (string, error) {
	return cart.ID, nil
}

// CancelGuestOrder cancels a guest order
func (f *FakeAdapter) CancelGuestOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos) error {
	return f.cancelOrder(orderInfos)
}

// CancelCustomerOrder cancels a customer order
func (f *FakeAdapter) CancelCustomerOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos, auth authDomain.Auth) error {
	return f.cancelOrder(orderInfos)
}

func (f *FakeAdapter) cancelOrder(orderInfos placeorder.PlacedOrderInfos) error {
	var toDelete []string
	for _, order := range orderInfos {
		_, found := f.placedOrders[order.OrderNumber]

		if !found {
			return errors.New("Order cancel not possible order #" + order.OrderNumber + " wasn't placed")
		}

		toDelete = append(toDelete, order.OrderNumber)
	}

	for _, orderNumber := range toDelete {
		delete(f.placedOrders, orderNumber)
	}

	return nil
}
