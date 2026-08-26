package graphql

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	graphqlDto "flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
)

// CommerceCheckoutQueryResolver resolves graphql checkout queries
type CommerceCheckoutQueryResolver struct {
	placeOrderHandler    *placeorder.Handler
	decoratedCartFactory *decorator.DecoratedCartFactory
	stateMapper          *dto.StateMapper
	logger               flamingo.Logger
}

// Inject dependencies
func (r *CommerceCheckoutQueryResolver) Inject(
	placeOrderHandler *placeorder.Handler,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	stateMapper *dto.StateMapper,
	logger flamingo.Logger,
) {
	r.placeOrderHandler = placeOrderHandler
	r.decoratedCartFactory = decoratedCartFactory
	r.stateMapper = stateMapper
	r.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "graphql")
}

// CommerceCheckoutActivePlaceOrder checks if there is an order in unfinished state
func (r *CommerceCheckoutQueryResolver) CommerceCheckoutActivePlaceOrder(ctx context.Context) (bool, error) {
	active, err := r.placeOrderHandler.HasUnfinishedProcess(ctx)
	if err != nil {
		r.logger.Error("Failed to check unfinished place order process", err)
		return false, interfaces.ErrCheckoutGeneral
	}

	return active, nil
}

// CommerceCheckoutCurrentContext returns the last saved context
func (r *CommerceCheckoutQueryResolver) CommerceCheckoutCurrentContext(ctx context.Context) (*dto.PlaceOrderContext, error) {
	pctx, err := r.placeOrderHandler.CurrentContext(ctx)
	if err != nil {
		// keep existing external error contract for expected process states
		if errors.Is(err, placeorder.ErrNoPlaceOrderProcess) {
			return nil, interfaces.ErrNoPlaceOrderProcess
		}

		r.logger.Error("Failed to get current place order context", err)

		return nil, interfaces.ErrCheckoutGeneral
	}

	dc := graphqlDto.NewDecoratedCart(r.decoratedCartFactory.Create(ctx, pctx.Cart))

	graphQLState, err := r.stateMapper.Map(*pctx)
	if err != nil {
		r.logger.Error("Failed to map place order context state", err)
		return nil, interfaces.ErrCheckoutGeneral
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
