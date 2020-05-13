package infrastructure

import (
	"context"

	customerDomain "flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo/v3/core/auth"

	"flamingo.me/flamingo/v3/core/oauth/domain"
)

type (
	// NilCustomerServiceAdapter for CustomerService and CustomerIdentityService that returns always NotFound
	NilCustomerServiceAdapter struct{}
)

var _ customerDomain.CustomerService = new(NilCustomerServiceAdapter)
var _ customerDomain.CustomerIdentityService = new(NilCustomerServiceAdapter)

// GetByAuth - implementation of required interface to get a customer based on Auth infos
func (n *NilCustomerServiceAdapter) GetByAuth(context.Context, domain.Auth) (customerDomain.Customer, error) {
	return nil, customerDomain.ErrCustomerNotFoundError
}

// GetByIdentity retrieves the authenticated customer by Identity
func (n *NilCustomerServiceAdapter) GetByIdentity(context.Context, auth.Identity) (customerDomain.Customer, error) {
	return nil, customerDomain.ErrCustomerNotFoundError
}
