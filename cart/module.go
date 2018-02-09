package cart

import (
	"encoding/gob"

	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/core/cart/infrastructure"
	"go.aoe.com/flamingo/core/cart/interfaces/controller"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// CartModule registers our profiler
	CartModule struct {
		RouterRegistry  *router.Registry `inject:""`
		UseInMemoryCart bool             `inject:"config:cart.useInMemoryCartServiceAdapters"`
		//EventRouter    event.Router     `inject:""`
	}
)

// Configure module
func (m *CartModule) Configure(injector *dingo.Injector) {
	if m.UseInMemoryCart {
		injector.Bind((*infrastructure.GuestCartStorage)(nil)).To(infrastructure.InmemoryGuestCartStorage{}).AsEagerSingleton()
		injector.Bind((*infrastructure.CustomerCartStorage)(nil)).To(infrastructure.InmemoryCustomerCartStorage{}).AsEagerSingleton()
		injector.Bind((*cart.GuestCartService)(nil)).To(infrastructure.GuestCartServiceAdapter{})
		injector.Bind((*cart.CustomerCartService)(nil)).To(infrastructure.CustomerCartServiceAdapter{})
	}
	//Register DomainEventPublisher
	injector.Bind((*cart.EventPublisher)(nil)).To(application.DomainEventPublisher{})

	m.RouterRegistry.Handle("cart.view", (*controller.CartViewController).ViewAction)
	m.RouterRegistry.Route("/cart", "cart.view")

	m.RouterRegistry.Handle("cart.add", (*controller.CartViewController).AddAndViewAction)
	m.RouterRegistry.Route("/cart/add/:marketplaceCode", `cart.add(marketplaceCode,variantMarketplaceCode?="",qty?="1")`)

	m.RouterRegistry.Handle("cart.updateQty", (*controller.CartViewController).UpdateQtyAndViewAction)
	m.RouterRegistry.Route("/cart/update/:id", `cart.updateQty(id,qty?="1")`)

	m.RouterRegistry.Handle("cart.deleteItem", (*controller.CartViewController).DeleteAndViewAction)
	m.RouterRegistry.Route("/cart/delete/:id", `cart.deleteItem(id)`)

	gob.Register(cart.Cart{})

	//DecoratedCart API:

	m.RouterRegistry.Handle("cart.api.get", (*controller.CartApiController).GetAction)
	m.RouterRegistry.Handle("cart.api.add", (*controller.CartApiController).AddAction)

	m.RouterRegistry.Route("/api/cart", "cart.api.get")
	m.RouterRegistry.Route("/api/cart/add/:marketplaceCode", `cart.api.add(marketplaceCode,variantMarketplaceCode?="",qty?="1")`)

	//Event
	injector.BindMulti((*event.Subscriber)(nil)).To(application.EventReceiver{})

}

// DefaultConfig enables inMemory cart service adapter
func (m *CartModule) DefaultConfig() config.Map {
	return config.Map{
		"cart": config.Map{
			"useInMemoryCartServiceAdapters": true,
		},
	}
}
