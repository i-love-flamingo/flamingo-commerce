package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
}

// Inject dependencies
func (r *CommerceCartQueryResolver) Inject(applicationCartReceiverService *application.CartReceiverService) {
	r.applicationCartReceiverService = applicationCartReceiverService
}

// CommerceCart getter for queries
func (r *CommerceCartQueryResolver) CommerceCart(ctx context.Context) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	return r.applicationCartReceiverService.ViewDecoratedCart(ctx, req.Session())
}

// CommerceCartMutationResolver resolves cart mutations
type CommerceCartMutationResolver struct {
	q                      *CommerceCartQueryResolver
	applicationCartService *application.CartService
}

// Inject dependencies
func (r *CommerceCartMutationResolver) Inject(q *CommerceCartQueryResolver, applicationCartService *application.CartService) *CommerceCartMutationResolver {
	r.q = q
	r.applicationCartService = applicationCartService
	return r
}

// CommerceAddToCart mutation for adding products to the current users cart
func (r *CommerceCartMutationResolver) CommerceAddToCart(ctx context.Context, marketplaceCode string, qty *int, deliveryCode string) (*decorator.DecoratedCart, error) {
	if qty == nil {
		one := 1
		qty = &one
	}

	req := web.RequestFromContext(ctx)

	addRequest := r.applicationCartService.BuildAddRequest(ctx, marketplaceCode, "", *qty)

	_, err := r.applicationCartService.AddProduct(ctx, req.Session(), deliveryCode, addRequest)
	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}
