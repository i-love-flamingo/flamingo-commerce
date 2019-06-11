package cart

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"

	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"

	formDomain "flamingo.me/form/domain"

	"flamingo.me/form"

	"flamingo.me/flamingo-commerce/v3/cart/infrastructure/email"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/core/oauth"
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
		UseInMemoryCart bool `inject:"config:commerce.cart.useInMemoryCartServiceAdapters,optional"`
		EnableCartCache bool `inject:"config:commerce.cart.enableCartCache,optional"`
		UseEmailAdapter bool `inject:"config:commerce.cart.useEmailPlaceOrderAdapter,optional"`
	},
) {
	m.routerRegistry = routerRegistry
	if config != nil {
		m.useInMemoryCart = config.UseInMemoryCart
		m.enableCartCache = config.EnableCartCache
		m.useEmailAdapter = config.UseEmailAdapter
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
		injector.Bind((*placeorder.Service)(nil)).To(email.PlaceOrderServiceAdapter{})
	}
	//Register Default EventPublisher
	injector.Bind((*events.EventPublisher)(nil)).To(events.DefaultEventPublisher{})

	//Event
	flamingo.BindEventSubscriber(injector).To(application.EventReceiver{})

	// TemplateFunction
	flamingo.BindTemplateFunc(injector, "getCart", new(templatefunctions.GetCart))
	flamingo.BindTemplateFunc(injector, "getDecoratedCart", new(templatefunctions.GetDecoratedCart))

	injector.Bind((*cart.DeliveryInfoBuilder)(nil)).To(cart.DefaultDeliveryInfoBuilder{})

	if m.enableCartCache {
		injector.Bind((*application.CartCache)(nil)).To(application.CartSessionCache{})
	}

	//Register Form Data Provider
	injector.BindMap(new(formDomain.FormService), "commerce.cart.deliveryFormService").To(forms.DeliveryFormService{})
	injector.BindMap(new(formDomain.FormService), "commerce.cart.billingFormService").To(forms.BillingAddressFormService{})
	injector.BindMap(new(formDomain.FormService), "commerce.cart.personaldataFormService").To(forms.DefaultPersonalDataFormService{})

	web.BindRoutes(injector, new(routes))
}

// DefaultConfig enables inMemory cart service adapter
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"commerce": config.Map{
			"cart": config.Map{
				"useInMemoryCartServiceAdapters": true,
				"useEmailPlaceOrderAdapter":      true,
				"cacheLifetime":                  float64(1200), // in seconds
				"enableCartCache":                true,
			},
		},
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(oauth.Module),
		new(form.Module),
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

	registry.HandleDelete("cart.clean", r.viewController.CleanAndViewAction)
	registry.Route("/cart/delivery/:deliveryCode", "cart.clean")

	registry.HandleAny("cart.cleanDelivery", r.viewController.CleanDeliveryAndViewAction)
	registry.Route("/cart/delete/delivery/:deliveryCode", `cart.cleanDelivery(deliveryCode?="")`)

	registry.HandleAny("cart.deleteItem", r.viewController.DeleteAndViewAction)
	registry.Route("/cart/delete/:id", `cart.deleteItem(id,deliveryCode?="")`)
	r.apiRoutes(registry)
}

func (r *routes) apiRoutes(registry *web.RouterRegistry) {

	registry.Route("/api/cart", "cart.api.get")
	registry.HandleDelete("cart.api.get", r.apiController.DeleteCartAction)
	registry.HandleGet("cart.api.get", r.apiController.GetAction)

	//add command under the delivery:
	registry.Route("/api/cart/delivery/:deliveryCode/additem", `cart.api.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)
	registry.HandlePost("cart.api.add", r.apiController.AddAction)

	registry.Route("/api/cart/applyvoucher", `cart.api.applyVoucher(couponCode)`)
	registry.HandlePost("cart.api.applyVoucher", r.apiController.ApplyVoucherAndGetAction)
	registry.HandlePut("cart.api.applyVoucher", r.apiController.ApplyVoucherAndGetAction)

	registry.Route("/api/cart/removevoucher", `cart.api.removeVoucher(couponCode)`)
	registry.HandlePost("cart.api.removeVoucher", r.apiController.RemoveVoucherAndGetAction)
	registry.HandleDelete("cart.api.removeVoucher", r.apiController.RemoveVoucherAndGetAction)

	registry.Route("/api/cart/billing", `cart.api.billing`)
	registry.HandlePost("cart.api.billing", r.apiController.BillingAction)

	registry.Route("/api/cart/delivery/:deliveryCode", `cart.api.delivery.delete`)
	registry.HandleDelete("cart.api.delivery.delete", r.apiController.DeleteDelivery)

	registry.Route("/api/cart/delivery/:deliveryCode/deliveryinfo", `cart.api.delivery.update`)
	registry.HandlePost("cart.api.delivery.update", r.apiController.UpdateDeliveryInfoAction)

	//registry.Route("/api/cart/delivery/:shipping", `cart.api.shipping(deliveryCode?="")`)
	//TODO registry.HandleDelete("cart.api.delivery", r.apiController.DeleteDelivery)
}
