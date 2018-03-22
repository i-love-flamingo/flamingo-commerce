package application

import (
	"go.aoe.com/flamingo/core/auth/domain"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (

	//EventReceiver - handles events from other packages
	EventReceiver struct {
		Logger              flamingo.Logger      `inject:""`
		CartService         *CartService         `inject:""`
		CartReceiverService *CartReceiverService `inject:""`
	}
)

//Notify should get called by flamingo Eventlogic
func (e *EventReceiver) Notify(event event.Event) {
	switch eventType := event.(type) {
	//Handle LoginEvent and Merge Cart
	case *domain.LoginEvent:
		if eventType == nil {
			return
		}
		if eventType.Context == nil {
			return
		}
		if !e.CartReceiverService.ShouldHaveGuestCart(eventType.Context) {
			return
		}
		guestCart, err := e.CartReceiverService.ViewGuestCart(eventType.Context)
		if err != nil {
			e.Logger.WithField("category", "cart").Errorf("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.CartReceiverService.UserService.IsLoggedIn(eventType.Context) {
			e.Logger.WithField("category", "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}

		for _, item := range guestCart.Cartitems {
			e.Logger.WithField("category", "cart").Debugf("Merging item from guest to user cart %v", item)

			e.CartService.AddProduct(eventType.Context, cartDomain.AddRequest{MarketplaceCode: item.MarketplaceCode, Qty: item.Qty, VariantMarketplaceCode: item.VariantMarketPlaceCode})
		}
		e.CartService.DeleteSessionGuestCart(eventType.Context)

	}
}
