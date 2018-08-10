package cart

import (
	"encoding/gob"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/cart/infrastructure"
	"flamingo.me/flamingo-commerce/cart/interfaces/controller"
	"flamingo.me/flamingo-commerce/cart/interfaces/templatefunctions"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/template"
)

type (
	// CartModule registers our profiler
	CartModule struct {
		RouterRegistry  *router.Registry `inject:""`
		UseInMemoryCart bool             `inject:"config:cart.useInMemoryCartServiceAdapters"`
		EnableCartCache bool             `inject:"config:cart.enableCartCache,optional"`
	}
)

// Configure module
func (m *CartModule) Configure(injector *dingo.Injector) {
	if m.UseInMemoryCart {
		injector.Bind((*infrastructure.CartStorage)(nil)).To(infrastructure.InMemoryCartStorage{}).AsEagerSingleton()
		injector.Bind((*cart.GuestCartService)(nil)).To(infrastructure.InMemoryGuestCartService{})
		injector.Bind((*cart.CustomerCartService)(nil)).To(infrastructure.InMemoryCustomerCartService{})
	}
	//Register Default EventPublisher
	injector.Bind((*application.EventPublisher)(nil)).To(application.DefaultEventPublisher{})

	//Event
	injector.BindMulti((*event.Subscriber)(nil)).To(application.EventReceiver{})

	// TemplateFunction
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.GetCart{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.GetDecoratedCart{})

	injector.Bind((*cart.DeliveryInfoBuilder)(nil)).To(cart.DefaultDeliveryInfoBuilder{})

	if m.EnableCartCache {
		injector.Bind((*application.CartCache)(nil)).To(application.CartSessionCache{})
	}

	router.Bind(injector, new(routes))
}

// DefaultConfig enables inMemory cart service adapter
func (m *CartModule) DefaultConfig() config.Map {
	return config.Map{
		"cart": config.Map{
			"useInMemoryCartServiceAdapters": true,
		},
	}
}

type routes struct {
	viewController *controller.CartViewController
	apiController  *controller.CartApiController
}

func (r *routes) Inject(viewController *controller.CartViewController, apiController *controller.CartApiController) {
	r.viewController = viewController
	r.apiController = apiController
}

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleAny("cart.view", r.viewController.ViewAction)
	registry.Route("/cart", "cart.view")

	registry.HandleAny("cart.add", r.viewController.AddAndViewAction)
	registry.Route("/cart/add/:marketplaceCode", `cart.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)

	registry.HandleAny("cart.updateQty", r.viewController.UpdateQtyAndViewAction)
	registry.Route("/cart/update/:id", `cart.updateQty(id,qty?="1")`)

	registry.HandleAny("cart.deleteAllItems", r.viewController.DeleteAllAndViewAction)
	registry.Route("/cart/delete/all", `cart.deleteAllItems`)

	registry.HandleAny("cart.deleteItem", r.viewController.DeleteAndViewAction)
	registry.Route("/cart/delete/:id", `cart.deleteItem(id,deliveryCode?="")`)

	gob.Register(cart.Cart{})

	//DecoratedCart API:

	registry.HandleAny("cart.api.get", r.apiController.GetAction)
	registry.HandleAny("cart.api.add", r.apiController.AddAction)
	registry.HandleAny("cart.api.applyVoucher", r.apiController.ApplyVoucherAndGetAction)

	registry.Route("/api/cart", "cart.api.get")
	registry.Route("/api/cart/add/:marketplaceCode", `cart.api.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)
	registry.Route("/api/cart/applyvoucher/:couponCode", `cart.api.applyVoucher(couponCode)`)
}
