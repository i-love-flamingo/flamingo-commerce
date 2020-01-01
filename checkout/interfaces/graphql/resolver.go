package graphql

import (
	"context"
	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	dto2 "flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// CommerceCheckoutMutationResolver resolves graphql checkout mutations
type CommerceCheckoutQueryResolver struct {
	orderService         *application.OrderService
	decoratedCartFactory *decorator.DecoratedCartFactory
	cartService          *cartApplication.CartService
	logger               flamingo.Logger
}

// Inject dependencies
func (r *CommerceCheckoutQueryResolver) Inject(
	orderService *application.OrderService,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	cartService *cartApplication.CartService,
	logger flamingo.Logger) {
	r.orderService = orderService
	r.decoratedCartFactory = decoratedCartFactory
	r.cartService = cartService
	r.logger = logger.WithField(flamingo.LogKeyModule, "om3oms").WithField(flamingo.LogKeyCategory, "graphql")

}

func (r *CommerceCheckoutQueryResolver) CommerceCheckoutPlaceOrderContext(ctx context.Context) (*dto.PlaceOrderContext, error) {
	return &dto.PlaceOrderContext{
		Cart:       nil,
		OrderInfos: nil,
		State:      nil,
	}, nil
}

/*
//CommerceCheckoutPlaceOrder places the order.
// TODO - handleEarlyplaceorder and allow multiple calls to this resolver
// TODO - eventually extract common logic from controller to application service
func (r *CommerceCheckoutMutationResolver) CommerceCheckoutPlaceOrder(ctx context.Context, returnURLRaw string) (*dto.PlaceOrderResult, error) {
	returnURL, err := url.Parse(returnURLRaw)
	if err != nil {
		return nil, err
	}

	// reserve an unique order id for later order placing
	_, err = r.cartService.ReserveOrderIDAndSave(ctx, web.SessionFromContext(ctx))
	if err != nil {
		r.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error ", err)
		return &dto.PlaceOrderResult{
			Status: dto.INVALID,
			Error:  &dto.Error{ErrorKey: errors.New("reserve-order-id-failed").Error(), IsPaymentError: false},
		}, nil
	}

	decoratedCart, err := r.cartService.GetCartReceiverService().ViewDecoratedCart(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return nil, err
	}

	if decoratedCart.Cart.IsEmpty() {
		return &dto.PlaceOrderResult{
			Status: dto.ERROR,
			Error:  &dto.Error{ErrorKey: errors.New("place-order_cart-is-empty").Error(), IsPaymentError: false},
		}, nil
	}

	if !decoratedCart.Cart.IsPaymentSelected() {
		return &dto.PlaceOrderResult{
			Status: dto.ERROR,
			Error:  &dto.Error{ErrorKey: errors.New("place-order_payment-selection-not-set").Error(), IsPaymentError: false},
		}, nil
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
			PlacedDecoratedCart: dto2.NewDecoratedCart(r.decoratedCartFactory.Create(ctx, info.Cart)),
		},
	}, nil
}

*/
