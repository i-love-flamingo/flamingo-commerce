package application

import (
	"errors"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/checkout/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
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

func (os *OrderService) SetSources(ctx web.Context) error {
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		os.Logger.Errorf("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return err
	}
	err = os.SourcingEngine.SetSourcesForCartItems(ctx, decoratedCart)
	if err != nil {
		os.Logger.WithField("category", "checkout.orderService").Errorf("Error while getting sources: %v", err)
		return errors.New("Error while setting sources.")
	}
	return nil
}

func (os *OrderService) PlaceOrder(ctx web.Context, decoratedCart *cart.DecoratedCart, payment *cart.CartPayment) (orderid string, orderError error) {
	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}
	return os.CartService.PlaceOrder(ctx, payment)
}

//OnStepCurrentCartPlaceOrder - probably the best choice for a simple checkout
// Assumptions: Only one BuildDeliveryInfo is used on the cart!
func (os *OrderService) OnStepCurrentCartPlaceOrder(ctx web.Context, billingAddress *cart.Address, shippingAddress *cart.Address, payment *cart.CartPayment, purchaser *cart.Person) (orderid string, orderError error) {
	os.Logger.Debugf("OnStepCurrentCartPlaceOrder call billingAddress:%v shippingAddress:%v payment:%v", billingAddress, shippingAddress, payment)

	if billingAddress == nil {
		os.Logger.Warn("OnStepCurrentCartPlaceOrder called without billing address")
		return "", errors.New("Billing Address is missing")
	}
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		os.Logger.Errorf("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return "", err
	}

	updateCommands, err := os.DeliveryInfoBuilder.BuildDeliveryInfoUpdateCommand(ctx, decoratedCart)
	if err != nil {
		os.Logger.Errorf("OnStepCurrentCartPlaceOrder BuildDeliveryInfoUpdateCommand Error %v", err)
		return "", err
	}

	//If an Address is given - add it to every DeliveryInfo(s)
	if shippingAddress != nil {
		if len(updateCommands) == 0 {
			os.Logger.Warnf("OnStepCurrentCartPlaceOrder Cart has no DeliveryInfoUpdates Build - cannot set shippingAddress")
			return "", errors.New("No DeliveryInfos Build - cannot set shippingAddress")
		}
		for k, _ := range updateCommands {
			updateCommands[k].DeliveryInfo.DeliveryLocation.Address = shippingAddress
		}
	}
	err = os.CartService.UpdateDeliveryInfosAndBilling(ctx, billingAddress, updateCommands)
	if err != nil {
		os.Logger.Errorf("OnStepCurrentCartPlaceOrder UpdateDeliveryInfosAndBilling Error %v", err)
		return "", err
	}

	err = os.CartService.UpdatePurchaser(ctx, purchaser, nil)
	if err != nil {
		os.Logger.Errorf("OnStepCurrentCartPlaceOrder UpdatePurchaser Error %v", err)
		return "", err
	}

	//After setting DeliveryInfos - call SourcingEnginge (this will reload the cart and update all items!)
	err = os.SetSources(ctx)
	if err != nil {
		os.Logger.Errorf("OnStepCurrentCartPlaceOrder SetSources Error %v", err)
		return "", err
	}

	if payment == nil {
		payment = os.PaymentService.GetDefaultCartPayment(&decoratedCart.Cart)
	}

	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}

	orderid, orderError = os.PlaceOrder(ctx, decoratedCart, payment)

	if orderError != nil {
		os.Logger.WithField("category", "checkout.orderService").Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while placing the order. Please contact customer support.")
	}
	return orderid, nil
}
