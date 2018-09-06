package checkout

import (
	"flamingo.me/flamingo-commerce/checkout/domain"
	paymentDomain "flamingo.me/flamingo-commerce/checkout/domain/payment"
	"flamingo.me/flamingo-commerce/checkout/infrastructure"
	paymentInfrastructure "flamingo.me/flamingo-commerce/checkout/infrastructure/payment"
	"flamingo.me/flamingo-commerce/checkout/interfaces/controller"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
	"github.com/go-playground/form"
)

type (
	// CheckoutModule registers our profiler
	CheckoutModule struct {
		UseFakeSourcingService bool `inject:"config:checkout.useFakeSourcingService,optional"`
	}
)

// Configure module
func (m *CheckoutModule) Configure(injector *dingo.Injector) {
	injector.BindMap((*paymentDomain.PaymentProvider)(nil), "offlinepayment").To(paymentInfrastructure.OfflinePaymentProvider{})

	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeSourcingService {
		injector.Override((*domain.SourcingService)(nil), "").To(infrastructure.FakeSourcingService{})
	}

	router.Bind(injector, new(routes))
}

// DefaultConfig
func (m *CheckoutModule) DefaultConfig() config.Map {
	return config.Map{
		"checkout": config.Map{
			"defaultPaymentMethod": "checkmo",
		},
	}
}

type routes struct {
	controller *controller.CheckoutController
}

func (r *routes) Inject(controller *controller.CheckoutController) {
	r.controller = controller
}

func (r *routes) Routes(registry *router.Registry) {
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
