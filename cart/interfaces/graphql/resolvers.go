package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
}

// Inject dependencies
func (r *CommerceCartQueryResolver) Inject(
	applicationCartReceiverService *application.CartReceiverService,
) {
	r.applicationCartReceiverService = applicationCartReceiverService
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
