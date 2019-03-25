package infrastructure

import (
	"context"

	customerDomain "flamingo.me/flamingo-commerce/v3/customer/domain"

	"flamingo.me/flamingo/v3/core/auth/domain"
)

type (
	NilCustomerServiceAdapter struct{}
)

func (n *NilCustomerServiceAdapter) GetByAuth(ctx context.Context, auth domain.Auth) (customerDomain.Customer, error) {
	return nil, customerDomain.ErrCustomerNotFoundError
}
