package checkout

import (
	"github.com/go-playground/form"
	"go.aoe.com/flamingo/core/checkout/application"
	"go.aoe.com/flamingo/core/checkout/infrastructure"
	"go.aoe.com/flamingo/core/checkout/interfaces/controller"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
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

	m.RouterRegistry.Handle("checkout.guest", (*controller.CheckoutController).SubmitGuestCheckoutAction)
	m.RouterRegistry.Route("/checkout/guest", "checkout.guest")

	m.RouterRegistry.Handle("checkout.user", (*controller.CheckoutController).SubmitUserCheckoutAction)
	m.RouterRegistry.Route("/checkout/user", "checkout.user")

	m.RouterRegistry.Handle("checkout.success", (*controller.CheckoutController).SuccessAction)
	m.RouterRegistry.Route("/checkout/success", "checkout.success")

	m.RouterRegistry.Handle("checkout.processpayment", (*controller.CheckoutController).ProcessPaymentAction)
	m.RouterRegistry.Route("/checkout/processpayment/:providercode/:methodcode", "checkout.processpayment")

	injector.Bind((*form.Decoder)(nil)).ToProvider(form.NewDecoder).AsEagerSingleton()
	if m.UseFakeDeliveryLocationsService {
		injector.Bind((*application.SourcingService)(nil)).To(infrastructure.FakeSourcingService{})
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
