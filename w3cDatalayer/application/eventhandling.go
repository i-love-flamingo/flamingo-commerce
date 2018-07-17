package application

import (
	"context"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/w3cDatalayer/domain"
	authDomain "flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (
	EventReceiver struct {
		factory              *Factory
		cartDecoratorFactory *cart.DecoratedCartFactory
		logger               flamingo.Logger
	}
)

func (e *EventReceiver) Inject(factory *Factory, cartFactory *cart.DecoratedCartFactory, logger flamingo.Logger) {
	e.factory = factory
	e.cartDecoratorFactory = cartFactory
	e.logger = logger
}

//NotifyWithContext should get called by flamingo Eventlogic.
// We use it to listen to Events that are relevant for the Datalayer
// In case the events might be asycron (e.g. the origin action does a redirect to a sucess page) - we save the datalayer Event to a Session Flash - to make sure it is still available the first time the DatalayerService.Get is calles
func (e *EventReceiver) NotifyWithContext(ctx context.Context, event event.Event) {
	switch currentEvent := event.(type) {
	//Handle OrderPlacedEvent and Set Transaction to current datalayer
	case *application.AddToCartEvent:
		e.logger.WithField("category", "w3cDatalayer").Debug("Receive Event AddToCartEvent")
		if webContext, ok := ctx.(web.Context); ok {
			// In case of Configurable: the MarketplaceCode which is interesting for the datalayer is the Variant that is selected
			saleableProductCode := currentEvent.MarketplaceCode
			if currentEvent.VariantMarketplaceCode != "" {
				saleableProductCode = currentEvent.VariantMarketplaceCode
			}
			dataLayerEvent := e.factory.BuildAddToBagEvent(saleableProductCode, currentEvent.ProductName, currentEvent.Qty)
			webContext.Session().AddFlash(
				dataLayerEvent,
				SESSION_EVENTS_KEY,
			)
		}
	case *application.ChangedQtyInCartEvent:
		e.logger.WithField("category", "w3cDatalayer").Debug("Receive Event ChangedQtyInCartEvent")

		if webContext, ok := ctx.(web.Context); ok {
			saleableProductCode := currentEvent.MarketplaceCode
			if currentEvent.VariantMarketplaceCode != "" {
				saleableProductCode = currentEvent.VariantMarketplaceCode
			}
			dataLayerEvent := e.factory.BuildChangeQtyEvent(saleableProductCode, currentEvent.ProductName, currentEvent.QtyAfter, currentEvent.QtyBefore, currentEvent.CartId)
			webContext.Session().AddFlash(
				dataLayerEvent,
				SESSION_EVENTS_KEY,
			)

		}
	case *authDomain.LoginEvent:
		e.logger.WithField("category", "w3cDatalayer").Debug("Receive Event LoginEvent")
		if webContext, ok := ctx.(web.Context); ok {

			dataLayerEvent := domain.Event{EventInfo: make(map[string]interface{})}
			dataLayerEvent.EventInfo["eventName"] = "login"
			webContext.Session().AddFlash(
				dataLayerEvent,
				SESSION_EVENTS_KEY,
			)

		}
	}

}
