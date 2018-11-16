package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/checkout/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/gorilla/sessions"
)

type (
	// OrderService defines the order service
	OrderService struct {
		SourcingEngine      *domain.SourcingEngine           `inject:""`
		PaymentService      *PaymentService                  `inject:""`
		Logger              flamingo.Logger                  `inject:""`
		CartService         *application.CartService         `inject:""`
		CartReceiverService *application.CartReceiverService `inject:""`
		DeliveryInfoBuilder cart.DeliveryInfoBuilder         `inject:""`
	}
)

// SetSources sets sources for sessions carts items
func (os *OrderService) SetSources(ctx context.Context, session *sessions.Session) error {
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return err
	}
	err = os.SourcingEngine.SetSourcesForCartItems(ctx, session, decoratedCart)
	if err != nil {
		os.Logger.WithField("category", "checkout.orderService").Error("Error while getting sources: %v", err)
		return errors.New("Error while setting sources.")
	}
	return nil
}

// PlaceOrder places the order
func (os *OrderService) PlaceOrder(ctx context.Context, session *sessions.Session, decoratedCart *cart.DecoratedCart, payment *cart.CartPayment) (cart.PlacedOrderInfos, error) {
	validationResult := os.CartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return nil, errors.New("Cart is Invalid.")
	}
	return os.CartService.PlaceOrder(ctx, session, payment)
}

// CurrentCartSaveInfos saves additional informations on current cart
func (os *OrderService) CurrentCartSaveInfos(ctx context.Context, session *sessions.Session, billingAddress *cart.Address, shippingAddress *cart.Address, purchaser *cart.Person, additionalCustomData map[string]string) error {
	os.Logger.Debug("CurrentCartSaveInfos call billingAddress:%v shippingAddress:%v payment:%v", billingAddress, shippingAddress)

	if billingAddress == nil {
		os.Logger.Warn("CurrentCartSaveInfos called without billing address")
		return errors.New("Billing Address is missing")
	}
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.Logger.Error("CurrentCartSaveInfos GetDecoratedCart Error %v", err)
		return err
	}

	//update Billing
	err = os.CartService.UpdateBillingAddress(ctx, session, billingAddress)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder UpdateBillingAddress Error %v", err)
		return err
	}

	//Update ShippingAddress on ALL Deliveries in the Cart if given
	// Maybe later we need to support different shipping addresses in the Checkout
	if shippingAddress != nil {
		for _, d := range decoratedCart.Cart.Deliveries {
			newDeliveryInfo := d.DeliveryInfo
			newDeliveryInfo.DeliveryLocation.Address = shippingAddress
			err = os.CartService.UpdateDeliveryInfo(ctx, session, d.DeliveryInfo.Code, newDeliveryInfo)
			if err != nil {
				os.Logger.Error("OnStepCurrentCartPlaceOrder UpdateDeliveryInfosAndBilling Error %v", err)
				return err
			}
		}

	}

	//Update Purchaser
	err = os.CartService.UpdatePurchaser(ctx, session, purchaser, additionalCustomData)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder UpdatePurchaser Error %v", err)
		return err
	}

	//After setting DeliveryInfos - call SourcingEnginge (this will reload the cart and update all items!)
	err = os.SetSources(ctx, session)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder SetSources Error %v", err)
		return err
	}
	return nil
}

//CurrentCartPlaceOrder - probably the best choice for a simple checkout
// Assumptions: Only one BuildDeliveryInfo is used on the cart!
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *sessions.Session, payment cart.CartPayment) (cart.PlacedOrderInfos, error) {
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return nil, err
	}

	validationResult := os.CartService.ValidateCart(ctx, session, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return nil, errors.New("Cart is Invalid.")
	}

	placedOrderInfos, err := os.PlaceOrder(ctx, session, decoratedCart, &payment)

	if err != nil {
		os.Logger.WithField("category", "checkout.orderService").Error("Error during place Order: %v", err)
		return nil, errors.New("Error while placing the order. Please contact customer support.")
	}
	return placedOrderInfos, nil
}
