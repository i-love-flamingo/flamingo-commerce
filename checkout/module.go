package checkout

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/checkout/domain"
	paymentDomain "flamingo.me/flamingo-commerce/v3/checkout/domain/payment"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure"
	paymentInfrastructure "flamingo.me/flamingo-commerce/v3/checkout/infrastructure/payment"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller/formdto"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/go-playground/form"
	formDomain "go.aoe.com/flamingo/form/domain"
)

type (
	// Module registers our profiler
	Module struct {
		UseFakeSourcingService bool `inject:"config:checkout.useFakeSourcingService,optional"`
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMap((*paymentDomain.Provider)(nil), "offlinepayment").To(paymentInfrastructure.OfflinePaymentProvider{})

	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeSourcingService {
		injector.Override((*domain.SourcingService)(nil), "").To(infrastructure.FakeSourcingService{})
	}

	injector.Bind((*formDomain.FormService)(nil)).To(formdto.CheckoutFormService{})

	web.BindRoutes(injector, new(routes))
}

// DefaultConfig for checkout module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"checkout": config.Map{
			"defaultPaymentMethod": "checkmo",
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
	registry.Route("/checkout", "checkout.start")

	registry.HandleAny("checkout.review", r.controller.ReviewAction)
	registry.Route("/checkout/review", `checkout.review`)

	registry.HandleAny("checkout.guest", r.controller.SubmitGuestCheckoutAction)
	registry.Route("/checkout/guest", "checkout.guest")

	registry.HandleAny("checkout.user", r.controller.SubmitUserCheckoutAction)
	registry.Route("/checkout/user", "checkout.user")

	registry.HandleAny("checkout.success", r.controller.SuccessAction)
	registry.Route("/checkout/success", "checkout.success")

	registry.HandleAny("checkout.expired", r.controller.ExpiredAction)
	registry.Route("/checkout/expired", "checkout.expired")

	registry.HandleAny("checkout.processpayment", r.controller.ProcessPaymentAction)
	registry.Route("/checkout/processpayment/:providercode/:methodcode", "checkout.processpayment")
}
