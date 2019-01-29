package fake

import (
	"context"
	"time"

	"flamingo.me/flamingo-commerce/order/infrastructure/inmemory"

	"flamingo.me/flamingo-commerce/order/domain"
	authDomain "flamingo.me/flamingo/core/auth/domain"
)

type (
	// CustomerOrders is the fake customer orders api service
	CustomerOrders struct{} 
)

var (
	_ domain.CustomerOrderService = (*CustomerOrders)(nil)
)

// GetBehaviour returns the InMemoryBehaviour
func (co *CustomerOrders) GetBehaviour(context.Context, authDomain.Auth) (domain.Behaviour, error) {
	return new(inmemory.Behaviour), nil
}

// Get returns a CustomerOrders struct
func (co *CustomerOrders) Get(ctx context.Context, authentication authDomain.Auth) ([]*domain.Order, error) {
	return []*domain.Order{
		{
			ID: "100",

			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		}, {
			ID: "100",

			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		}, {
			ID: "100",

			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		}, {
			ID: "100",

			CreationTime: time.Now(),
			UpdateTime:   time.Now(),
			Status:       "Teststatus",
			OrderItems:   make([]*domain.OrderItem, 0),
			Total:        123.45,
			CurrencyCode: "EUR",
		},
	}, nil
}

// GetByID fetches a faked customer order by id
func (co *CustomerOrders) GetByID(ctx context.Context, authentication authDomain.Auth, id string) (*domain.Order, error) {
	return &domain.Order{
		ID:           "100",
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
		Status:       "Teststatus",
		OrderItems:   make([]*domain.OrderItem, 0),
		Total:        123.45,
		CurrencyCode: "EUR",
	}, nil
}
