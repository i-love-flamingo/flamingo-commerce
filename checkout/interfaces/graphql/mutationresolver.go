package graphql

import (
	"context"
	"net/url"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	graphqlDto "flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutMutationResolver struct {
	placeorderHandler    *placeorder.Handler
	cartService          *cartApplication.CartService
	stateMapping         map[string]dto.State
	logger               flamingo.Logger
	decoratedCartFactory *decorator.DecoratedCartFactory
}

// Inject dependencies
func (r *CommerceCheckoutMutationResolver) Inject(
	placeorderHandler *placeorder.Handler,
	cartService *cartApplication.CartService,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	stateMapping map[string]dto.State,
	logger flamingo.Logger) {
	r.placeorderHandler = placeorderHandler
	r.decoratedCartFactory = decoratedCartFactory
	r.cartService = cartService
	r.stateMapping = stateMapping
	r.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "graphql")

}

// CommerceCheckoutRefreshPlaceOrder refreshes the current place order and proceeds the process
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutRefreshPlaceOrder(ctx context.Context) (*dto.PlaceOrderContext, error) {
	return r.refresh(ctx, r.placeorderHandler.RefreshPlaceOrder)
}

func (r *CommerceCheckoutMutationResolver) refresh(
	ctx context.Context,
	refreshFnc func(context.Context, placeorder.RefreshPlaceOrderCommand) (*process.Context, error),
) (*dto.PlaceOrderContext, error) {

	poctx, err := refreshFnc(ctx, placeorder.RefreshPlaceOrderCommand{})
	if err != nil {
		return nil, err
	}

	dc := graphqlDto.NewDecoratedCart(r.decoratedCartFactory.Create(ctx, poctx.Cart))

	return &dto.PlaceOrderContext{
		Cart:       dc,
		OrderInfos: nil,
		State:      r.mapStateToGraphQL(*poctx),
		UUID:       poctx.UUID,
	}, nil
}

// CommerceCheckoutRefreshPlaceOrderBlocking refreshes the current place order blocking
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutRefreshPlaceOrderBlocking(ctx context.Context) (*dto.PlaceOrderContext, error) {
	return r.refresh(ctx, r.placeorderHandler.RefreshPlaceOrderBlocking)
}

// CommerceCheckoutStartPlaceOrder starts a new process (if not running)
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutStartPlaceOrder(ctx context.Context, returnURLRaw string) (*dto.StartPlaceOrderResult, error) {
	session := web.SessionFromContext(ctx)
	cart, err := r.cartService.GetCartReceiverService().ViewCart(ctx, session)
	if err != nil {
		return nil, err
	}
	var returnURL *url.URL
	if returnURLRaw != "" {
		returnURL, err = url.Parse(returnURLRaw)
		if err != nil {
			return nil, err
		}
	}
	startPlaceOrderCommand := placeorder.StartPlaceOrderCommand{Cart: *cart, ReturnURL: returnURL}
	pctx, err := r.placeorderHandler.StartPlaceOrder(ctx, startPlaceOrderCommand)
	if err == placeorder.ErrAnotherPlaceOrderProcessRunning {
		dtopctx, err := r.CommerceCheckoutRefreshPlaceOrder(ctx)
		if err != nil {
			return nil, err
		}
		return &dto.StartPlaceOrderResult{
			UUID: dtopctx.UUID,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &dto.StartPlaceOrderResult{
		UUID: pctx.UUID,
	}, nil
}

// CommerceCheckoutCancelPlaceOrder cancels a running place order
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutCancelPlaceOrder(ctx context.Context) (bool, error) {
	err := r.placeorderHandler.CancelPlaceOrder(ctx, placeorder.CancelPlaceOrderCommand{})

	return err == nil, err
}

func (r *CommerceCheckoutMutationResolver) mapStateToGraphQL(pctx process.Context) dto.State {
	resultState := r.stateMapping[pctx.State]
	resultState.MapFrom(pctx)

	return resultState
}
