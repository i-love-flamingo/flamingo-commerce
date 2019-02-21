package cart

import (
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/cart/infrastructure/email"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module registers our profiler
	Module struct {
		routerRegistry  *web.RouterRegistry
		useInMemoryCart bool
		useEmailAdapter bool
		enableCartCache bool
	}
)

// Inject dependencies
func (m *Module) Inject(
	routerRegistry *web.RouterRegistry,
	config *struct {
		UseInMemoryCart bool `inject:"config:cart.useInMemoryCartServiceAdapters"`
		EnableCartCache bool `inject:"config:cart.enableCartCache,optional"`
		UseEmailAdapter bool `inject:"config:cart.useEmailPlaceOrderAdapter,optional"`
	},
) {
	m.routerRegistry = routerRegistry
	if config != nil {
		m.useInMemoryCart = config.UseInMemoryCart
		m.enableCartCache = config.EnableCartCache
	}
}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.useInMemoryCart {
		injector.Bind((*infrastructure.CartStorage)(nil)).To(infrastructure.InMemoryCartStorage{}).AsEagerSingleton()
		injector.Bind((*cart.GuestCartService)(nil)).To(infrastructure.InMemoryGuestCartService{})
		injector.Bind((*cart.CustomerCartService)(nil)).To(infrastructure.InMemoryCustomerCartService{})
	}
	if m.useEmailAdapter {
		injector.Bind((*cart.PlaceOrderService)(nil)).To(email.PlaceOrderServiceAdapter{})
	}
	//Register Default EventPublisher
	injector.Bind((*application.EventPublisher)(nil)).To(application.DefaultEventPublisher{})

	//Event
	flamingo.BindEventSubscriber(injector).To(application.EventReceiver{})

	// TemplateFunction
	flamingo.BindTemplateFunc(injector, "getCart", new(templatefunctions.GetCart))
	flamingo.BindTemplateFunc(injector, "getDecoratedCart", new(templatefunctions.GetDecoratedCart))

	injector.Bind((*cart.DeliveryInfoBuilder)(nil)).To(cart.DefaultDeliveryInfoBuilder{})

	if m.enableCartCache {
		injector.Bind((*application.CartCache)(nil)).To(application.CartSessionCache{})
	}

	web.BindRoutes(injector, new(routes))
}

// DefaultConfig enables inMemory cart service adapter
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cart": config.Map{
			"useInMemoryCartServiceAdapters": true,
			"useEmailPlaceOrderAdapter":      true,
			"cacheLifetime":                  float64(1200), // in seconds
			"enableCartCache":                false,
		},
	}
}

type routes struct {
	viewController *controller.CartViewController
	apiController  *controller.CartAPIController
}

func (r *routes) Inject(viewController *controller.CartViewController, apiController *controller.CartAPIController) {
	r.viewController = viewController
	r.apiController = apiController
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleAny("cart.view", r.viewController.ViewAction)
	registry.Route("/cart", "cart.view")

	registry.HandleAny("cart.add", r.viewController.AddAndViewAction)
	registry.Route("/cart/add/:marketplaceCode", `cart.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)

	registry.HandleAny("cart.updateQty", r.viewController.UpdateQtyAndViewAction)
	registry.Route("/cart/update/:id", `cart.updateQty(id,qty?="1",deliveryCode?="")`)

	registry.HandleAny("cart.deleteAllItems", r.viewController.DeleteAllAndViewAction)
	registry.Route("/cart/delete/all", "cart.deleteAllItems")

	registry.HandleAny("cart.clean", r.viewController.CleanAndViewAction)
	registry.Route("/cart/clean", "cart.clean")

	registry.HandleAny("cart.cleanDelivery", r.viewController.CleanDeliveryAndViewAction)
	registry.Route("/cart/delete/delivery/:deliveryCode", `cart.cleanDelivery(deliveryCode?="")`)

	registry.HandleAny("cart.deleteItem", r.viewController.DeleteAndViewAction)
	registry.Route("/cart/delete/:id", `cart.deleteItem(id,deliveryCode?="")`)
	gob.Register(cart.Cart{})

	// DecoratedCart API:

	registry.HandleGet("cart.api.get", r.apiController.GetAction)
	registry.HandleDelete("cart.api.get", r.apiController.CleanAndGetAction)
	registry.HandleDelete("cart.api.delivery", r.apiController.CleanDeliveryAndGetAction)
	registry.HandleAny("cart.api.add", r.apiController.AddAction)
	registry.HandleAny("cart.api.applyVoucher", r.apiController.ApplyVoucherAndGetAction)

	registry.Route("/api/cart", "cart.api.get")
	registry.Route("/api/cart/delivery/:deliveryCode", `cart.api.get(deliveryCode?="")`)
	registry.Route("/api/cart/add/:marketplaceCode", `cart.api.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)
	registry.Route("/api/cart/applyvoucher/:couponCode", `cart.api.applyVoucher(couponCode)`)
}
