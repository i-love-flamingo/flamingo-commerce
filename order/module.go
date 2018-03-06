package order

import (
	"go.aoe.com/flamingo/core/order/domain"
	"go.aoe.com/flamingo/core/order/infrastructure"
	"go.aoe.com/flamingo/core/order/interfaces/controller"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// OrdersModule for Orders
	Module struct {
		RouterRegistry  *router.Registry `inject:""`
		UseFakeAdapters bool             `inject:"config:order.useFakeAdapters,optional"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("customerorders", new(controller.DataControllerCustomerOrders))
	if m.UseFakeAdapters {
		injector.Bind((*domain.CustomerOrderService)(nil)).To(infrastructure.FakeCustomerOrders{})
	}

}
