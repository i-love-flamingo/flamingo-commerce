package application

import (
	"errors"
	"log"

	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (

	// PaymentService
	OrderService struct {
		SourcingEngine SourcingEngine  `inject:""`
		PaymentService PaymentService  `inject:""`
		Logger         flamingo.Logger `inject:""`
	}
)

func (os OrderService) PlaceOrder(ctx web.Context, decoratedCart cart.DecoratedCart, shippingMethod string, shippingCarrierCode string, billingAddress *cart.Address, shippingAddress *cart.Address) (orderid string, orderError error) {
	err := os.SourcingEngine.SetSourcesForCartItems(ctx, decoratedCart)
	if err != nil {
		log.Printf("Error while getting pickup sources: %v", err)
		return "", errors.New("Error while getting pickup sources.")
	}

	err = decoratedCart.Cart.SetShippingInformation(ctx, shippingAddress, billingAddress, shippingCarrierCode, shippingMethod)
	if err != nil {
		os.Logger.Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while setting shipping informations.")
	}

	orderid, orderError = decoratedCart.Cart.PlaceOrder(ctx, os.PaymentService.GetPayment())

	if orderError != nil {
		os.Logger.Errorf("Error during place Order: %v", err)
		return "", errors.New("Error while placing the order.")
	}
	return orderid, nil

}
