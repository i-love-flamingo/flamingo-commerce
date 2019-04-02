package infrastructure

import (
	"context"

	customerDomain "flamingo.me/flamingo-commerce/v3/customer/domain"

	"flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	//NilCustomerServiceAdapter - a Adpater for CustomerSerive that returns always NotFound
	NilCustomerServiceAdapter struct{}
)

//GetByAuth - implementation of required interface to get a customer based on Auth infos
func (n *NilCustomerServiceAdapter) GetByAuth(ctx context.Context, auth domain.Auth) (customerDomain.Customer, error) {
	return nil, customerDomain.ErrCustomerNotFoundError
}
