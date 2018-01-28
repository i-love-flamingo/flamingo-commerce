package application

import (
	"context"

	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	//DomainEventPublisher implements the event publisher of the domain and uses the framework event router
	DomainEventPublisher struct {
		Logger flamingo.Logger `inject:""`
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
