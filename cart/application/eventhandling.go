package application

import (
	"context"

	"go.aoe.com/flamingo/core/auth/domain"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	//DomainEventPublisher implements the event publisher of the domain and uses the framework event router
	DomainEventPublisher struct {
		Logger flamingo.Logger `inject:""`
	}

	//EventReceiver
	EventReceiver struct {
		Logger      flamingo.Logger `inject:""`
		CartService *CartService    `inject:""`
	}
)

func (d *DomainEventPublisher) PublishOrderPlacedEvent(ctx context.Context, carto *cart.Cart, orderId string) {
	event := cart.OrderPlacedEvent{
		Cart:    carto,
		OrderId: orderId,
	}
	if webContext, ok := ctx.(web.Context); ok {
		d.Logger.Infof("Publish Event OrderPlacedEvent for Order: %v", orderId)
		event.CurrentContext = webContext
		//For now we publish only to Flamingo default Event Router
		webContext.EventRouter().Dispatch(event)
	}

}

//Notify should get called by flamingo Eventlogic
func (e *EventReceiver) Notify(event event.Event) {
	switch eventType := event.(type) {
	//Handle OrderPlacedEvent and Set Transaction to current datalayer
	case *domain.LoginEvent:
		if eventType == nil {
			return
		}
		if eventType.Context == nil {
			return
		}
		if !e.CartService.HasSessionAGuestCart(eventType.Context) {
			return
		}
		guestCart, err := e.CartService.GetSessionGuestCart(eventType.Context)
		if err != nil {
			e.Logger.WithField("category", "cart").Errorf("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.CartService.UserService.IsLoggedIn(eventType.Context) {
			e.Logger.WithField("category", "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}

		for _, item := range guestCart.Cartitems {
			e.Logger.WithField("category", "cart").Debugf("Merging item from guest to user cart %v", item)

			e.CartService.AddProduct(eventType.Context, cart.AddRequest{MarketplaceCode: item.MarketplaceCode, Qty: item.Qty, VariantMarketplaceCode: item.VariantMarketPlaceCode})
		}
		e.CartService.DeleteSessionGuestCart(eventType.Context)

	}
}
