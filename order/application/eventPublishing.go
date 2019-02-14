package application

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/order/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// OrderPlacedEvent defines event properties
	OrderPlacedEvent struct {
		Cart             *cart.Cart
		PlacedOrderInfos domain.PlacedOrderInfos
	}

	//EventPublisher - technology free interface to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *cart.Cart, placedOrderInfos domain.PlacedOrderInfos)
	}

	//DefaultEventPublisher implements the event publisher of the domain and uses the framework event router
	DefaultEventPublisher struct {
		logger      flamingo.Logger
		eventRouter event.Router
	}
)

var (
	_ flamingo.Event    = (*OrderPlacedEvent)(nil)
	_ EventPublisher = (*DefaultEventPublisher)(nil)
)

// Inject the default event publisher dependencies
func (dep *DefaultEventPublisher) Inject(
	Logger flamingo.Logger,
	EventRouter event.Router,
) {
	dep.logger = Logger
	dep.eventRouter = EventRouter
}

// PublishOrderPlacedEvent publishes an event for placed orders
func (dep *DefaultEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cart.Cart, placedOrderInfos domain.PlacedOrderInfos) {
	eventObject := OrderPlacedEvent{
		Cart:             cart,
		PlacedOrderInfos: placedOrderInfos,
	}

	dep.logger.Info("Publish Event OrderPlacedEvent for Order: %#v", placedOrderInfos)

	//For now we publish only to Flamingo default Event Router
	dep.eventRouter.Dispatch(ctx, &eventObject)
}
