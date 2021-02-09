package controller

import (
	"context"
	"fmt"

	// cart type is referenced in swag comment and requires empty import
	_ "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	// domain type is referenced in swag comment and requires empty import
	_ "flamingo.me/flamingo-commerce/v3/payment/domain"

	"flamingo.me/flamingo-commerce/v3/checkout/application"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// PaymentAPIController for payment api
	PaymentAPIController struct {
		responder    *web.Responder
		orderService *application.OrderService
		logger       flamingo.Logger
	}

	resultError struct {
		Message string
		Code    string
	} // @name paymentResultError
)

// Inject dependencies
func (pc *PaymentAPIController) Inject(
	responder *web.Responder,
	Logger flamingo.Logger,
	orderService *application.OrderService,
) {
	pc.responder = responder
	pc.logger = Logger.WithField("category", "PaymentApiController")
	pc.orderService = orderService
}

// Status Get Payment Status
// @Summary Get the payment status of current cart (or last placed cart)
// @Tags Payment
// @Produce json
// @Success 200 {object} domain.FlowStatus{data=cart.Cart}
// @Failure 500 {object} resultError
// @Router /api/v1/payment/status [get]
func (pc *PaymentAPIController) Status(ctx context.Context, r *web.Request) web.Result {
	decoratedCart, err := pc.orderService.LastPlacedOrCurrentCart(ctx)

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
