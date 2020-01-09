package checkout

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/locker"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	flamingographql "flamingo.me/graphql"
	"github.com/go-playground/form"
)

type (
	// Module registers our profiler
	Module struct {
		UseFakeSourcingService bool   `inject:"config:checkout.useFakeSourcingService,optional"`
		PlaceOrderLockType     string `inject:"config:checkout.placeorder.lockType,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {

	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeSourcingService {
		injector.Override((*domain.SourcingService)(nil), "").To(infrastructure.FakeSourcingService{})
	}

	if m.PlaceOrderLockType == "clusterlock" {
		//TODO provide cluster imple
		injector.Bind(new(placeorder.TryLock)).To(locker.Simple{})
	} else {
		injector.Bind(new(placeorder.TryLock)).To(locker.Simple{})
	}

	injector.Bind(new(process.State)).AnnotatedWith("startState").To(states.New{})

	web.BindRoutes(injector, new(routes))

	injector.BindMulti(new(flamingographql.Service)).To(graphql.Service{})
}

// DefaultConfig for checkout module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"checkout": config.Map{
			"useDeliveryForms":    true,
			"usePersonalDataForm": false,
			"skipReviewAction":    false,
		},
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
	registry.Route("/checkout/start", "checkout.start")

	registry.HandleAny("checkout.review", r.controller.ReviewAction)
	registry.Route("/checkout/review", `checkout.review`)

	registry.HandleAny("checkout", r.controller.SubmitCheckoutAction)
	registry.Route("/checkout", "checkout")

	registry.HandleAny("checkout.payment", r.controller.PaymentAction)
	registry.Route("/checkout/payment", "checkout.payment")

	registry.HandleAny("checkout.success", r.controller.SuccessAction)
	registry.Route("/checkout/success", "checkout.success")

	registry.HandleAny("checkout.expired", r.controller.ExpiredAction)
	registry.Route("/checkout/expired", "checkout.expired")

	registry.HandleAny("checkout.placeorder", r.controller.PlaceOrderAction)
	registry.Route("/checkout/placeorder", "checkout.placeorder")
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(cart.Module),
	}
}
