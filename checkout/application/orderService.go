package application

import (
	"context"
	"errors"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// OrderService defines the order service
	OrderService struct {
		sourcingEngine      *domain.SourcingEngine
		paymentService      *PaymentService
		logger              flamingo.Logger
		cartService         *application.CartService
		cartReceiverService *application.CartReceiverService
		deliveryInfoBuilder cart.DeliveryInfoBuilder
	}
)

// Inject dependencies
func (os *OrderService) Inject(
	SourcingEngine *domain.SourcingEngine,
	PaymentService *PaymentService,
	Logger flamingo.Logger,
	CartService *application.CartService,
	CartReceiverService *application.CartReceiverService,
	DeliveryInfoBuilder cart.DeliveryInfoBuilder,
) {
	os.sourcingEngine = SourcingEngine
	os.paymentService = PaymentService
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
		return errors.New("Error while setting sources.")
	}
	return nil
}

// PlaceOrder places the order
func (os *OrderService) PlaceOrder(ctx context.Context, session *web.Session, decoratedCart *cart.DecoratedCart, payment *cart.Payment) (cart.PlacedOrderInfos, error) {
	validationResult := os.cartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		os.logger.Warn("Try to place an invalid cart")
		return nil, errors.New("Cart is Invalid.")
	}
	return os.cartService.PlaceOrder(ctx, session, payment)
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

//CurrentCartPlaceOrder - probably the best choice for a simple checkout
// Assumptions: Only one BuildDeliveryInfo is used on the cart!
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *web.Session, payment cart.Payment) (cart.PlacedOrderInfos, error) {
	decoratedCart, err := os.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return nil, err
	}

	validationResult := os.cartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		os.logger.Warn("Try to place an invalid cart")
		return nil, errors.New("Cart is Invalid.")
	}

	placedOrderInfos, err := os.PlaceOrder(ctx, session, decoratedCart, &payment)

	if err != nil {
		os.logger.WithField("category", "checkout.orderService").Error("Error during place Order: %v", err)
		return nil, errors.New("Error while placing the order. Please contact customer support.")
	}
	return placedOrderInfos, nil
}
