package cart

import (
	"encoding/gob"
	"flamingo/core/cart/domain"
	"flamingo/core/cart/infrastructure"
	controller "flamingo/core/cart/interfaces/controller"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		//EventRouter    event.Router     `inject:""`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.CartService)(nil)).In(dingo.Singleton).To(infrastructure.InMemoryCartService{})

	m.RouterRegistry.Handle("cart.view", new(controller.CartViewController))
	m.RouterRegistry.Route("/cart", "cart.view")

	gob.Register(domain.Cart{})

	//Cart API:

	m.RouterRegistry.Handle("cart.api.get", new(controller.CartApiGetController))
	m.RouterRegistry.Handle("cart.api.add", new(controller.CartApiAddController))

	m.RouterRegistry.Route("/api/cart", "cart.api.get")
	m.RouterRegistry.Route("/api/cart/add", "cart.api.add(productcode)")

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
