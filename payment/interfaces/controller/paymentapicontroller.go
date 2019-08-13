package controller

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/application"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// PaymentAPIController for payment api
	PaymentAPIController struct {
		responder                      *web.Responder
		cartReceiverService            *cartApplication.CartReceiverService
		orderService                   *application.OrderService
		logger                         flamingo.Logger
		decoratedCartFactory           *decorator.DecoratedCartFactory
		applicationCartReceiverService *cartApplication.CartReceiverService
	}

	resultError struct {
		Message string
		Code    string
	}
)

// Inject dependencies
func (pc *PaymentAPIController) Inject(
	responder *web.Responder,
	Logger flamingo.Logger,
	cartReceiver *cartApplication.CartReceiverService,
	orderService *application.OrderService,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	applicationCartReceiverService *cartApplication.CartReceiverService,
) {
	pc.responder = responder
	pc.logger = Logger.WithField("category", "PaymentApiController")
	pc.cartReceiverService = cartReceiver
	pc.orderService = orderService
	pc.decoratedCartFactory = decoratedCartFactory
	pc.applicationCartReceiverService = applicationCartReceiverService
}

// Status Get Payment Status
func (pc *PaymentAPIController) Status(ctx context.Context, r *web.Request) web.Result {
	decoratedCart, err := pc.lastPlacedOrCurrentCart(ctx)

	if err != nil {
		pc.logger.Warn("Error when getting last used cart", err)
		return pc.responder.Data(resultError{
			Message: fmt.Sprintf("Cart not found: %v", err),
			Code:    "status.polling.cart.error",
		})
	}

	if decoratedCart.Cart.PaymentSelection == nil {
		pc.logger.Warn("Error because payment selection is empty")
		return pc.responder.Data(resultError{
			Message: "Payment selection is empty",
			Code:    "status.polling.paymentselection.error",
		})
	}

	gateway, err := pc.orderService.GetPaymentGateway(ctx, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		pc.logger.Warn("Error because payment gateway is not set", err)
		return pc.responder.Data(resultError{
			Message: fmt.Sprintf("Payment Gateway is not set: %v", err),
			Code:    "status.polling.paymentgateway.error",
		})
	}

	flowStatus, err := gateway.FlowStatus(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID)
	if err != nil {
		pc.logger.Warn("Error because flow status is unknown", err)
		return pc.responder.Data(resultError{
			Message: fmt.Sprintf("Flow status unknown: %v", err),
			Code:    "status.polling.flowstatus.error",
		})
	}

	return pc.responder.Data(flowStatus)
}

func (pc *PaymentAPIController) lastPlaceOrderInfo(ctx context.Context) (*application.PlaceOrderInfo, error) {
	lastPlacedOrder, err := pc.orderService.LastPlacedOrder(ctx)
	if err != nil {
		pc.logger.Warn("couldn't get last placed order from orderService:", err)
		return nil, err
	}

	return lastPlacedOrder, nil
}

// lastPlacedOrCurrentCart returns the decorated cart of the last placed order if there is one if not return the current cart
func (pc *PaymentAPIController) lastPlacedOrCurrentCart(ctx context.Context) (*decorator.DecoratedCart, error) {
	lastPlacedOrder, err := pc.lastPlaceOrderInfo(ctx)
	if err != nil {
		pc.logger.Warn("couldn't get last placed order from orderService:", err)
		return nil, err
	}

	if lastPlacedOrder != nil {
		// cart has been placed early use it
		return pc.decoratedCartFactory.Create(ctx, lastPlacedOrder.Cart), nil
	}

	// cart wasn't placed early, fetch it from service
	decoratedCart, err := pc.applicationCartReceiverService.ViewDecoratedCart(ctx, web.SessionFromContext(ctx))
	if err != nil {
		pc.logger.WithContext(ctx).Error("lastPlacedOrCurrentCart -> ViewDecoratedCart Error:", err)
		return nil, err
	}

	return decoratedCart, nil
}
