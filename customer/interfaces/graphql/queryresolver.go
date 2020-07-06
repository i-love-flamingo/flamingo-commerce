package graphql

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/customer/application"
	"flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql/dtocustomer"
)

type (
	// CustomerResolver graphql resolver
	CustomerResolver struct {
		service *application.Service
	}
)

// Inject dependencies
func (r *CustomerResolver) Inject(
	service *application.Service,
) *CustomerResolver {
	r.service = service

	return r
}

// CommerceCustomerStatus resolves the commerce customer query
// Deprecated: use commerce customer query resolver instead
func (r *CustomerResolver) CommerceCustomerStatus(ctx context.Context) (*dtocustomer.CustomerStatusResult, error) {
	userID, err := r.service.GetUserID(ctx, web.RequestFromContext(ctx))
	if errors.Is(err, application.ErrNoIdentity) {
		return &dtocustomer.CustomerStatusResult{IsLoggedIn: false}, nil
	}

	if err != nil {
		return nil, err
	}

	return &dtocustomer.CustomerStatusResult{
		IsLoggedIn: true,
		UserID:     userID,
	}, nil
}

// CommerceCustomer resolver the commerce customer
func (r *CustomerResolver) CommerceCustomer(ctx context.Context) (*dtocustomer.CustomerResult, error) {
	user, err := r.service.GetForIdentity(ctx, web.RequestFromContext(ctx))
	if errors.Is(err, application.ErrNoIdentity) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	result := &dtocustomer.CustomerResult{
		ID:           user.GetID(),
		PersonalData: user.GetPersonalData(),
		Addresses:    user.GetAddresses(),
	}

	if address := user.GetDefaultShippingAddress(); address != nil {
		result.DefaultShippingAddress = *address
	}

	if address := user.GetDefaultBillingAddress(); address != nil {
		result.DefaultBillingAddress = *address
	}

	return result, nil
}
