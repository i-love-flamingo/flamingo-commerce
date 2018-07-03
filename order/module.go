package order

import (
	"flamingo.me/flamingo-commerce/order/domain"
	"flamingo.me/flamingo-commerce/order/infrastructure"
	"flamingo.me/flamingo-commerce/order/interfaces/controller"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
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
	injector.Bind((*domain.OrderDecoratorInterface)(nil)).To(domain.OrderDecorator{})
	if m.UseFakeAdapters {
		injector.Bind((*domain.CustomerOrderService)(nil)).To(infrastructure.FakeCustomerOrders{})
	}

}
