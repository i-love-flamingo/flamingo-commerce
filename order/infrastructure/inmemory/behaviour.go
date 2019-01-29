package inmemory

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/order/domain"
)

type (
	// Behaviour defines the in memory order behaviour
	Behaviour struct {
		storage Storager
	}

	// Storager interface for in memory order storage
	Storager interface {
		GetOrder(id string) (*domain.Order, error)
		HasOrder(id string) bool
		StoreOrder(order *domain.Order) error
	}

	// Storage as default implementation of an order storage
	Storage struct {
		orders map[string]*domain.Order
	}
)

var (
	_ domain.Behaviour = (*Behaviour)(nil)
)

// Inject dependencies
func (imb *Behaviour) Inject(
	OrderStorage Storager,
) {
	imb.storage = OrderStorage
}

// PlaceOrder handles the in memory order service
func (imb *Behaviour) PlaceOrder(ctx context.Context, cart *cart.Cart, payment *cart.CartPayment) (domain.PlacedOrderInfos, error) {
	return nil, errors.New("not yet implemented")
}

var (
	_ Storager = (*Storage)(nil)
)

// init the in memory order storage if required
func (os *Storage) init() {
	if os.orders == nil {
		os.orders = make(map[string]*domain.Order)
	}
}

// GetOrder gets an order from the storage
func (os *Storage) GetOrder(id string) (*domain.Order, error) {
	os.init()
	if !os.HasOrder(id) {
		return nil, errors.New("no such order")
	}

	result := os.orders[id]

	return result, nil
}

// HasOrder checks if an order with `id` is in the storage
func (os *Storage) HasOrder(id string) bool {
	os.init()
	_, result := os.orders[id]

	return result
}

// StoreOrder puts an order into the in memory order storage
func (os *Storage) StoreOrder(order *domain.Order) error {
	os.init()
	os.orders[order.ID] = order

	return nil
}
