package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/checkout/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
)

type (

	// PaymentService
	OrderService struct {
		SourcingEngine      *domain.SourcingEngine           `inject:""`
		PaymentService      *PaymentService                  `inject:""`
		Logger              flamingo.Logger                  `inject:""`
		CartService         *application.CartService         `inject:""`
		CartReceiverService *application.CartReceiverService `inject:""`
		DeliveryInfoBuilder cart.DeliveryInfoBuilder         `inject:""`
	}
)

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

func (os *OrderService) PlaceOrder(ctx context.Context, session *sessions.Session, decoratedCart *cart.DecoratedCart, payment *cart.CartPayment) (orderid string, orderError error) {
	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}
	return os.CartService.PlaceOrder(ctx, session, payment)
}

func (os *OrderService) CurrentCartSaveInfos(ctx context.Context, session *sessions.Session, billingAddress *cart.Address, shippingAddress *cart.Address, purchaser *cart.Person) error {
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

	updateCommands, err := os.DeliveryInfoBuilder.BuildDeliveryInfoUpdateCommand(web.ToContext(ctx), decoratedCart)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder BuildDeliveryInfoUpdateCommand Error %v", err)
		return err
	}

	//If an Address is given - add it to every DeliveryInfo(s)
	if shippingAddress != nil {
		if len(updateCommands) == 0 {
			os.Logger.Warn("OnStepCurrentCartPlaceOrder Cart has no DeliveryInfoUpdates Build - cannot set shippingAddress")
			return errors.New("No DeliveryInfos Build - cannot set shippingAddress")
		}
		for k, _ := range updateCommands {
			updateCommands[k].DeliveryInfo.DeliveryLocation.Address = shippingAddress
		}
	}
	err = os.CartService.UpdateDeliveryInfosAndBilling(ctx, session, billingAddress, updateCommands)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder UpdateDeliveryInfosAndBilling Error %v", err)
		return err
	}

	err = os.CartService.UpdatePurchaser(ctx, session, purchaser, nil)
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

//OnStepCurrentCartPlaceOrder - probably the best choice for a simple checkout
// Assumptions: Only one BuildDeliveryInfo is used on the cart!
func (os *OrderService) CurrentCartPlaceOrder(ctx context.Context, session *sessions.Session, payment cart.CartPayment) (orderid string, orderError error) {
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return "", err
	}

	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}

	orderid, orderError = os.PlaceOrder(ctx, session, decoratedCart, &payment)

	if orderError != nil {
		os.Logger.WithField("category", "checkout.orderService").Error("Error during place Order: %v", err)
		return "", errors.New("Error while placing the order. Please contact customer support.")
	}
	return orderid, nil
}
