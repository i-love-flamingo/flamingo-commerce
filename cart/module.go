package cart

import (
	"flamingo/core/cart/application"
	"flamingo/core/cart/interfaces/controller"
	"flamingo/framework/dingo"
	"flamingo/framework/event"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
		EventRouter    event.Router           `inject:""`
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("cart.view", new(controller.CartViewController))
	m.RouterRegistry.Route("/cart", "cart.view")

	m.RouterRegistry.Handle("cart.api", new(controller.CartApiController))
	m.RouterRegistry.Route("/api/cart", "cart.view")
	m.RouterRegistry.Handle("/api/cart/item/add", new(controller.CartItemAddApiController))
	m.RouterRegistry.Route("/api/cart", "cart.item.add.api")

	m.RouterRegistry.Handle("logintest", new(controller.TestLoginController))
	m.RouterRegistry.Route("/logintest", "logintest")

	m.EventRouter.AddSubscriber(new(application.EventOrchestration))
}
