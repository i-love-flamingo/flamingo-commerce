package graphql

import (
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"

)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutMutationResolver struct {
	graphQLService *Service
}

// Inject dependencies
func (r *CommerceCheckoutMutationResolver) Inject(
	cartService *cartApplication.CartService,
	graphQLService *Service) {
	r.graphQLService = graphQLService
}
