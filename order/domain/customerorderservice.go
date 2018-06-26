package domain

import (
	"context"

	"flamingo.me/flamingo/core/auth/domain"
)

// CustomerOrderService for customer order retrieval
type CustomerOrderService interface {
	Get(context.Context, domain.Auth) ([]*Order, error)
	GetById(context.Context, domain.Auth, string) (*Order, error)
}
