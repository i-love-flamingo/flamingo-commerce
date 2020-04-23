package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/customer/interfaces/dtocustomer"
)

type(
	// CustomerResolver graphql resolver
	CustomerResolver struct{

	}
)

// CommerceCustomerStatus ...
func (r *CustomerResolver) CommerceCustomerStatus (ctx context.Context) (*dtocustomer.CustomerStatusResult, error) {
	return &dtocustomer.CustomerStatusResult{}, nil
}

// CommerceCustomer ...
func (r *CustomerResolver) CommerceCustomer (ctx context.Context) (*dtocustomer.CustomerResult, error) {
	return &dtocustomer.CustomerResult{}, nil
}
