package email

import (
	"context"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	authDomain "flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	//TODO - need to be implemented
	EMailAdapter struct {
		emailAddress string
		logger flamingo.Logger
	}
)

var (
	_ cartDomain.PlaceOrderService = new(EMailAdapter)
)


func (e *EMailAdapter) Inject(logger flamingo.Logger, config *struct {
	EmailAddress string `inject:"config:cart.emailAdapter.emailAddress"`
})  {
	e.emailAddress = config.EmailAddress
	e.logger = logger
}

func (e *EMailAdapter) PlaceGuestCart(ctx context.Context, cart *cartDomain.Cart, payment *cartDomain.CartPayment) (cartDomain.PlacedOrderInfos, error) {
	var placedOrders cartDomain.PlacedOrderInfos
	placedOrders = append(placedOrders, cartDomain.PlacedOrderInfo{
		OrderNumber: "1",
	})
	e.getLogger().Warn("send mail not implemented")
	return placedOrders, nil

}
func (e *EMailAdapter) PlaceCustomerCart(ctx context.Context, auth authDomain.Auth, cart *cartDomain.Cart, payment *cartDomain.CartPayment) (cartDomain.PlacedOrderInfos, error) {
	var placedOrders cartDomain.PlacedOrderInfos
	placedOrders = append(placedOrders, cartDomain.PlacedOrderInfo{
		OrderNumber: "1",
	})
	e.getLogger().Warn("send mail not implemented")
	return placedOrders, nil
}


func (e *EMailAdapter) getLogger() flamingo.Logger {
	return e.logger.WithField("module","cart").WithField("category","emailAdapter")
}