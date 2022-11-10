package checkout

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	flamingographql "flamingo.me/graphql"
	"github.com/go-playground/form/v4"

	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/contextstore"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/payment"

	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/locker"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql"
)

type (
	// Module registers our checkout module
	Module struct {
		UseFakeSourcingService bool   `inject:"config:commerce.checkout.useFakeSourcingService,optional"`
		PlaceOrderLockType     string `inject:"config:commerce.checkout.placeorder.lock.type"`
		PlaceOrderContextStore string `inject:"config:commerce.checkout.placeorder.contextstore.type"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeSourcingService {
		injector.Override((*domain.SourcingService)(nil), "").To(infrastructure.FakeSourcingService{})
	}

	if m.PlaceOrderLockType == "redis" {
		injector.Bind(new(locker.Redis)).ToProvider(locker.NewRedis).In(dingo.Singleton)
		injector.Bind(new(placeorder.TryLocker)).To(new(locker.Redis))
		injector.BindMap(new(healthcheck.Status), "placeorder.locker.redis").To(new(locker.Redis))
	} else {
		injector.Bind(new(placeorder.TryLocker)).ToProvider(locker.NewMemory).In(dingo.Singleton)
	}

	if m.PlaceOrderContextStore == "redis" {
		injector.Bind(new(contextstore.Redis)).In(dingo.Singleton)
		injector.Bind(new(process.ContextStore)).To(new(contextstore.Redis))
		injector.BindMap(new(healthcheck.Status), "placeorder.contextstore.redis").To(new(contextstore.Redis))
	} else {
		injector.Bind(new(process.ContextStore)).To(new(contextstore.Memory)).In(dingo.Singleton)
	}

	injector.Bind(new(process.PaymentValidatorFunc)).ToInstance(placeorder.PaymentValidator)

	injector.Bind(new(process.State)).AnnotatedWith("startState").To(states.New{})
	injector.Bind(new(process.State)).AnnotatedWith("failedState").To(states.Failed{})
	injector.BindMap(new(process.State), new(states.New).Name()).To(states.New{})
	injector.BindMap(new(process.State), new(states.PrepareCart).Name()).To(states.PrepareCart{})
	injector.BindMap(new(process.State), new(states.ValidateCart).Name()).To(states.ValidateCart{})
	injector.BindMap(new(process.State), new(states.ValidatePaymentSelection).Name()).To(states.ValidatePaymentSelection{})
	injector.BindMap(new(process.State), new(states.CreatePayment).Name()).To(states.CreatePayment{})
	injector.BindMap(new(process.State), new(states.CompleteCart).Name()).To(states.CompleteCart{})
	injector.BindMap(new(process.State), new(states.CompletePayment).Name()).To(states.CompletePayment{})
	injector.BindMap(new(process.State), new(states.PlaceOrder).Name()).To(states.PlaceOrder{})
	injector.BindMap(new(process.State), new(states.ValidatePayment).Name()).To(states.ValidatePayment{})
	injector.BindMap(new(process.State), new(states.WaitForCustomer).Name()).To(states.WaitForCustomer{})
	injector.BindMap(new(process.State), new(states.Success).Name()).To(states.Success{})
	injector.BindMap(new(process.State), new(states.Failed).Name()).To(states.Failed{})
	injector.BindMap(new(process.State), new(states.ShowIframe).Name()).To(states.ShowIframe{})
	injector.BindMap(new(process.State), new(states.ShowHTML).Name()).To(states.ShowHTML{})
	injector.BindMap(new(process.State), new(states.ShowWalletPayment).Name()).To(states.ShowWalletPayment{})
	injector.BindMap(new(process.State), new(states.Redirect).Name()).To(states.Redirect{})
	injector.BindMap(new(process.State), new(states.PostRedirect).Name()).To(states.PostRedirect{})
	injector.BindMap(new(process.State), new(states.TriggerClientSDK).Name()).To(states.TriggerClientSDK{})

	// bind internal states to graphQL states
	injector.BindMap(new(dto.State), new(states.New).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.PrepareCart).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.ValidateCart).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.ValidatePaymentSelection).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.CreatePayment).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.CompleteCart).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.CompletePayment).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.PlaceOrder).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.ValidatePayment).Name()).To(dto.Wait{})
	injector.BindMap(new(dto.State), new(states.WaitForCustomer).Name()).To(dto.WaitForCustomer{})
	injector.BindMap(new(dto.State), new(states.Success).Name()).To(dto.Success{})
	injector.BindMap(new(dto.State), new(states.Failed).Name()).To(dto.Failed{})
	injector.BindMap(new(dto.State), new(states.ShowHTML).Name()).To(dto.ShowHTML{})
	injector.BindMap(new(dto.State), new(states.ShowWalletPayment).Name()).To(dto.ShowWalletPayment{})
	injector.BindMap(new(dto.State), new(states.ShowIframe).Name()).To(dto.ShowIframe{})
	injector.BindMap(new(dto.State), new(states.Redirect).Name()).To(dto.Redirect{})
	injector.BindMap(new(dto.State), new(states.PostRedirect).Name()).To(dto.PostRedirect{})
	injector.BindMap(new(dto.State), new(states.TriggerClientSDK).Name()).To(dto.TriggerClientSDK{})

	web.BindRoutes(injector, new(routes))
	web.BindRoutes(injector, new(apiRoutes))

	injector.BindMulti(new(flamingographql.Service)).To(graphql.Service{})
}

// CueConfig definition
func (m *Module) CueConfig() string {
	// language=cue
	return `
commerce: checkout: {
	Redis :: {
		maxIdle:                 number | *25
		idleTimeoutMilliseconds: number | *240000
		network:                 string | *"tcp"
		address:                 string | *"localhost:6379"
		database:                number | *0
		ttl:                     string | *"2h"
	}
	activateDeprecatedSourcing:		    bool | *false
	useDeliveryForms:                 bool | *true
	usePersonalDataForm:              bool | *false
	skipReviewAction:                 bool | *false
	skipStartAction?:                 bool 
	showReviewStepAfterPaymentError?: bool
	showEmptyCartPageIfNoItems?:      bool
	redirectToCartOnInvalidCart?:     bool
	privacyPolicyRequired?:           bool
	placeorder: {
		lock: {
			type: *"memory" | "redis"
			if type == "redis" {
				redis: Redis
			}
		}
		contextstore: {
			type: *"memory" | "redis"
			if type == "redis" {
				redis: Redis
			}
		}
		states: {
			placeorder: {
				cancelOrdersDuringRollback: bool | *false		
			}
		}
	}
}`
}

// FlamingoLegacyConfigAlias mapping
func (m *Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"checkout.useDeliveryForms":                "commerce.checkout.useDeliveryForms",
		"checkout.usePersonalDataForm":             "commerce.checkout.usePersonalDataForm",
		"checkout.skipReviewAction":                "commerce.checkout.skipReviewAction",
		"checkout.skipStartAction":                 "commerce.checkout.skipStartAction",
		"checkout.showReviewStepAfterPaymentError": "commerce.checkout.showReviewStepAfterPaymentError",
		"checkout.showEmptyCartPageIfNoItems":      "commerce.checkout.showEmptyCartPageIfNoItems",
		"checkout.redirectToCartOnInvalideCart":    "commerce.checkout.redirectToCartOnInvalidCart",
		"checkout.privacyPolicyRequired":           "commerce.checkout.privacyPolicyRequired",
	}
}

type routes struct {
	controller *controller.CheckoutController
}

// Inject required controller
func (r *routes) Inject(controller *controller.CheckoutController) {
	r.controller = controller
}

// Routes  configuration for checkout controllers
func (r *routes) Routes(registry *web.RouterRegistry) {
	// routes
	registry.HandleAny("checkout.start", r.controller.StartAction)
	registry.MustRoute("/checkout/start", "checkout.start")

	registry.HandleAny("checkout.review", r.controller.ReviewAction)
	registry.MustRoute("/checkout/review", `checkout.review`)

	registry.HandleAny("checkout", r.controller.SubmitCheckoutAction)
	registry.MustRoute("/checkout", "checkout")

	registry.HandleAny("checkout.payment", r.controller.PaymentAction)
	registry.MustRoute("/checkout/payment", "checkout.payment")

	registry.HandleAny("checkout.success", r.controller.SuccessAction)
	registry.MustRoute("/checkout/success", "checkout.success")

	registry.HandleAny("checkout.expired", r.controller.ExpiredAction)
	registry.MustRoute("/checkout/expired", "checkout.expired")

	registry.HandleAny("checkout.placeorder", r.controller.PlaceOrderAction)
	registry.MustRoute("/checkout/placeorder", "checkout.placeorder")
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(cart.Module),
		new(payment.Module),
		new(flamingo.SessionModule),
	}
}

type apiRoutes struct {
	apiController *controller.APIController
}

func (r *apiRoutes) Inject(apiController *controller.APIController) {
	r.apiController = apiController
}

func (r *apiRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/api/v1/checkout/placeorder", "checkout.api.placeorder")
	registry.HandleGet("checkout.api.placeorder", r.apiController.CurrentPlaceOrderContextAction)
	registry.HandlePut("checkout.api.placeorder", r.apiController.StartPlaceOrderAction)
	registry.HandleDelete("checkout.api.placeorder", r.apiController.ClearPlaceOrderAction)

	registry.MustRoute("/api/v1/checkout/placeorder/cancel", "checkout.api.placeorder.cancel")
	registry.HandlePost("checkout.api.placeorder.cancel", r.apiController.CancelPlaceOrderAction)

	registry.MustRoute("/api/v1/checkout/placeorder/refresh", "checkout.api.placeorder.refresh")
	registry.HandlePost("checkout.api.placeorder.refresh", r.apiController.RefreshPlaceOrderAction)

	registry.MustRoute("/api/v1/checkout/placeorder/refresh-blocking", "checkout.api.placeorder.refreshblocking")
	registry.HandlePost("checkout.api.placeorder.refreshblocking", r.apiController.RefreshPlaceOrderBlockingAction)
}
