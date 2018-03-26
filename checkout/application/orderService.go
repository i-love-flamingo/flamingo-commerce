package application

import (
	"errors"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (

	// PaymentService
	OrderService struct {
		SourcingEngine      *SourcingEngine                  `inject:""`
		PaymentService      *PaymentService                  `inject:""`
		Logger              flamingo.Logger                  `inject:""`
		CartService         *application.CartService         `inject:""`
		CartReceiverService *application.CartReceiverService `inject:""`
	}
)

func (os OrderService) PlaceOrder(ctx web.Context, shippingMethod string, shippingCarrierCode string, billingAddress *cart.Address, shippingAddress *cart.Address) (orderid string, orderError error) {

	decoratedCart, behaviour, err := os.CartReceiverService.GetDecoratedCart(ctx)
	if err != nil {
		return "", err
	}
	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}

	err = os.SourcingEngine.SetSourcesForCartItems(ctx, decoratedCart, behaviour)
	if err != nil {
		os.Logger.WithField("category", "checkout.orderService").Errorf("Error while getting pickup sources: %v", err)
		return "", errors.New("Error while getting pickup sources.")
	}

	err = os.CartService.SetShippingInformation(ctx, shippingAddress, billingAddress, shippingCarrierCode, shippingMethod)
	if err != nil {
		os.Logger.WithField("category", "checkout.orderService").Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while setting shipping informations.")
	}

	orderid, orderError = os.CartService.PlaceOrder(ctx, os.PaymentService.GetPayment())

	if orderError != nil {
		os.Logger.WithField("category", "checkout.orderService").Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while placing the order. Please contact customer support.")
	}
	return orderid, nil
}
