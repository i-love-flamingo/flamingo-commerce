package order

import (
	"go.aoe.com/flamingo/core/order/interfaces/controller"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// OrdersModule for Orders
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("customerorders", new(controller.DataControllerCustomerOrders))
}
