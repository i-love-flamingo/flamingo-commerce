package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
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

func (r *CommerceCartMutationResolver) CommerceDeleteItem(ctx context.Context, itemId string, deliveryCode string) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	err := r.applicationCartService.DeleteItem(ctx, req.Session(), itemId, deliveryCode)

	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}

// CommerceDeleteCartDelivery mutation for removing deliveries from current users cart
func (r *CommerceCartMutationResolver) CommerceDeleteCartDelivery(ctx context.Context, deliveryCode string) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	_, err := r.applicationCartService.DeleteDelivery(ctx, req.Session(), deliveryCode)
	if err != nil {
		return nil, err
	}
	return r.q.CommerceCart(ctx)
}

// CommerceUpdateItemQty mutation for updating item quantity
func (r *CommerceCartMutationResolver) CommerceUpdateItemQty(ctx context.Context, itemID string, deliveryCode string, qty int) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	err := r.applicationCartService.UpdateItemQty(ctx, req.Session(), itemID, deliveryCode, qty)
	if err != nil {
		return nil, err
	}
	return r.q.CommerceCart(ctx)
}

type CommerceCartItemResolver struct {
	productResolver *graphql.CommerceProductQueryResolver
}

// Inject dependencies
func (c *CommerceCartItemResolver) Inject(productResolver *graphql.CommerceProductQueryResolver) {
	c.productResolver = productResolver
}

// Product returns the product for the corresponding cart item
func (c *CommerceCartItemResolver) Product(ctx context.Context, item *cart.Item) (domain.BasicProduct, error) {
	marketplaceCode := item.MarketplaceCode

	return c.productResolver.CommerceProduct(ctx, marketplaceCode)
}
