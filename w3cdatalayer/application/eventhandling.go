package application

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo/v3/core/auth"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// EventReceiver struct with required dependencies
	EventReceiver struct {
		factory              *Factory
		cartDecoratorFactory *decorator.DecoratedCartFactory
		logger               flamingo.Logger
	}
)

// Inject method
func (e *EventReceiver) Inject(factory *Factory, cartFactory *decorator.DecoratedCartFactory, logger flamingo.Logger) {
	e.factory = factory
	e.cartDecoratorFactory = cartFactory
	e.logger = logger.WithField("category", "w3cDatalayer").WithField(flamingo.LogKeyModule, "w3cdatalayer")
}

// Notify should get called by flamingo Eventlogic.
// We use it to listen to Events that are relevant for the Datalayer
// In case the events might be asycron (e.g. the origin action does a redirect to a success page) - we save the datalayer Event to a Session Flash - to make sure it is still available the first time the DatalayerService.Get is calles
func (e *EventReceiver) Notify(ctx context.Context, event flamingo.Event) {
	switch currentEvent := event.(type) {
	// Handle OrderPlacedEvent and Set Transaction to current datalayer
	case *events.AddToCartEvent:
		e.logger.WithContext(ctx).Debug("Receive Event AddToCartEvent")
		session := web.SessionFromContext(ctx)
		if session != nil {
			// In case of Configurable: the MarketplaceCode which is interesting for the datalayer is the Variant that is selected
			saleableProductCode := currentEvent.MarketplaceCode
			if currentEvent.VariantMarketplaceCode != "" {
				saleableProductCode = currentEvent.VariantMarketplaceCode
			}
			dataLayerEvent := e.factory.BuildAddToBagEvent(saleableProductCode, currentEvent.ProductName, currentEvent.Qty)
			session.AddFlash(
				dataLayerEvent,
				SessionEventsKey,
			)
		}
	case *events.ChangedQtyInCartEvent:
		e.logger.WithContext(ctx).WithField("category", "w3cDatalayer").Debug("Receive Event ChangedQtyInCartEvent")

		session := web.SessionFromContext(ctx)
		if session != nil {
			saleableProductCode := currentEvent.MarketplaceCode
			if currentEvent.VariantMarketplaceCode != "" {
				saleableProductCode = currentEvent.VariantMarketplaceCode
			}
			dataLayerEvent := e.factory.BuildChangeQtyEvent(saleableProductCode, currentEvent.ProductName, currentEvent.QtyAfter, currentEvent.QtyBefore, currentEvent.Cart.ID)
			session.AddFlash(
				dataLayerEvent,
				SessionEventsKey,
			)

		}
	case *auth.WebLoginEvent:
		e.logger.WithContext(ctx).WithField("category", "w3cDatalayer").Debug("Receive Event WebLoginEvent")
		session := web.SessionFromContext(ctx)
		if session != nil {

			dataLayerEvent := domain.Event{EventInfo: make(map[string]interface{})}
			dataLayerEvent.EventInfo["eventName"] = "login"
			session.AddFlash(
				dataLayerEvent,
				SessionEventsKey,
			)
		}
	}
}
