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
		UseFakeAdapters bool `inject:"config:order.useFakeAdapters,optional"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.OrderDecoratorInterface)(nil)).To(domain.OrderDecorator{})
	if m.UseFakeAdapters {
		injector.Bind((*domain.CustomerOrderService)(nil)).To(infrastructure.FakeCustomerOrders{})
	}
	router.Bind(injector, new(routes))
}

type routes struct {
	controller *controller.DataControllerCustomerOrders
}

func (r *routes) Inject(controller *controller.DataControllerCustomerOrders) {
	r.controller = controller
}

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleData("customerorders", r.controller.Data)
}
