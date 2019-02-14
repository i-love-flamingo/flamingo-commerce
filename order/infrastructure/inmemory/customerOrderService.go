package inmemory

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/v3/order/domain"
	authDomain "flamingo.me/flamingo/v3/core/auth/domain"
)

type (
	// CustomerOrderService defines the inmemory customer order service
	CustomerOrderService struct {
		behaviour *Behaviour
	}
)

var (
	_ domain.CustomerOrderService = (*CustomerOrderService)(nil)
)

// Inject dependencies
func (os *CustomerOrderService) Inject(
	Behaviour *Behaviour,
) {
	os.behaviour = Behaviour
}

// Get all orders depending on the authentication
func (os *CustomerOrderService) Get(context.Context, authDomain.Auth) ([]*domain.Order, error) {
	return nil, errors.New("not yet implemented")
}

// GetByID an order by ID depending on the authentication
func (os *CustomerOrderService) GetByID(context.Context, authDomain.Auth, string) (*domain.Order, error) {
	return nil, errors.New("not yet implemented")
}

// GetBehaviour the service behaviour
func (os *CustomerOrderService) GetBehaviour(context.Context, authDomain.Auth) (domain.Behaviour, error) {
	return os.behaviour, nil
}
