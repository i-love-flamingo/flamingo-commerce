package controller

import (
	"context"
	"net/http"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	placeorderDomain "flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
)

type (
	// APIController for checkout rest api
	APIController struct {
		responder            *web.Responder
		placeorderHandler    *placeorder.Handler
		cartService          *cartApplication.CartService
		logger               flamingo.Logger
		decoratedCartFactory *decorator.DecoratedCartFactory
	}

	// startPlaceOrderResult result of start place order
	startPlaceOrderResult struct {
		UUID string
	}

	// placeOrderContext infos
	placeOrderContext struct {
		Cart                 *cart.Cart
		OrderInfos           *placedOrderInfos
		State                string
		StateData            process.StateData
		UUID                 string
		FailedReason         string
		CartValidationResult *validation.Result
	}

	// placedOrderInfos infos
	placedOrderInfos struct {
		PaymentInfos        []application.PlaceOrderPaymentInfo
		PlacedOrderInfos    []placeorderDomain.PlacedOrderInfo
		Email               string
		PlacedDecoratedCart *decorator.DecoratedCart
	}

	// errorResponse format
	errorResponse struct {
		Code    string
		Message string
	} // @name checkoutError
)

// Inject dependencies
func (c *APIController) Inject(
	responder *web.Responder,
	placeorderHandler *placeorder.Handler,
	cartService *cartApplication.CartService,
	decoratedCartFactory *decorator.DecoratedCartFactory,
	logger flamingo.Logger,
) {
	c.responder = responder
	c.placeorderHandler = placeorderHandler
	c.decoratedCartFactory = decoratedCartFactory
	c.cartService = cartService
	c.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "apicontroller")
}

// StartPlaceOrderAction starts a new process
// @Summary Starts the place order process, which is a background process handling payment and rollbacks if required.
// @Tags Checkout
// @Produce json
// @Success 201 {object} startPlaceOrderResult "201 if new process was started"
// @Failure 500 {object} errorResponse
// @Failure 400 {object} errorResponse
// @Param returnURL query string true "the returnURL that should be used after an external payment flow"
// @Router /api/v1/checkout/placeorder [put]
func (c *APIController) StartPlaceOrderAction(ctx context.Context, r *web.Request) web.Result {
	session := web.SessionFromContext(ctx)
	cart, err := c.cartService.GetCartReceiverService().ViewCart(ctx, session)
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}
	returnURLRaw, err := r.Query1("returnURL")
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "400", Message: "returnURL missing"})
		response.Status(http.StatusBadRequest)
		return response
	}
	var returnURL *url.URL
	returnURL, err = url.Parse(returnURLRaw)
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "400", Message: err.Error()})
		response.Status(http.StatusBadRequest)
		return response
	}

	startPlaceOrderCommand := placeorder.StartPlaceOrderCommand{Cart: *cart, ReturnURL: returnURL}
	pctx, err := c.placeorderHandler.StartPlaceOrder(ctx, startPlaceOrderCommand)
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}
	response := c.responder.Data(startPlaceOrderResult{
		UUID: pctx.UUID,
	})
	response.Status(http.StatusCreated)
	return response
}

// CancelPlaceOrderAction cancels a running place order process
// @Summary Cancels a running place order process
// @Tags Checkout
// @Produce json
// @Success 200 {boolean} boolean
// @Failure 500 {object} errorResponse
// @Router /api/v1/checkout/placeorder/cancel [post]
func (c *APIController) CancelPlaceOrderAction(ctx context.Context, r *web.Request) web.Result {
	err := c.placeorderHandler.CancelPlaceOrder(ctx, placeorder.CancelPlaceOrderCommand{})
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}
	return c.responder.Data(true)
}

// ClearPlaceOrderAction clears the last place order if in final state
// @Summary Clears the last placed order if in final state
// @Tags Checkout
// @Produce json
// @Success 200 {boolean} boolean
// @Failure 500 {object} errorResponse
// @Router /api/v1/checkout/placeorder [delete]
func (c *APIController) ClearPlaceOrderAction(ctx context.Context, r *web.Request) web.Result {
	err := c.placeorderHandler.ClearPlaceOrder(ctx)
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}
	return c.responder.Data(true)
}

// CurrentPlaceOrderContextAction returns the last saved context
// @Summary Returns the last saved context
// @Tags Checkout
// @Produce json
// @Success 200 {object} placeOrderContext
// @Failure 500 {object} errorResponse
// @Router /api/v1/checkout/placeorder [get]
func (c *APIController) CurrentPlaceOrderContextAction(ctx context.Context, r *web.Request) web.Result {
	pctx, err := c.placeorderHandler.CurrentContext(ctx)
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}

	return c.responder.Data(c.getPlaceOrderContext(ctx, pctx))
}

func (c *APIController) getPlaceOrderContext(ctx context.Context, pctx *process.Context) placeOrderContext {
	var orderInfos *placedOrderInfos
	if pctx.PlaceOrderInfo != nil {
		decoratedCart := c.decoratedCartFactory.Create(ctx, pctx.Cart)
		orderInfos = &placedOrderInfos{
			PaymentInfos:        pctx.PlaceOrderInfo.PaymentInfos,
			PlacedOrderInfos:    pctx.PlaceOrderInfo.PlacedOrders,
			Email:               pctx.PlaceOrderInfo.ContactEmail,
			PlacedDecoratedCart: decoratedCart,
		}
	}

	var validationResult *validation.Result
	var failedReason string
	if pctx.FailedReason != nil {
		failedReason = pctx.FailedReason.Reason()
		if reason, ok := pctx.FailedReason.(process.CartValidationErrorReason); ok {
			validationResult = &reason.ValidationResult
		}
	}

	return placeOrderContext{
		Cart:                 &pctx.Cart,
		OrderInfos:           orderInfos,
		State:                pctx.CurrentStateName,
		StateData:            pctx.CurrentStateData,
		FailedReason:         failedReason,
		CartValidationResult: validationResult,
		UUID:                 pctx.UUID,
	}
}

// RefreshPlaceOrderAction returns the current place order context and proceeds the process in a non blocking way
// @Summary Returns the current place order context and proceeds the process in a non blocking way
// @Tags Checkout
// @Produce json
// @Success 200 {object} placeOrderContext
// @Failure 500 {object} errorResponse
// @Router /api/v1/checkout/placeorder/refresh [post]
func (c *APIController) RefreshPlaceOrderAction(ctx context.Context, r *web.Request) web.Result {
	pctx, err := c.placeorderHandler.RefreshPlaceOrder(ctx, placeorder.RefreshPlaceOrderCommand{})
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}
	return c.responder.Data(c.getPlaceOrderContext(ctx, pctx))
}

// RefreshPlaceOrderBlockingAction proceeds the process and returns the place order context afterwards (blocking)
// @Summary Proceeds the process and returns the place order context afterwards (blocking)
// @Description This is useful to get the most recent place order context, for example after returning from an external payment provider
// @Tags Checkout
// @Produce json
// @Success 200 {object} placeOrderContext
// @Failure 500 {object} errorResponse
// @Router /api/v1/checkout/placeorder/refresh-blocking [post]
func (c *APIController) RefreshPlaceOrderBlockingAction(ctx context.Context, r *web.Request) web.Result {
	pctx, err := c.placeorderHandler.RefreshPlaceOrderBlocking(ctx, placeorder.RefreshPlaceOrderCommand{})
	if err != nil {
		response := c.responder.Data(errorResponse{Code: "500", Message: err.Error()})
		response.Status(http.StatusInternalServerError)
		return response
	}
	return c.responder.Data(c.getPlaceOrderContext(ctx, pctx))
}
