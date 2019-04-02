package domain

import (
	"context"

	"flamingo.me/flamingo/v3/core/oauth/domain"
)

type (

	// CustomerOrderService interface
	CustomerOrderService interface {
		Get(context.Context, domain.Auth) ([]*Order, error)
		GetByID(context.Context, domain.Auth, string) (*Order, error)
	}
)
