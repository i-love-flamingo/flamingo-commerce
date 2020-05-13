package graphql

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
)

// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
	applicationCartService         *application.CartService
}

// Inject dependencies
func (r *CommerceCartQueryResolver) Inject(
	applicationCartReceiverService *application.CartReceiverService,
	cartService *application.CartService,
) {
	r.applicationCartReceiverService = applicationCartReceiverService
	r.applicationCartService = cartService
}

// CommerceCart getter for queries
func (r *CommerceCartQueryResolver) CommerceCart(ctx context.Context) (*dto.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)
	dc, err := r.applicationCartReceiverService.ViewDecoratedCart(ctx, req.Session())
	if err != nil {
		return nil, err
	}

	return dto.NewDecoratedCart(dc), nil
}

// CommerceCartValidator to trigger the cart validation service
func (r *CommerceCartQueryResolver) CommerceCartValidator(ctx context.Context) (*validation.Result, error) {
	session := web.SessionFromContext(ctx)

	decoratedCart, err := r.applicationCartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		return nil, err
	}

	result := r.applicationCartService.ValidateCart(ctx, session, decoratedCart)

	return &result, nil
}
