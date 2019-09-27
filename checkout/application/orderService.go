package application

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	paymentDomain "flamingo.me/flamingo-commerce/v3/payment/domain"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// OrderService defines the order service
	OrderService struct {
		sourcingEngine         *domain.SourcingEngine
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

	//PlaceOrderPaymentInfo holding payment infos
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

var orderFailedStat = stats.Int64("flamingo-commerce/orderfailed", "my stat records 1 occurrence per error", stats.UnitDimensionless)

func init() {
	gob.Register(PlaceOrderInfo{})
	opencensus.View("flamingo-commerce/orderfailed/count", orderFailedStat, view.Count())
}

// Inject dependencies
func (os *OrderService) Inject(
	SourcingEngine *domain.SourcingEngine,
	logger flamingo.Logger,
	CartService *application.CartService,
	CartReceiverService *application.CartReceiverService,
	DeliveryInfoBuilder cart.DeliveryInfoBuilder,
	webCartPaymentGatewayProvider interfaces.WebCartPaymentGatewayProvider,
	decoratedCartFactory *decorator.DecoratedCartFactory,
) {
	os.sourcingEngine = SourcingEngine
	os.logger = logger.WithField(flamingo.LogKeyCategory, "checkout.OrderService").WithField(flamingo.LogKeyModule, "checkout")
	os.cartService = CartService
	os.cartReceiverService = CartReceiverService
	os.webCartPaymentGateways = webCartPaymentGatewayProvider()
	os.deliveryInfoBuilder = DeliveryInfoBuilder
	os.decoratedCartFactory = decoratedCartFactory
}

// SetSources sets sources for sessions carts items
func (os *OrderService) SetSources(ctx context.Context, session *web.Session) error {
	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return err
	}
	err = os.sourcingEngine.SetSourcesForCartItems(ctx, session, decoratedCart)
	if err != nil {
		os.logger.WithContext(ctx).WithField("category", "checkout.orderService").Error("Error while getting sources: %v", err)
		return errors.New("error while setting sources")
	}
	return nil
}

// CurrentCartSaveInfos saves additional information on current cart
func (os *OrderService) CurrentCartSaveInfos(ctx context.Context, session *web.Session, billingAddress *cart.Address, shippingAddress *cart.Address, purchaser *cart.Person, additionalData *cart.AdditionalData) error {
	os.logger.WithContext(ctx).Debug("CurrentCartSaveInfos call billingAddress:%v shippingAddress:%v payment:%v", billingAddress, shippingAddress)

	if billingAddress == nil {
		os.logger.WithContext(ctx).Warn("CurrentCartSaveInfos called without billing address")
		return errors.New("Billing Address is missing")
	}

	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.logger.WithContext(ctx).Error("CurrentCartSaveInfos GetDecoratedCart Error %v", err)
		return err
	}

	// update Billing
	err = os.cartService.UpdateBillingAddress(ctx, session, billingAddress)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder UpdateBillingAddress Error %v", err)
		return err
	}

	// Update ShippingAddress on ALL Deliveries in the Cart if given
	// Maybe later we need to support different shipping addresses in the Checkout
	if shippingAddress != nil {
		for _, d := range decoratedCart.Cart.Deliveries {
			newDeliveryInfoUpdateCommand := cart.CreateDeliveryInfoUpdateCommand(d.DeliveryInfo)
			newDeliveryInfoUpdateCommand.DeliveryInfo.DeliveryLocation.Address = shippingAddress
			err = os.cartService.UpdateDeliveryInfo(ctx, session, d.DeliveryInfo.Code, newDeliveryInfoUpdateCommand)
			if err != nil {
				os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder UpdateDeliveryInfosAndBilling Error %v", err)
				return err
			}
		}

	}

	// Update Purchaser
	err = os.cartService.UpdatePurchaser(ctx, session, purchaser, additionalData)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder UpdatePurchaser Error %v", err)
		return err
	}

	// After setting DeliveryInfos - call SourcingEnginge (this will reload the cart and update all items!)
	err = os.SetSources(ctx, session)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder SetSources Error %v", err)
		return err
	}
	return nil
}

// GetPaymentGateway tries to get the supplied payment gateway by code from the registered payment gateways
func (os *OrderService) GetPaymentGateway(ctx context.Context, paymentGatewayCode string) (interfaces.WebCartPaymentGateway, error) {
	gateway, ok := os.webCartPaymentGateways[paymentGatewayCode]
	if !ok {
		return nil, errors.New("Payment gateway " + paymentGatewayCode + " not found")
	}

	return gateway, nil
}

// GetAvailablePaymentGateways returns the list of registered WebCartPaymentGateway
func (os *OrderService) GetAvailablePaymentGateways(ctx context.Context) map[string]interfaces.WebCartPaymentGateway {
	return os.webCartPaymentGateways
}

// CurrentCartPlaceOrder places the current cart without additional payment processing
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *web.Session, cartPayment placeorder.Payment) (*PlaceOrderInfo, error) {
	var info *PlaceOrderInfo
	var err error
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = func() (*PlaceOrderInfo, error) {
			decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(placeOrderContext, session)
			if err != nil {
				// record fail count metric
				stats.Record(placeOrderContext, orderFailedStat.M(1))
				os.logger.WithContext(placeOrderContext).Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
				return nil, err
			}

			validationResult := os.cartService.ValidateCart(placeOrderContext, session, decoratedCart)
			if !validationResult.IsValid() {
				// record fail count metric
				stats.Record(placeOrderContext, orderFailedStat.M(1))
				os.logger.WithContext(placeOrderContext).Warn("Try to place an invalid cart")
				return nil, errors.New("cart is invalid")
			}

			placedOrderInfos, err := os.cartService.PlaceOrder(placeOrderContext, session, &cartPayment)
			if err != nil {
				// record fail count metric
				stats.Record(placeOrderContext, orderFailedStat.M(1))
				os.logger.WithContext(placeOrderContext).Error("Error during place Order:" + err.Error())
				return nil, errors.New("error while placing the order. please contact customer support")
			}

			placeOrderInfo := os.preparePlaceOrderInfo(ctx, decoratedCart.Cart, placedOrderInfos, cartPayment)
			os.storeLastPlacedOrder(ctx, placeOrderInfo)

			return placeOrderInfo, nil
		}()
	})
	return info, err
}

// CancelOrder cancels an previously placed order and returns the cart with the order content
func (os *OrderService) CancelOrder(ctx context.Context, session *web.Session, order *PlaceOrderInfo) (*cart.Cart, error) {
	return os.cartService.CancelOrder(ctx, session, order.PlacedOrders, order.Cart)
}

// CurrentCartPlaceOrderWithPaymentProcessing places the current cart which is fetched from the context
func (os *OrderService) CurrentCartPlaceOrderWithPaymentProcessing(ctx context.Context, session *web.Session) (*PlaceOrderInfo, error) {
	var info *PlaceOrderInfo
	var err error
	// use a background context from here on to prevent the place order canceled by context cancel
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = func() (*PlaceOrderInfo, error) {
			// fetch decorated cart either via cache or freshly from cart receiver service
			decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(placeOrderContext, session)
			if err != nil {
				// record fail count metric
				stats.Record(placeOrderContext, orderFailedStat.M(1))
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
	var info *PlaceOrderInfo
	var err error
	// use a background context from here on to prevent the place order canceled by context cancel
	web.RunWithDetachedContext(ctx, func(placeOrderContext context.Context) {
		info, err = os.placeOrderWithPaymentProcessing(placeOrderContext, decoratedCart, session)
	})
	return info, err
}

// GetContactMail returns the contact mail from the shipping address with fall back to the billing address
func (os *OrderService) GetContactMail(cart cart.Cart) string {
	//Get Email from either the cart
	shippingEmail := cart.GetMainShippingEMail()
	if shippingEmail == "" {
		shippingEmail = cart.BillingAddress.Email
	}
	return shippingEmail
}

// storeLastPlacedOrder stores the last placed order/cart in the session
func (os *OrderService) storeLastPlacedOrder(ctx context.Context, info *PlaceOrderInfo) {
	session := web.SessionFromContext(ctx)

	_ = session.Store(LastPlacedOrderSessionKey, info)
}

// LastPlacedOrder returns the last placed order/cart if available
func (os *OrderService) LastPlacedOrder(ctx context.Context) (*PlaceOrderInfo, error) {
	session := web.SessionFromContext(ctx)

	lastPlacedOrder, found := session.Load(LastPlacedOrderSessionKey)
	if found == false {
		return nil, nil
	}

	placeOrderInfo, ok := lastPlacedOrder.(PlaceOrderInfo)
	if ok == false {
		return nil, errors.New("placeOrderInfo couldn't be received from session")
	}

	return &placeOrderInfo, nil
}

// HasLastPlacedOrder returns if a order has been previously placed
func (os *OrderService) HasLastPlacedOrder(ctx context.Context) bool {
	lastPlaced, err := os.LastPlacedOrder(ctx)
	return lastPlaced != nil && err == nil
}

// ClearLastPlacedOrder clears the last placed cart, this can be useful if an cart / order is finished
func (os *OrderService) ClearLastPlacedOrder(ctx context.Context) {
	session := web.SessionFromContext(ctx)
	session.Delete(LastPlacedOrderSessionKey)
}

// LastPlacedOrCurrentCart returns the decorated cart of the last placed order if there is one if not return the current cart
func (os *OrderService) LastPlacedOrCurrentCart(ctx context.Context) (*decorator.DecoratedCart, error) {
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
	if !decoratedCart.Cart.IsPaymentSelected() {
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.WithContext(ctx).Error("cart.checkoutcontroller.submitaction: Error Gateway not in carts PaymentSelection")
		return nil, errors.New("no payment gateway selected")
	}

	validationResult := os.cartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		// record fail count metric
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.WithContext(ctx).Warn("Try to place an invalid cart")
		return nil, errors.New("cart is invalid")
	}

	gateway, err := os.GetPaymentGateway(ctx, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		// record fail count metric
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.WithContext(ctx).Error(fmt.Sprintf("cart.checkoutcontroller.submitaction: Error %v  Gateway: %v", err, decoratedCart.Cart.PaymentSelection.Gateway()))
		return nil, errors.New("selected gateway not available")
	}

	flowStatus, err := gateway.FlowStatus(ctx, &decoratedCart.Cart, PaymentFlowStandardCorrelationID)
	if err != nil {
		// record fail count metric
		stats.Record(ctx, orderFailedStat.M(1))
		return nil, err
	}

	if flowStatus.Status == paymentDomain.PaymentFlowStatusFailed || flowStatus.Status == paymentDomain.PaymentFlowStatusCancelled {
		// record fail count metric
		stats.Record(ctx, orderFailedStat.M(1))
		return nil, flowStatus.Error
	}

	if flowStatus.Status == paymentDomain.PaymentFlowStatusAborted {
		return nil, flowStatus.Error
	}

	cartPayment, err := gateway.OrderPaymentFromFlow(ctx, &decoratedCart.Cart, PaymentFlowStandardCorrelationID)
	if err != nil {
		// record fail count metric
		stats.Record(ctx, orderFailedStat.M(1))
		return nil, err
	}

	placedOrderInfos, err := os.cartService.PlaceOrder(ctx, session, cartPayment)
	if err != nil {
		// record fail count metric
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.WithContext(ctx).Error("Error during place Order:" + err.Error())
		return nil, err
	}

	placeOrderInfo := os.preparePlaceOrderInfo(ctx, decoratedCart.Cart, placedOrderInfos, *cartPayment)
	os.storeLastPlacedOrder(ctx, placeOrderInfo)

	err = gateway.ConfirmResult(ctx, &decoratedCart.Cart, cartPayment)
	if err != nil {
		os.logger.WithContext(ctx).Error("Error during gateway.ConfirmResult:" + err.Error())
	}

	return placeOrderInfo, nil
}

func (os *OrderService) preparePlaceOrderInfo(ctx context.Context, currentCart cart.Cart, placedOrderInfos placeorder.PlacedOrderInfos, cartPayment placeorder.Payment) *PlaceOrderInfo {
	email := os.GetContactMail(currentCart)

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
