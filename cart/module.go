package cart

import (
	"encoding/gob"

	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/cart/infrastructure"
	"go.aoe.com/flamingo/core/cart/interfaces/controller"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry  *router.Registry `inject:""`
		UseInMemoryCart bool             `inject:"config:cart.useInMemoryCartServiceAdapters"`
		//EventRouter    event.Router     `inject:""`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.UseInMemoryCart {
		injector.Bind((*cart.GuestCartService)(nil)).In(dingo.Singleton).To(infrastructure.InMemoryCartServiceAdapter{})
	}

	m.RouterRegistry.Handle("cart.view", (*controller.CartViewController).ViewAction)
	m.RouterRegistry.Route("/cart", "cart.view")

	m.RouterRegistry.Handle("cart.add", (*controller.CartViewController).AddAndViewAction)
	m.RouterRegistry.Route("/cart/add/:marketplaceCode", `cart.add(marketplaceCode,variantMarketplaceCode?="",qty?="1")`)

	gob.Register(cart.Cart{})

	//DecoratedCart API:

	m.RouterRegistry.Handle("cart.api.get", (*controller.CartApiController).GetAction)
	m.RouterRegistry.Handle("cart.api.add", (*controller.CartApiController).AddAction)

	m.RouterRegistry.Route("/api/cart", "cart.api.get")
	m.RouterRegistry.Route("/api/cart/add/:marketplaceCode", `cart.api.add(marketplaceCode,variantMarketplaceCode?="",qty?="1")`)

}

// DefaultConfig enables inMemory cart service adapter
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cart": config.Map{
			"useInMemoryCartServiceAdapters": true,
		},
	}
}
