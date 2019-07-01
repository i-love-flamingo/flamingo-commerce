package application

import (
	"context"
	"errors"
	"fmt"

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
	}

	// PlaceOrderInfo struct defines the data of payments on placed orders
	PlaceOrderInfo struct {
		PaymentInfos []PlaceOrderPaymentInfo
		PlacedOrders placeorder.PlacedOrderInfos
		ContactEmail string
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
	//PaymentFlowStandardCorrelationID - used as correlationid for the start of the payment (session scoped)
	PaymentFlowStandardCorrelationID = "checkout"
)

var orderFailedStat = stats.Int64("flamingo-commerce/orderfailed", "my stat records 1 occurences per error", stats.UnitDimensionless)

func init() {
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
) {
	os.sourcingEngine = SourcingEngine
	os.logger = logger.WithField(flamingo.LogKeyCategory, "checkout.OrderService").WithField(flamingo.LogKeyModule, "checkout")
	os.cartService = CartService
	os.cartReceiverService = CartReceiverService
	os.webCartPaymentGateways = webCartPaymentGatewayProvider()
	os.deliveryInfoBuilder = DeliveryInfoBuilder
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

	//update Billing
	err = os.cartService.UpdateBillingAddress(ctx, session, billingAddress)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder UpdateBillingAddress Error %v", err)
		return err
	}

	//Update ShippingAddress on ALL Deliveries in the Cart if given
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

	//Update Purchaser
	err = os.cartService.UpdatePurchaser(ctx, session, purchaser, additionalData)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder UpdatePurchaser Error %v", err)
		return err
	}

	//After setting DeliveryInfos - call SourcingEnginge (this will reload the cart and update all items!)
	err = os.SetSources(ctx, session)
	if err != nil {
		os.logger.WithContext(ctx).Error("OnStepCurrentCartPlaceOrder SetSources Error %v", err)
		return err
	}
	return nil
}

//CurrentCartPlaceOrder - use to place the current cart
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *web.Session, payment placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	// use a background context from here on to prevent the place order canceled by context cancel
	placeOrderContext := web.ContextWithSession(
		context.Background(),
		web.SessionFromContext(ctx),
	)

	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(placeOrderContext, session)

	if err != nil {
		// record failcount metric
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return nil, err
	}

	validationResult := os.cartService.ValidateCart(placeOrderContext, session, decoratedCart)
	if !validationResult.IsValid() {
		// record failcount metric
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Warn("Try to place an invalid cart")
		return nil, errors.New("cart is invalid")
	}

	placedOrderInfos, err := os.cartService.PlaceOrder(placeOrderContext, session, &payment)

	if err != nil {
		// record failcount metric
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Error("Error during place Order:" + err.Error())
		return nil, errors.New("error while placing the order. please contact customer support")
	}
	return placedOrderInfos, nil
}

// GetPaymentGateway helper
func (os *OrderService) GetPaymentGateway(ctx context.Context, paymentGatewayCode string) (interfaces.WebCartPaymentGateway, error) {

	gateway, ok := os.webCartPaymentGateways[paymentGatewayCode]
	if !ok {
		return nil, errors.New("Payment gateway " + paymentGatewayCode + " not found")
	}

	return gateway, nil
}

//GetAvailablePaymentGateways - returns the list of registered WebCartPaymentGateway
func (os *OrderService) GetAvailablePaymentGateways(ctx context.Context) map[string]interfaces.WebCartPaymentGateway {
	return os.webCartPaymentGateways
}

//CurrentCartPlaceOrderWithPaymentProcessing - use to place the current cart
func (os *OrderService) CurrentCartPlaceOrderWithPaymentProcessing(ctx context.Context, session *web.Session) (*PlaceOrderInfo, error) {
	// use a background context from here on to prevent the place order canceled by context cancel
	placeOrderContext := web.ContextWithSession(
		context.Background(),
		web.SessionFromContext(ctx),
	)

	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(placeOrderContext, session)
	if !decoratedCart.Cart.IsPaymentSelected() {
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Error("cart.checkoutcontroller.submitaction: Error Gateway not in carts PaymentSelection")
		return nil, errors.New("no payment gateway selected")
	}

	validationResult := os.cartService.ValidateCart(placeOrderContext, session, decoratedCart)
	if !validationResult.IsValid() {
		// record failcount metric
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Warn("Try to place an invalid cart")
		return nil, errors.New("cart is invalid")
	}

	gateway, err := os.GetPaymentGateway(placeOrderContext, decoratedCart.Cart.PaymentSelection.Gateway())
	if err != nil {
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Error(fmt.Sprintf("cart.checkoutcontroller.submitaction: Error %v  Gateway: %v", err,decoratedCart.Cart.PaymentSelection.Gateway()))
		return nil, errors.New("selected gateway not available")
	}

	cartPayment, err := gateway.GetFlowResult(placeOrderContext, &decoratedCart.Cart, PaymentFlowStandardCorrelationID)
	if err != nil {
		return nil, err
	}
	err = gateway.ConfirmResult(placeOrderContext, &decoratedCart.Cart, cartPayment)
	if err != nil {
		return nil, err
	}

	placedOrderInfos, err := os.cartService.PlaceOrder(placeOrderContext, session, cartPayment)
	if err != nil {
		// record failcount metric
		stats.Record(placeOrderContext, orderFailedStat.M(1))
		os.logger.WithContext(placeOrderContext).Error("Error during place Order:" + err.Error())
		return nil, errors.New("error while placing the order. please contact customer support")
	}

	email := os.GetContactMail(decoratedCart.Cart)

	placeOrderInfo := PlaceOrderInfo{
		ContactEmail: email,
		PlacedOrders: placedOrderInfos,
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
	return &placeOrderInfo, nil
}

// GetContactMail helper with fallback
func (os *OrderService) GetContactMail(cart cart.Cart) string {
	//Get Email from either the cart
	shippingEmail := cart.GetMainShippingEMail()
	if shippingEmail == "" {
		shippingEmail = cart.BillingAdress.Email
	}
	return shippingEmail
}
