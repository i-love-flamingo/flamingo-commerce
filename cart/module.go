package cart

import (
	"encoding/gob"
	"flamingo/core/cart/domain/cart"
	"flamingo/core/cart/infrastructure"
	controller "flamingo/core/cart/interfaces/controller"
	"flamingo/framework/config"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		Config         config.Map       `inject:"config:cart"`
		//EventRouter    event.Router     `inject:""`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if v, ok := m.Config["useInMemoryCartServiceAdapters"].(bool); v && ok {
		injector.Bind((*cart.GuestCartService)(nil)).In(dingo.Singleton).To(infrastructure.InMemoryCartServiceAdapter{})
	}

	m.RouterRegistry.Handle("cart.view", new(controller.CartViewController))
	m.RouterRegistry.Route("/cart", "cart.view")

	gob.Register(cart.Cart{})

	//DecoratedCart API:

	m.RouterRegistry.Handle("cart.api.get", (*controller.CartApiController).GetAction)
	m.RouterRegistry.Handle("cart.api.add", (*controller.CartApiController).AddAction)

	m.RouterRegistry.Route("/api/cart", "cart.api.get")
	m.RouterRegistry.Route("/api/cart/add", "cart.api.add(marketplaceCode)")

	//
	//m.RouterRegistry.Handle("cart.api", new(controller.CartApiController))
	//m.RouterRegistry.Route("/api/cart", "cart.view")
	//m.RouterRegistry.Handle("/api/cart/item/add", new(controller.CartItemAddApiController))
	//m.RouterRegistry.Route("/api/cart", "cart.item.add.api")
	//
	////	m.RouterRegistry.Handle("foo", controller.CartItemAddApiController.AddToBasketAction)
	//
	//m.RouterRegistry.Handle("logintest", new(controller.TestLoginController))
	//m.RouterRegistry.Route("/logintest", "logintest")
	//
	//m.EventRouter.AddSubscriber(new(application.EventOrchestration))
	//
	////	a := controller.CartItemAddApiController.AddToBasketAction
	////	a

	//m.RouterRegistry.Mount("/api/cart", new(controller.CartApiController))
	//m.RouterRegistry.Mount("/cart", new(controller.CartController))
}
