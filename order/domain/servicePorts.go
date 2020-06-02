package domain

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth"
)

type (
	// CustomerIdentityOrderService loads orders for an authenticated user
	CustomerIdentityOrderService interface {
		// Get all orders for a customer
		Get(ctx context.Context, identity auth.Identity) ([]*Order, error)
		// GetByID returns a single order for a customer
		GetByID(ctx context.Context, identity auth.Identity, orderID string) (*Order, error)
	}
)
