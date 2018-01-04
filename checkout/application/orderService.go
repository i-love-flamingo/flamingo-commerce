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
		SourcingEngine SourcingEngine          `inject:""`
		PaymentService PaymentService          `inject:""`
		Logger         flamingo.Logger         `inject:""`
		CartService    application.CartService `inject:""`
	}
)

func (os OrderService) PlaceOrder(ctx web.Context, decoratedCart cart.DecoratedCart, shippingMethod string, shippingCarrierCode string, billingAddress *cart.Address, shippingAddress *cart.Address) (orderid string, orderError error) {
	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}

	err := os.SourcingEngine.SetSourcesForCartItems(ctx, decoratedCart)
	if err != nil {
		os.Logger.Errorf("Error while getting pickup sources: %v", err)
		return "", errors.New("Error while getting pickup sources.")
	}

	err = decoratedCart.Cart.SetShippingInformation(ctx, os.CartService.Auth(ctx), shippingAddress, billingAddress, shippingCarrierCode, shippingMethod)
	if err != nil {
		os.Logger.Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while setting shipping informations.")
	}

	orderid, orderError = decoratedCart.Cart.PlaceOrder(ctx, os.CartService.Auth(ctx), os.PaymentService.GetPayment())

	if orderError != nil {
		os.Logger.Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while placing the order.")
	}
	return orderid, nil

}
