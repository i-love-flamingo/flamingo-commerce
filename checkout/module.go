package checkout

import (
	"github.com/go-playground/form"
	"flamingo.me/flamingo-commerce/checkout/domain"
	paymentDomain "flamingo.me/flamingo-commerce/checkout/domain/payment"
	"flamingo.me/flamingo-commerce/checkout/infrastructure"
	paymentInfrastructure "flamingo.me/flamingo-commerce/checkout/infrastructure/payment"
	"flamingo.me/flamingo-commerce/checkout/interfaces/controller"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

type (
	// CheckoutModule registers our profiler
	CheckoutModule struct {
		RouterRegistry                  *router.Registry `inject:""`
		UseFakeDeliveryLocationsService bool             `inject:"config:checkout.useFakeDeliveryLocationsService,optional"`
	}
)

// Configure module
func (m *CheckoutModule) Configure(injector *dingo.Injector) {

	m.RouterRegistry.Handle("checkout.start", (*controller.CheckoutController).StartAction)
	m.RouterRegistry.Route("/checkout", "checkout.start")

	m.RouterRegistry.Handle("checkout.review", (*controller.CheckoutController).ReviewAction)
	m.RouterRegistry.Route("/checkout/review", `checkout.review`)

	m.RouterRegistry.Handle("checkout.guest", (*controller.CheckoutController).SubmitGuestCheckoutAction)
	m.RouterRegistry.Route("/checkout/guest", "checkout.guest")

	m.RouterRegistry.Handle("checkout.user", (*controller.CheckoutController).SubmitUserCheckoutAction)
	m.RouterRegistry.Route("/checkout/user", "checkout.user")

	m.RouterRegistry.Handle("checkout.success", (*controller.CheckoutController).SuccessAction)
	m.RouterRegistry.Route("/checkout/success", "checkout.success")

	m.RouterRegistry.Handle("checkout.expired", (*controller.CheckoutController).ExpiredAction)
	m.RouterRegistry.Route("/checkout/expired", "checkout.expired")

	m.RouterRegistry.Handle("checkout.processpayment", (*controller.CheckoutController).ProcessPaymentAction)
	m.RouterRegistry.Route("/checkout/processpayment/:providercode/:methodcode", "checkout.processpayment")

	injector.BindMap((*paymentDomain.PaymentProvider)(nil), "offlinepayment").To(paymentInfrastructure.OfflinePaymentProvider{})

	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeDeliveryLocationsService {
		injector.Override((*domain.SourcingService)(nil), "").To(infrastructure.FakeSourcingService{})
	}
}

// DefaultConfig
func (m *CheckoutModule) DefaultConfig() config.Map {
	return config.Map{
		"checkout": config.Map{
			"defaultPaymentMethod": "checkmo",
		},
	}
}
