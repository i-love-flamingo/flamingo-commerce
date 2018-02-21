package application

import (
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	EventReceiver struct {
		Service              *Service                   `inject:""`
		CartDecoratorFactory *cart.DecoratedCartFactory `inject:""`
		Logger               flamingo.Logger            `inject:""`
	}
)

//Notify should get called by flamingo Eventlogic
// - OrderPlacedEvent is used to attach TransactionData - This is only useful in case where not directly redirected to a success page for example
func (e *EventReceiver) Notify(event event.Event) {
	return

	switch currentEvent := event.(type) {
	//Handle OrderPlacedEvent and Set Transaction to current datalayer
	case *cart.OrderPlacedEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")
		if currentEvent.CurrentContext == nil {
			break
		}
		decoratedCart := e.CartDecoratorFactory.Create(currentEvent.CurrentContext, *currentEvent.Cart)
		if decoratedCart != nil {
			e.Service.CurrentContext = currentEvent.CurrentContext
			e.Service.SetTransaction(decoratedCart.Cart.GetCartTotals(), decoratedCart.DecoratedItems, currentEvent.OrderId)
			e.Service.AddEvent("orderplaced")
		}
	case *cart.AddToCartEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")
		if currentEvent.CurrentContext == nil {
			break
		}

		e.Service.CurrentContext = currentEvent.CurrentContext
		e.Service.CurrentContext.Session().AddFlash(
			currentEvent,
			"addToCart",
		)
	case *cart.ChangedQtyInCartEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")
		if currentEvent.CurrentContext == nil {
			break
		}

		e.Service.CurrentContext = currentEvent.CurrentContext
		e.Service.CurrentContext.Session().AddFlash(
			currentEvent,
			"changedQtyInCart",
		)
	}
}
