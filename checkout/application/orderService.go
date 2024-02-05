package application

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// OrderService defines the order service
	OrderService struct {
		logger                 flamingo.Logger
		cartService            *application.CartService
		cartReceiverService    *application.CartReceiverService
		deliveryInfoBuilder    cart.DeliveryInfoBuilder
		webCartPaymentGateways map[string]interfaces.WebCartPaymentGateway
		decoratedCartFactory   *decorator.DecoratedCartFactory
	}

	// PlaceOrderInfo struct defines the data of payments on placed orders
	PlaceOrderInfo struct {
		PaymentInfos []PlaceOrderPaymentInfo
		PlacedOrders placeorder.PlacedOrderInfos
		ContactEmail string
		Cart         cart.Cart
	}

	// PlaceOrderPaymentInfo holding payment infos
	PlaceOrderPaymentInfo struct {
		Gateway         string
		PaymentProvider string
		Method          string
		CreditCardInfo  *placeorder.CreditCardInfo
		Amount          priceDomain.Price
		Title           string
	}
)

const (
	// PaymentFlowStandardCorrelationID used as correlationid for the start of the payment (session scoped)
	PaymentFlowStandardCorrelationID = "checkout"

	// LastPlacedOrderSessionKey is the session key for storing the last placed order
	LastPlacedOrderSessionKey = "orderservice_last_placed"
)

var (
	// cartValidationFailCount counts validation failures on carts
	cartValidationFailCount = stats.Int64("flamingo-commerce/checkout/orders/cart_validation_failed", "Count of failures while validating carts", stats.UnitDimensionless)

	// noPaymentSelectionCount counts error for orders without payment selection
	noPaymentSelectionCount = stats.Int64("flamingo-commerce/checkout/orders/no_payment_selection", "Count of carts without having a selected payment", stats.UnitDimensionless)

	// paymentGatewayNotFoundCount counts errors if payment gateway for selected payment could not be found
	paymentGatewayNotFoundCount = stats.Int64("flamingo-commerce/checkout/orders/payment_gateway_not_found", "The selected payment gateway could not be found", stats.UnitDimensionless)

	// paymentFlowStatusErrorCount counts errors while fetching payment flow status
	paymentFlowStatusErrorCount = stats.Int64("flamingo-commerce/checkout/orders/payment_flow_status_error", "Count of failures while fetching payment flow status", stats.UnitDimensionless)

	// orderPaymentFromFlowErrorCount counts errors while fetching payment from flow
	orderPaymentFromFlowErrorCount = stats.Int64("flamingo-commerce/checkout/orders/order_payment_from_flow_error", "Count of failures while fetching payment from flow", stats.UnitDimensionless)

	// paymentFlowStatusFailedCanceledCount counts orders trying to be placed with payment status either failed or canceled
	paymentFlowStatusFailedCanceledCount = stats.Int64("flamingo-commerce/checkout/orders/payment_flow_status_failed_canceled", "Count of payments with status failed or canceled", stats.UnitDimensionless)

	// paymentFlowStatusAbortedCount counts orders trying to be placed with payment status aborted
	paymentFlowStatusAbortedCount = stats.Int64("flamingo-commerce/checkout/orders/payment_flow_status_aborted", "Count of payments with status aborted", stats.UnitDimensionless)

	// placeOrderFailCount counts failed placed orders
	placeOrderFailCount = stats.Int64("flamingo-commerce/checkout/orders/place_order_failed", "Count of failures while placing orders", stats.UnitDimensionless)

	// placeOrderSuccessCount counts successfully placed orders
	placeOrderSuccessCount = stats.Int64("flamingo-commerce/checkout/orders/place_order_successful", "Count of successfully placed orders", stats.UnitDimensionless)
)

func init() {
	gob.Register(PlaceOrderInfo{})
	openCensusViews := map[string]*stats.Int64Measure{
		"flamingo-commerce/checkout/orders/cart_validation_failed":              cartValidationFailCount,
		"flamingo-commerce/checkout/orders/no_payment_selection":                noPaymentSelectionCount,
		"flamingo-commerce/checkout/orders/payment_gateway_not_found":           paymentGatewayNotFoundCount,
		"flamingo-commerce/checkout/orders/payment_flow_status_error":           paymentFlowStatusErrorCount,
		"flamingo-commerce/checkout/orders/order_payment_from_flow_error":       orderPaymentFromFlowErrorCount,
		"flamingo-commerce/checkout/orders/payment_flow_status_failed_canceled": paymentFlowStatusFailedCanceledCount,
		"flamingo-commerce/checkout/orders/payment_flow_status_aborted":         paymentFlowStatusAbortedCount,
		"flamingo-commerce/checkout/orders/place_order_failed":                  placeOrderFailCount,
		"flamingo-commerce/checkout/orders/place_order_successful":              placeOrderSuccessCount,
	}

	for name, measure := range openCensusViews {
		err := opencensus.View(name, measure, view.Sum())
		if err != nil {
			panic(err)
		}

		stats.Record(context.Background(), measure.M(0))
	}
}

// Inject dependencies
func (os *OrderService) Inject(
	logger flamingo.Logger,
	CartService *application.CartService,
	CartReceiverService *application.CartReceiverService,
	DeliveryInfoBuilder cart.DeliveryInfoBuilder,
	webCartPaymentGatewayProvider interfaces.WebCartPaymentGatewayProvider,
	decoratedCartFactory *decorator.DecoratedCartFactory,

) {
	os.logger = logger.WithField(flamingo.LogKeyCategory, "checkout.OrderService").WithField(flamingo.LogKeyModule, "checkout")
	os.cartService = CartService
	os.cartReceiverService = CartReceiverService
	os.webCartPaymentGateways = webCartPaymentGatewayProvider()
	os.deliveryInfoBuilder = DeliveryInfoBuilder
	os.decoratedCartFactory = decoratedCartFactory
}

// GetPaymentGateway tries to get the supplied payment gateway by code from the registered payment gateways
func (os *OrderService) GetPaymentGateway(ctx context.Context, paymentGatewayCode string) (interfaces.WebCartPaymentGateway, error) {
	_, span := trace.StartSpan(ctx, "checkout/OrderService/GetPaymentGateway")
	defer span.End()

	gateway, ok := os.webCartPaymentGateways[paymentGatewayCode]
	if !ok {
		return nil, errors.New("Payment gateway " + paymentGatewayCode + " not found")
	}

	return gateway, nil
}

// GetAvailablePaymentGateways returns the list of registered WebCartPaymentGateway
func (os *OrderService) GetAvailablePaymentGateways(ctx context.Context) map[string]interfaces.WebCartPaymentGateway {
	_, span := trace.StartSpan(ctx, "checkout/OrderService/GetAvailablePaymentGateways")
	defer span.End()

	return os.webCartPaymentGateways
}

// CurrentCartPlaceOrder places the current cart without additional payment processing
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *web.Session, cartPayment placeorder.Payment) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/CurrentCartPlaceOrder")
	defer span.End()

	var info *PlaceOrderInfo
	var err error
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = func() (*PlaceOrderInfo, error) {
			decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(placeOrderContext, session)
			if err != nil {
				os.logger.WithContext(placeOrderContext).Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error ", err)

				return nil, err
			}

			return os.placeOrder(placeOrderContext, session, decoratedCart, cartPayment)
		}()
	})

	return info, err
}

func (os *OrderService) placeOrder(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, payment placeorder.Payment) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/placeOrder")
	defer span.End()

	validationResult := os.cartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		// record cartValidationFailCount metric
		stats.Record(ctx, cartValidationFailCount.M(1))
		os.logger.WithContext(ctx).Warn("Try to place an invalid cart")

		return nil, errors.New("cart is invalid")
	}

	placedOrderInfos, err := os.cartService.PlaceOrderWithCart(ctx, session, &decoratedCart.Cart, &payment)
	if err != nil {
		// record placeOrderFailCount metric
		stats.Record(ctx, placeOrderFailCount.M(1))
		os.logger.WithContext(ctx).Error("Error during place Order:" + err.Error())

		return nil, errors.New("error while placing the order. please contact customer support")
	}

	placeOrderInfo := os.preparePlaceOrderInfo(ctx, decoratedCart.Cart, placedOrderInfos, payment)
	os.storeLastPlacedOrder(ctx, placeOrderInfo)

	// record placeOrderSuccessCount metric
	stats.Record(ctx, placeOrderSuccessCount.M(1))

	return placeOrderInfo, nil
}

// CancelOrder cancels an previously placed order and returns the restored cart with the order content
func (os *OrderService) CancelOrder(ctx context.Context, session *web.Session, order *PlaceOrderInfo) (*cart.Cart, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/CancelOrder")
	defer span.End()

	return os.cartService.CancelOrder(ctx, session, order.PlacedOrders, order.Cart)
}

// CancelOrderWithoutRestore cancels an previously placed order
func (os *OrderService) CancelOrderWithoutRestore(ctx context.Context, session *web.Session, order *PlaceOrderInfo) error {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/CancelOrderWithoutRestore")
	defer span.End()

	return os.cartService.CancelOrderWithoutRestore(ctx, session, order.PlacedOrders)
}

// CurrentCartPlaceOrderWithPaymentProcessing places the current cart which is fetched from the context
func (os *OrderService) CurrentCartPlaceOrderWithPaymentProcessing(ctx context.Context, session *web.Session) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/CurrentCartPlaceOrderWithPaymentProcessing")
	defer span.End()

	var info *PlaceOrderInfo
	var err error
	// use a background context from here on to prevent the place order canceled by context cancel
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = func() (*PlaceOrderInfo, error) {
			// fetch decorated cart either via cache or freshly from cart receiver service
			decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(placeOrderContext, session)
			if err != nil {
				os.logger.WithContext(placeOrderContext).Warn("Cannot create decorated cart from cart")

				return nil, errors.New("cart is invalid")
			}

			return os.placeOrderWithPaymentProcessing(placeOrderContext, decoratedCart, session)
		}()
	})

	return info, err
}

// CartPlaceOrderWithPaymentProcessing places the cart passed to the function
// this function enables clients to pass a cart as is, without the usage of the cartReceiverService
func (os *OrderService) CartPlaceOrderWithPaymentProcessing(ctx context.Context, decoratedCart *decorator.DecoratedCart,
	session *web.Session) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/CartPlaceOrderWithPaymentProcessing")
	defer span.End()

	var info *PlaceOrderInfo
	var err error
	// use a background context from here on to prevent the place order canceled by context cancel
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = os.placeOrderWithPaymentProcessing(placeOrderContext, decoratedCart, session)
	})

	return info, err
}

// CartPlaceOrder places the cart passed to the function
// this function enables clients to pass a cart as is, without the usage of the cartReceiverService
func (os *OrderService) CartPlaceOrder(ctx context.Context, decoratedCart *decorator.DecoratedCart, payment placeorder.Payment) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/CartPlaceOrder")
	defer span.End()

	var info *PlaceOrderInfo
	var err error
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = os.placeOrder(placeOrderContext, web.SessionFromContext(ctx), decoratedCart, payment)
	})

	return info, err
}

// storeLastPlacedOrder stores the last placed order/cart in the session
func (os *OrderService) storeLastPlacedOrder(ctx context.Context, info *PlaceOrderInfo) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/storeLastPlacedOrder")
	defer span.End()

	session := web.SessionFromContext(ctx)

	_ = session.Store(LastPlacedOrderSessionKey, info)
}

// LastPlacedOrder returns the last placed order/cart if available
func (os *OrderService) LastPlacedOrder(ctx context.Context) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/LastPlacedOrder")
	defer span.End()

	session := web.SessionFromContext(ctx)

	lastPlacedOrder, found := session.Load(LastPlacedOrderSessionKey)
	if !found {
		return nil, nil
	}

	placeOrderInfo, ok := lastPlacedOrder.(PlaceOrderInfo)
	if !ok {
		return nil, errors.New("placeOrderInfo couldn't be received from session")
	}

	return &placeOrderInfo, nil
}

// HasLastPlacedOrder returns if a order has been previously placed
func (os *OrderService) HasLastPlacedOrder(ctx context.Context) bool {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/HasLastPlacedOrder")
	defer span.End()

	lastPlaced, err := os.LastPlacedOrder(ctx)

	return lastPlaced != nil && err == nil
}

// ClearLastPlacedOrder clears the last placed cart, this can be useful if an cart / order is finished
func (os *OrderService) ClearLastPlacedOrder(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/ClearLastPlacedOrder")
	defer span.End()

	session := web.SessionFromContext(ctx)
	session.Delete(LastPlacedOrderSessionKey)
}

// LastPlacedOrCurrentCart returns the decorated cart of the last placed order if there is one if not return the current cart
func (os *OrderService) LastPlacedOrCurrentCart(ctx context.Context) (*decorator.DecoratedCart, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/LastPlacedOrCurrentCart")
	defer span.End()

	lastPlacedOrder, err := os.LastPlacedOrder(ctx)
	if err != nil {
		os.logger.Warn("couldn't get last placed order:", err)

		return nil, err
	}

	if lastPlacedOrder != nil {
		// cart has been placed early use it
		return os.decoratedCartFactory.Create(ctx, lastPlacedOrder.Cart), nil
	}

	// cart wasn't placed early, fetch it from service
	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, web.SessionFromContext(ctx))
	if err != nil {
		os.logger.WithContext(ctx).Error("ViewDecoratedCart Error:", err)

		return nil, err
	}

	return decoratedCart, nil
}

// placeOrderWithPaymentProcessing after generating the decorated cart, the place order flow
// is the same for the interface functions, therefore the common flow is placed in this private helper function
func (os *OrderService) placeOrderWithPaymentProcessing(ctx context.Context, decoratedCart *decorator.DecoratedCart,
	session *web.Session) (*PlaceOrderInfo, error) {
	ctx, span := trace.StartSpan(ctx, "checkout/OrderService/placeOrderWithPaymentProcessing")
	defer span.End()

	if !decoratedCart.Cart.IsPaymentSelected() {
		// record noPaymentSelectionCount metric
		stats.Record(ctx, noPaymentSelectionCount.M(1))
		os.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error Gateway not in carts PaymentSelection")

		return nil, errors.New("no payment gateway selected")
	}

	validationResult := os.cartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		// record cartValidationFailCount metric
		stats.Record(ctx, cartValidationFailCount.M(1))
		os.logger.WithContext(ctx).Warn("Try to place an invalid cart")

		return nil, errors.New("cart is invalid")
	}

	gateway, err := os.GetPaymentGateway(ctx, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		// record paymentGatewayNotFoundCount metric
		stats.Record(ctx, paymentGatewayNotFoundCount.M(1))
		os.logger.WithContext(ctx).Error(fmt.Sprintf("cart.checkoutcontroller.submitaction: Error %v  Gateway: %v", err, decoratedCart.Cart.PaymentSelection.Gateway()))

		return nil, errors.New("selected gateway not available")
	}

	flowStatus, err := gateway.FlowStatus(ctx, &decoratedCart.Cart, PaymentFlowStandardCorrelationID)
	if err != nil {
		// record paymentFlowStatusErrorCount metric
		stats.Record(ctx, paymentFlowStatusErrorCount.M(1))

		return nil, err
	}

	if flowStatus.Status == paymentDomain.PaymentFlowStatusFailed || flowStatus.Status == paymentDomain.PaymentFlowStatusCancelled {
		// record paymentFlowStatusFailedCanceledCount metric
		stats.Record(ctx, paymentFlowStatusFailedCanceledCount.M(1))
		os.logger.WithContext(ctx).Info("cart.checkoutcontroller.submitaction: PaymentFlowStatusFailed or PaymentFlowStatusCancelled: Error ", flowStatus.Error)

		return nil, flowStatus.Error
	}

	if flowStatus.Status == paymentDomain.PaymentFlowStatusAborted {
		// record paymentFlowStatusAbortedCount metric
		stats.Record(ctx, paymentFlowStatusAbortedCount.M(1))
		os.logger.WithContext(ctx).Info("cart.checkoutcontroller.submitaction: PaymentFlowStatusAborted: Error ", flowStatus.Error)

		return nil, flowStatus.Error
	}

	cartPayment, err := gateway.OrderPaymentFromFlow(ctx, &decoratedCart.Cart, PaymentFlowStandardCorrelationID)
	if err != nil {
		// record orderPaymentFromFlowErrorCount metric
		stats.Record(ctx, orderPaymentFromFlowErrorCount.M(1))

		return nil, err
	}

	placedOrderInfos, err := os.cartService.PlaceOrderWithCart(ctx, session, &decoratedCart.Cart, cartPayment)
	if err != nil {
		// record placeOrderFailCount metric
		stats.Record(ctx, placeOrderFailCount.M(1))
		os.logger.WithContext(ctx).Error("Error during place Order: " + err.Error())

		return nil, err
	}

	os.logger.WithContext(ctx).Info("Placed Order: ", placedOrderInfos)

	placeOrderInfo := os.preparePlaceOrderInfo(ctx, decoratedCart.Cart, placedOrderInfos, *cartPayment)
	os.storeLastPlacedOrder(ctx, placeOrderInfo)

	if flowStatus.Status != paymentDomain.PaymentFlowStatusCompleted {
		err = gateway.ConfirmResult(ctx, &decoratedCart.Cart, cartPayment)
		if err != nil {
			os.logger.WithContext(ctx).Error("Error during gateway.ConfirmResult: " + err.Error())

			return nil, err
		}
	}

	// record placeOrderSuccessCount metric
	stats.Record(ctx, placeOrderSuccessCount.M(1))

	return placeOrderInfo, nil
}

func (os *OrderService) preparePlaceOrderInfo(ctx context.Context, currentCart cart.Cart, placedOrderInfos placeorder.PlacedOrderInfos, cartPayment placeorder.Payment) *PlaceOrderInfo {
	_, span := trace.StartSpan(ctx, "checkout/OrderService/preparePlaceOrderInfo")
	defer span.End()

	email := currentCart.GetContactMail()

	placeOrderInfo := &PlaceOrderInfo{
		ContactEmail: email,
		PlacedOrders: placedOrderInfos,
		Cart:         currentCart,
	}

	for _, transaction := range cartPayment.Transactions {
		placeOrderInfo.PaymentInfos = append(placeOrderInfo.PaymentInfos, PlaceOrderPaymentInfo{
			Gateway:         cartPayment.Gateway,
			Method:          transaction.Method,
			PaymentProvider: transaction.PaymentProvider,
			Title:           transaction.Title,
			Amount:          transaction.AmountPayed,
			CreditCardInfo:  transaction.CreditCardInfo,
		})
	}

	return placeOrderInfo
}
