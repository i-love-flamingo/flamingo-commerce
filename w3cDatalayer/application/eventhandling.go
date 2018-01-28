package application

import (
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	EventReceiver struct {
		Service              Service                   `inject:""`
		CartDecoratorFactory cart.DecoratedCartFactory `inject:""`
		Logger               flamingo.Logger           `inject:""`
	}
)

//Notify should get called by flamingo Eventlogic
// - OrderPlacedEvent is used to attach TransactionData - This is only useful in case where not directly redirected to a success page for example
func (e *EventReceiver) Notify(event event.Event) {
	switch event := event.(type) {
	//Handle OrderPlacedEvent and Set Transaction to current datalayer
	case *cart.OrderPlacedEvent:
		e.Logger.WithField("category", "w3cDatalayer").Debugf("Receive Event")
		if event.CurrentContext == nil {
			break
		}
		decoratedCart := e.CartDecoratorFactory.Create(event.CurrentContext, *event.Cart)
		if decoratedCart != nil {
			e.Service.CurrentContext = event.CurrentContext
			e.Service.SetTransaction(decoratedCart.Cart.GetCartTotals(), decoratedCart.DecoratedItems, event.OrderId)
			e.Service.AddEvent("orderplaced")
		}
	}
}
