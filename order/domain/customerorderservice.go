package domain

import (
	"context"

	"go.aoe.com/flamingo/core/auth/domain"
)

// CustomerOrderService for customer order retrieval
type CustomerOrderService interface {
	Get(context.Context, domain.Auth) ([]*Order, error)
}
