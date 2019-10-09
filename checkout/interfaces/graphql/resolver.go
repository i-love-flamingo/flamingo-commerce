package graphql

import (
	"context"
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo/v3/framework/web"
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

func (r *CommerceCheckoutMutationResolver) CommerceCheckoutPlaceOrder(ctx context.Context) (*application.PlaceOrderInfo, error) {
	req := web.RequestFromContext(ctx)
	session := req.Session()
	placeOrderInfo, err := application.OrderService{}.PlaceOrder(ctx, session)

	if err != nil {
		return placeOrderInfo, nil
	}

	return nil, err
}
