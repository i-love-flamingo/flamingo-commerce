// +build !graphql

package graphql

import (
	graphql3 "flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql"
	graphql2 "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql"
	graphql4 "flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql"
	graphql6 "flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql"
	graphql5 "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"

	"flamingo.me/graphql"
)

// this file is a starting point for the project specific resolvers
// it will not be regenerated!

type rootResolver struct {
	queryResolver    *queryResolver
	mutationResolver *mutationResolver
}

// interface guard
var _ ResolverRoot = new(rootResolver)

// Inject root resolver
func (r *rootResolver) Inject(queryResolver *queryResolver, mutationResolver *mutationResolver) {
	r.queryResolver = queryResolver
	r.mutationResolver = mutationResolver
}

// Query getter
func (r *rootResolver) Query() QueryResolver {
	return r.queryResolver
}

// Mutation getter
func (r *rootResolver) Mutation() MutationResolver {
	return r.mutationResolver
}

type queryResolver struct {
	*graphql.FlamingoQueryResolver
	*graphql2.CommerceCategoryQueryResolver `inject:""`
	*graphql3.CommerceCartQueryResolver     `inject:""`
	*graphql4.CommerceCheckoutQueryResolver `inject:""`
	*graphql5.CommerceProductQueryResolver  `inject:""`
	*graphql6.CustomerResolver              `inject:""`
}

// Inject dependencies
func (r *queryResolver) Inject(flamingoQueryResolver *graphql.FlamingoQueryResolver) {
	r.FlamingoQueryResolver = flamingoQueryResolver
}

type mutationResolver struct {
	*graphql.FlamingoQueryResolver
	*graphql3.CommerceCartMutationResolver     `inject:""`
	*graphql4.CommerceCheckoutMutationResolver `inject:""`
}

// Inject dependencies
func (r *mutationResolver) Inject(flamingoQueryResolver *graphql.FlamingoQueryResolver) {
	r.FlamingoQueryResolver = flamingoQueryResolver
}
