package checkout

import (
	"go.aoe.com/flamingo/core/checkout/interfaces/controller"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (

	// CheckoutModule registers our profiler
	CheckoutModule struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure module
func (m *CheckoutModule) Configure(injector *dingo.Injector) {

	m.RouterRegistry.Handle("checkout.start", (*controller.CheckoutController).SubmitAction)
	m.RouterRegistry.Route("/checkout", "checkout.start")

	m.RouterRegistry.Handle("checkout.success", (*controller.CheckoutController).SuccessAction)
	m.RouterRegistry.Route("/checkout/success", "checkout.success")

}

// DefaultConfig
func (m *CheckoutModule) DefaultConfig() config.Map {
	return config.Map{
		"checkout": config.Map{
			"defaultPaymentMethod": "checkmo",
		},
	}
}
