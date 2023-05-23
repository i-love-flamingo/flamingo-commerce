package cart

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/form"
	formDomain "flamingo.me/form/domain"
	flamingographql "flamingo.me/graphql"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	placeorderAdapter "flamingo.me/flamingo-commerce/v3/cart/infrastructure/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/templatefunctions"
	"flamingo.me/flamingo-commerce/v3/customer"
	"flamingo.me/flamingo-commerce/v3/product"
)

type (
	// Module registers our profiler
	Module struct {
		routerRegistry                *web.RouterRegistry
		enableDefaultCartAdapter      bool
		enablePlaceOrderLoggerAdapter bool
		enableCartCache               bool
		cartMergeStrategy             string
	}
)

// Inject dependencies
func (m *Module) Inject(
	routerRegistry *web.RouterRegistry,
	config *struct {
		EnableDefaultCartAdapter      bool   `inject:"config:commerce.cart.defaultCartAdapter.enabled,optional"`
		EnableCartCache               bool   `inject:"config:commerce.cart.enableCartCache,optional"`
		CartMergeStrategy             string `inject:"config:commerce.cart.mergeStrategy,optional"`
		EnablePlaceOrderLoggerAdapter bool   `inject:"config:commerce.cart.placeOrderLogger.enabled,optional"`
	},
) {
	m.routerRegistry = routerRegistry
	if config != nil {
		m.enableDefaultCartAdapter = config.EnableDefaultCartAdapter
		m.enableCartCache = config.EnableCartCache
		m.cartMergeStrategy = config.CartMergeStrategy
		m.enablePlaceOrderLoggerAdapter = config.EnablePlaceOrderLoggerAdapter
	}
}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.enableDefaultCartAdapter {
		injector.Bind((*infrastructure.CartStorage)(nil)).To(infrastructure.InMemoryCartStorage{}).AsEagerSingleton()
		injector.Bind((*infrastructure.GiftCardHandler)(nil)).To(infrastructure.DefaultGiftCardHandler{})
		injector.Bind((*infrastructure.VoucherHandler)(nil)).To(infrastructure.DefaultVoucherHandler{})
		injector.Bind((*cart.GuestCartService)(nil)).To(infrastructure.DefaultGuestCartService{})
		injector.Bind((*cart.CustomerCartService)(nil)).To(infrastructure.DefaultCustomerCartService{})
	}

	if m.enablePlaceOrderLoggerAdapter {
		injector.Bind((*placeorder.Service)(nil)).To(placeorderAdapter.PlaceOrderLoggerAdapter{})
	}
	// Register Default EventPublisher
	injector.Bind((*events.EventPublisher)(nil)).To(events.DefaultEventPublisher{})

	// Event
	flamingo.BindEventSubscriber(injector).To(application.EventReceiver{})

	// Cart merge strategy that is used by the event receiver to merge carts on during login
	switch m.cartMergeStrategy {
	case "replace":
		injector.Bind((*application.CartMerger)(nil)).To(application.CartMergeStrategyReplace{})
	case "none":
		injector.Bind((*application.CartMerger)(nil)).To(application.CartMergeStrategyNone{})
	default:
		injector.Bind((*application.CartMerger)(nil)).To(application.CartMergeStrategyMerge{})
	}

	// TemplateFunction
	flamingo.BindTemplateFunc(injector, "getCart", new(templatefunctions.GetCart))
	flamingo.BindTemplateFunc(injector, "getDecoratedCart", new(templatefunctions.GetDecoratedCart))
	flamingo.BindTemplateFunc(injector, "getQuantityAdjustmentDeletedItemsMessages", new(templatefunctions.GetQuantityAdjustmentDeletedItemsMessages))
	flamingo.BindTemplateFunc(injector, "getQuantityAdjustmentUpdatedItemsMessages", new(templatefunctions.GetQuantityAdjustmentUpdatedItemsMessage))
	flamingo.BindTemplateFunc(injector, "getQuantityAdjustmentCouponCodesRemoved", new(templatefunctions.GetQuantityAdjustmentCouponCodesRemoved))
	flamingo.BindTemplateFunc(injector, "removeQuantityAdjustmentMessages", new(templatefunctions.RemoveQuantityAdjustmentMessages))

	injector.Bind((*cart.DeliveryInfoBuilder)(nil)).To(cart.DefaultDeliveryInfoBuilder{})

	if m.enableCartCache {
		injector.Bind((*application.CartCache)(nil)).To(application.CartSessionCache{})
	}

	// Register Form Data Provider
	injector.BindMap(new(formDomain.FormService), "commerce.cart.deliveryFormService").To(forms.DeliveryFormService{})
	injector.BindMap(new(formDomain.FormService), "commerce.cart.billingFormService").To(forms.BillingAddressFormService{})
	injector.BindMap(new(formDomain.FormService), "commerce.cart.personaldataFormService").To(forms.DefaultPersonalDataFormService{})

	web.BindRoutes(injector, new(routes))

	injector.BindMulti(new(flamingographql.Service)).To(graphql.Service{})

	injector.Bind(new(application.Receiver)).To(application.BaseCartReceiver{})
	injector.Bind(new(application.Service)).To(application.CartService{})
}

// CueConfig defines the cart module configuration
func (*Module) CueConfig() string {
	return `
commerce: {
	cart: {
		defaultCartAdapter: {
			enabled: bool | *true
			storage: "inmemory"
			defaultTaxRate?: number
			productPrices: *"gross" | "net"
			defaultCurrency: string | *"€"
		}
		placeOrderLogger: {
			enabled: bool | *true
			useFlamingoLog: bool | *true
			logAsFile: bool | *true
			logDirectory: string | *"./orders/"
		}
		enableCartCache: bool | *true
		cacheLifetime: number | *1200
		defaultUseBillingAddress: bool | *false
		defaultDeliveryCode: string | *"delivery"
		deleteEmptyDelivery: bool | *false
		mergeStrategy: "none" | "replace" | *"merge"
		showEmptyCartPageIfNoItems?: bool
		adjustItemsToRestrictedQty?: bool
		personalDataForm: {
			additionalFormFields: [...string] | *[]
			dateOfBirthRequired: bool | *false
			passportCountryRequired: bool | *false
			passportNumberRequired: bool | *false
			minAge?: number
		}
		simplePaymentForm: {
			giftCardPaymentMethod: string | *"voucher"
		}
	}
}`
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"commerce.cart.useEmailPlaceOrderAdapter":                 "commerce.cart.placeOrderLogger.enabled",
		"commerce.cart.useInMemoryCartServiceAdapters":            "commerce.cart.defaultCartAdapter.enabled",
		"commerce.cart.inMemoryCartServiceAdapter.defaultTaxRate": "commerce.cart.defaultCartAdapter.defaultTaxRate",
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(product.Module),
		new(form.Module),
		new(customer.Module),
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
	registry.MustRoute("/cart", "cart.view")

	registry.HandleAny("cart.add", r.viewController.AddAndViewAction)
	registry.MustRoute("/cart/add/:marketplaceCode", `cart.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)

	registry.HandleAny("cart.updateQty", r.viewController.UpdateQtyAndViewAction)
	registry.MustRoute("/cart/update/:id", `cart.updateQty(id,qty?="1",deliveryCode?="")`)

	registry.HandleAny("cart.deleteAllItems", r.viewController.DeleteAllAndViewAction)
	registry.MustRoute("/cart/delete/all", "cart.deleteAllItems")

	registry.HandleAny("cart.clean", r.viewController.CleanAndViewAction)
	registry.MustRoute("/cart/clean", "cart.clean")

	registry.HandleDelete("cart.clean", r.viewController.CleanAndViewAction)
	registry.MustRoute("/cart/delivery/:deliveryCode", "cart.clean")

	registry.HandleAny("cart.cleanDelivery", r.viewController.CleanDeliveryAndViewAction)
	registry.MustRoute("/cart/delete/delivery/:deliveryCode", `cart.cleanDelivery(deliveryCode?="")`)

	registry.HandleAny("cart.deleteItem", r.viewController.DeleteAndViewAction)
	registry.MustRoute("/cart/delete/:id", `cart.deleteItem(id,deliveryCode?="")`)
	r.apiRoutes(registry)
}

func (r *routes) apiRoutes(registry *web.RouterRegistry) {
	// v1 Routes:
	registry.MustRoute("/api/v1/cart", "cart.api.cart")
	registry.HandleDelete("cart.api.cart", r.apiController.DeleteCartAction)
	registry.HandleGet("cart.api.cart", r.apiController.GetAction)

	registry.MustRoute("/api/v1/cart/billing", `cart.api.billing`)
	registry.HandlePut("cart.api.billing", r.apiController.BillingAction)

	registry.MustRoute("/api/v1/cart/payment-selection", `cart.api.payment-selection`)
	registry.HandlePut("cart.api.payment-selection", r.apiController.UpdatePaymentSelectionAction)

	registry.MustRoute("/api/v1/cart/deliveries/items", "cart.api.deliveries.items")
	registry.HandleDelete("cart.api.deliveries.items", r.apiController.DeleteAllItemsAction)

	registry.MustRoute("/api/v1/cart/delivery/:deliveryCode", `cart.api.delivery`)
	registry.HandleDelete("cart.api.delivery", r.apiController.DeleteDelivery)
	registry.HandlePut("cart.api.delivery", r.apiController.UpdateDeliveryInfoAction)

	registry.MustRoute("/api/v1/cart/delivery/:deliveryCode/item", `cart.api.item(marketplaceCode?="",variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)
	registry.HandlePost("cart.api.item", r.apiController.AddAction)
	registry.HandleDelete("cart.api.item", r.apiController.DeleteItemAction)
	registry.HandlePut("cart.api.item", r.apiController.UpdateItemAction)

	registry.MustRoute("/api/v1/cart/voucher", `cart.api.voucher(couponCode)`)
	registry.HandlePost("cart.api.voucher", r.apiController.ApplyVoucherAndGetAction)
	registry.HandleDelete("cart.api.voucher", r.apiController.RemoveVoucherAndGetAction)

	registry.MustRoute("/api/v1/cart/gift-card", `cart.api.gift-card(couponCode)`)
	registry.HandlePost("cart.api.gift-card", r.apiController.ApplyGiftCardAndGetAction)
	registry.HandleDelete("cart.api.gift-card", r.apiController.RemoveGiftCardAndGetAction)

	registry.MustRoute("/api/v1/cart/voucher-gift-card", `cart.api.voucher-gift-card(couponCode)`)
	registry.HandlePost("cart.api.voucher-gift-card", r.apiController.ApplyCombinedVoucherGift)

	// Legacy Routes:
	registry.MustRoute("/api/cart", "cart.api.get")
	registry.HandleDelete("cart.api.get", r.apiController.DeleteAllItemsAction)
	registry.HandleGet("cart.api.get", r.apiController.GetAction)

	registry.MustRoute("/api/cart/delivery/:deliveryCode/additem", `cart.api.add(marketplaceCode,variantMarketplaceCode?="",qty?="1",deliveryCode?="")`)
	registry.HandlePost("cart.api.add", r.apiController.AddAction)

	registry.MustRoute("/api/cart/applyvoucher", `cart.api.applyVoucher(couponCode)`)
	registry.HandlePost("cart.api.applyVoucher", r.apiController.ApplyVoucherAndGetAction)
	registry.HandlePut("cart.api.applyVoucher", r.apiController.ApplyVoucherAndGetAction)

	registry.MustRoute("/api/cart/removevoucher", `cart.api.removeVoucher(couponCode)`)
	registry.HandlePost("cart.api.removeVoucher", r.apiController.RemoveVoucherAndGetAction)
	registry.HandleDelete("cart.api.removeVoucher", r.apiController.RemoveVoucherAndGetAction)

	registry.MustRoute("/api/cart/applygiftcard", `cart.api.applyGiftCard(couponCode)`)
	registry.HandlePost("cart.api.applyGiftCard", r.apiController.ApplyGiftCardAndGetAction)
	registry.HandlePut("cart.api.applyGiftCard", r.apiController.ApplyGiftCardAndGetAction)

	registry.MustRoute("/api/cart/removegiftcard", `cart.api.removeGiftCard(couponCode)`)
	registry.HandlePost("cart.api.removeGiftCard", r.apiController.RemoveGiftCardAndGetAction)
	registry.HandleDelete("cart.api.removeGiftCard", r.apiController.RemoveGiftCardAndGetAction)

	registry.MustRoute("/api/cart/applycombinedvouchergift", `cart.api.applyCombinedVoucherGift(couponCode)`)
	registry.HandlePost("cart.api.applyCombinedVoucherGift", r.apiController.ApplyCombinedVoucherGift)

	registry.MustRoute("/api/cart/billing", `cart.api.billing`)
	registry.HandlePost("cart.api.billing", r.apiController.BillingAction)

	registry.MustRoute("/api/cart/delivery/:deliveryCode", `cart.api.delivery.delete`)
	registry.HandleDelete("cart.api.delivery.delete", r.apiController.DeleteDelivery)

	registry.MustRoute("/api/cart/delivery/:deliveryCode/deliveryinfo", `cart.api.delivery.update`)
	registry.HandlePost("cart.api.delivery.update", r.apiController.UpdateDeliveryInfoAction)
}
