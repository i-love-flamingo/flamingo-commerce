package checkout

import (
	"go.aoe.com/flamingo/core/checkout/interfaces/controller"
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

	m.RouterRegistry.Handle("checkout.start", (*controller.CheckoutController).StartAction)
	m.RouterRegistry.Route("/checkout", "checkout.start")

	m.RouterRegistry.Handle("checkout.submit", (*controller.CheckoutController).SubmitAction)
	m.RouterRegistry.Route("/checkout", "checkout.submit")

}
