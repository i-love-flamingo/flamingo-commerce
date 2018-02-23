package application

import (
	"context"

	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	EventReceiver struct {
		Service              *Service                   `inject:""`
		CartDecoratorFactory *cart.DecoratedCartFactory `inject:""`
		Logger               flamingo.Logger            `inject:""`
	}
)

//NotifyWithContext should get called by flamingo Eventlogic
// - OrderPlacedEvent is used to attach TransactionData - This is only useful in case where not directly redirected to a success page for example
func (e *EventReceiver) NotifyWithContext(ctx context.Context, event event.Event) {
	switch currentEvent := event.(type) {
	//Handle OrderPlacedEvent and Set Transaction to current datalayer
	case *cart.OrderPlacedEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")

		decoratedCart := e.CartDecoratorFactory.Create(ctx, *currentEvent.Cart)
		if decoratedCart != nil {
			if webContext, ok := ctx.(web.Context); ok {
				e.Service.CurrentContext = webContext
				e.Service.SetTransaction(decoratedCart.Cart.GetCartTotals(), decoratedCart.DecoratedItems, currentEvent.OrderId)
				e.Service.AddEvent("orderplaced")
			}
		}
	case *cart.AddToCartEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")

		if webContext, ok := ctx.(web.Context); ok {
			webContext.Session().AddFlash(
				event,
				"addToCart",
			)
		}
	case *cart.ChangedQtyInCartEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")

		if webContext, ok := ctx.(web.Context); ok {
			webContext.Session().AddFlash(
				event,
				"changedQtyInCart",
			)
		}
	}
}
