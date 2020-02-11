package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	graphqlDto "flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
)

// CommerceCheckoutQueryResolver resolves graphql checkout queries
type CommerceCheckoutQueryResolver struct {
	placeOrderHandler    *placeorder.Handler
	decoratedCartFactory *decorator.DecoratedCartFactory
	stateMapper          *dto.StateMapper
}

// Inject dependencies
func (r *CommerceCheckoutQueryResolver) Inject(
	placeOrderHandler *placeorder.Handler,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	stateMapper *dto.StateMapper,
) {
	r.placeOrderHandler = placeOrderHandler
	r.decoratedCartFactory = decoratedCartFactory
	r.stateMapper = stateMapper
}

// CommerceCheckoutActivePlaceOrder checks if there is an order in unfinished state
func (r *CommerceCheckoutQueryResolver) CommerceCheckoutActivePlaceOrder(ctx context.Context) (bool, error) {
	return r.placeOrderHandler.HasUnfinishedProcess(ctx)
}

// CommerceCheckoutCurrentContext returns the last saved context
func (r *CommerceCheckoutQueryResolver) CommerceCheckoutCurrentContext(ctx context.Context) (*dto.PlaceOrderContext, error) {
	pctx, err := r.placeOrderHandler.CurrentContext(ctx)
	if err != nil {
		return nil, err
	}

	dc := graphqlDto.NewDecoratedCart(r.decoratedCartFactory.Create(ctx, pctx.Cart))

	graphQLState, err := r.stateMapper.Map(*pctx)
	if err != nil {
		return nil, err
	}

	var orderInfos *dto.PlacedOrderInfos
	if pctx.PlaceOrderInfo != nil {
		orderInfos = &dto.PlacedOrderInfos{
			PaymentInfos:        pctx.PlaceOrderInfo.PaymentInfos,
			PlacedOrderInfos:    pctx.PlaceOrderInfo.PlacedOrders,
			Email:               pctx.PlaceOrderInfo.ContactEmail,
			PlacedDecoratedCart: dc,
		}
	}

	return &dto.PlaceOrderContext{
		Cart:       dc,
		OrderInfos: orderInfos,
		State:      graphQLState,
		UUID:       pctx.UUID,
	}, nil
}
