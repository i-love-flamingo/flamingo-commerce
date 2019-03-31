package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

type (
	// OrderService defines the order service
	OrderService struct {
		sourcingEngine      *domain.SourcingEngine
		logger              flamingo.Logger
		cartService         *application.CartService
		cartReceiverService *application.CartReceiverService
		deliveryInfoBuilder cart.DeliveryInfoBuilder
	}
)

var orderFailedStat = stats.Int64("flamingo-commerce/orderfailed", "my stat records 1 occurences per error", stats.UnitDimensionless)

func init() {
	opencensus.View("flamingo-commerce/orderfailed/count", orderFailedStat, view.Count())
}

// Inject dependencies
func (os *OrderService) Inject(
	SourcingEngine *domain.SourcingEngine,
	Logger flamingo.Logger,
	CartService *application.CartService,
	CartReceiverService *application.CartReceiverService,
	DeliveryInfoBuilder cart.DeliveryInfoBuilder,
) {
	os.sourcingEngine = SourcingEngine
	os.logger = Logger
	os.cartService = CartService
	os.cartReceiverService = CartReceiverService
	os.deliveryInfoBuilder = DeliveryInfoBuilder
}

// SetSources sets sources for sessions carts items
func (os *OrderService) SetSources(ctx context.Context, session *web.Session) error {
	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return err
	}
	err = os.sourcingEngine.SetSourcesForCartItems(ctx, session, decoratedCart)
	if err != nil {
		os.logger.WithField("category", "checkout.orderService").Error("Error while getting sources: %v", err)
		return errors.New("error while setting sources")
	}
	return nil
}

// CurrentCartSaveInfos saves additional informations on current cart
func (os *OrderService) CurrentCartSaveInfos(ctx context.Context, session *web.Session, billingAddress *cart.Address, shippingAddress *cart.Address, purchaser *cart.Person, additionalData *cart.AdditionalData) error {
	os.logger.Debug("CurrentCartSaveInfos call billingAddress:%v shippingAddress:%v payment:%v", billingAddress, shippingAddress)

	if billingAddress == nil {
		os.logger.Warn("CurrentCartSaveInfos called without billing address")
		return errors.New("Billing Address is missing")
	}
	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.logger.Error("CurrentCartSaveInfos GetDecoratedCart Error %v", err)
		return err
	}

	//update Billing
	err = os.cartService.UpdateBillingAddress(ctx, session, billingAddress)
	if err != nil {
		os.logger.Error("OnStepCurrentCartPlaceOrder UpdateBillingAddress Error %v", err)
		return err
	}

	//Update ShippingAddress on ALL Deliveries in the Cart if given
	// Maybe later we need to support different shipping addresses in the Checkout
	if shippingAddress != nil {
		for _, d := range decoratedCart.Cart.Deliveries {
			newDeliveryInfoUpdateCommand := cart.DeliveryInfoUpdateCommand{
				DeliveryInfo: d.DeliveryInfo,
			}
			newDeliveryInfoUpdateCommand.DeliveryInfo.DeliveryLocation.Address = shippingAddress
			err = os.cartService.UpdateDeliveryInfo(ctx, session, d.DeliveryInfo.Code, newDeliveryInfoUpdateCommand)
			if err != nil {
				os.logger.Error("OnStepCurrentCartPlaceOrder UpdateDeliveryInfosAndBilling Error %v", err)
				return err
			}
		}

	}

	//Update Purchaser
	err = os.cartService.UpdatePurchaser(ctx, session, purchaser, additionalData)
	if err != nil {
		os.logger.Error("OnStepCurrentCartPlaceOrder UpdatePurchaser Error %v", err)
		return err
	}

	//After setting DeliveryInfos - call SourcingEnginge (this will reload the cart and update all items!)
	err = os.SetSources(ctx, session)
	if err != nil {
		os.logger.Error("OnStepCurrentCartPlaceOrder SetSources Error %v", err)
		return err
	}
	return nil
}

//CurrentCartPlaceOrder - use to place the current cart
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *web.Session, payment cart.Payment) (cart.PlacedOrderInfos, error) {
	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, session)

	if err != nil {
		// record failcount metric
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return nil, err
	}

	validationResult := os.cartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		// record failcount metric
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.Warn("Try to place an invalid cart")
		return nil, errors.New("cart is invalid")
	}

	placedOrderInfos, err := os.cartService.PlaceOrder(ctx, session, &payment)

	if err != nil {
		// record failcount metric
		stats.Record(ctx, orderFailedStat.M(1))
		os.logger.WithField("category", "checkout.orderService").Error("Error during place Order:" + err.Error())
		return nil, errors.New("error while placing the order. please contact customer support")
	}
	return placedOrderInfos, nil
}
