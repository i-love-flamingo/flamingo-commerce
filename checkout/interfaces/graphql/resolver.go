package graphql

import (
	"context"
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/flamingo/v3/framework/web"
	"net/url"
)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutMutationResolver struct {
	orderService         *application.OrderService
	decoratedCartFactory *decorator.DecoratedCartFactory
	cartService          *cartApplication.CartService
}

// Inject dependencies
func (r *CommerceCheckoutMutationResolver) Inject(
	orderService *application.OrderService,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	cartService *cartApplication.CartService) {
	r.orderService = orderService
	r.decoratedCartFactory = decoratedCartFactory
	r.cartService = cartService

}

//CommerceCheckoutPlaceOrder places the order.
// TODO - handleEarlyplaceorder and allow multiple calls to this resolver
// TODO - eventually extract common logic from controller to application service
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutPlaceOrder(ctx context.Context, returnURLRaw string) (*dto.PlaceOrderResult, error) {
	returnURL, err := url.Parse(returnURLRaw)
	if err != nil {
		return nil, err
	}

	decoratedCart, err := r.cartService.GetCartReceiverService().ViewDecoratedCart(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return nil, err
	}
	validationResult := r.cartService.ValidateCart(ctx, web.SessionFromContext(ctx), decoratedCart)
	if !validationResult.IsValid() {
		return &dto.PlaceOrderResult{
			Status:               dto.INVALID,
			CartValidationResult: &validationResult,
		}, nil
	}

	gateway, err := r.orderService.GetPaymentGateway(ctx, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		return &dto.PlaceOrderResult{
			Status: dto.ERROR,
			Error:  &dto.Error{ErrorKey: err.Error(), IsPaymentError: true},
		}, nil
	}

	// start the payment flow
	_, err = gateway.StartFlow(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID, returnURL)
	if err != nil {
		return &dto.PlaceOrderResult{
			Status: dto.ERROR,
			Error:  &dto.Error{ErrorKey: err.Error(), IsPaymentError: true},
		}, nil
	}

	//TODO flowResult.Earlyplace order check!

	info, err := r.orderService.CurrentCartPlaceOrderWithPaymentProcessing(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return &dto.PlaceOrderResult{
		Status: dto.ORDERSUCCESS,
		OrderSuccessData: &dto.OrderSuccessData{
			PaymentInfos:        info.PaymentInfos,
			PlacedOrderInfos:    info.PlacedOrders,
			Email:               info.ContactEmail,
			PlacedDecoratedCart: *r.decoratedCartFactory.Create(ctx, info.Cart),
		},
	}, nil
}
