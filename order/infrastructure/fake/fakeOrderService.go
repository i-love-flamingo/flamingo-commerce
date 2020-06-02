package fake

import (
	"context"
	"time"

	"flamingo.me/flamingo-commerce/v3/order/domain"
	"flamingo.me/flamingo/v3/core/auth"
)

type (
	// CustomerOrders is the fake customer order adapter
	CustomerOrders struct{}
)

var (
	_ domain.CustomerIdentityOrderService = (*CustomerOrders)(nil)
)

// Get all orders for a customer
func (co *CustomerOrders) Get(_ context.Context, _ auth.Identity) ([]*domain.Order, error) {
	return []*domain.Order{
		{
			ID:           "100",
			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		}, {
			ID:           "101",
			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		}, {
			ID:           "102",
			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		}, {
			ID:           "103",
			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		},
	}, nil
}

// GetByID returns a single customer order
func (co *CustomerOrders) GetByID(_ context.Context, _ auth.Identity, orderID string) (*domain.Order, error) {
	return &domain.Order{
		ID:           orderID,
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		Status:       "Teststatus",
		OrderItems:   make([]*domain.OrderItem, 0),
		Total:        123.45,
		CurrencyCode: "EUR",
	}, nil
}
