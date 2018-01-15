package order

import (
	"go.aoe.com/flamingo/core/magento/infrastructure/orderservice"
	"go.aoe.com/flamingo/core/order/domain"
	"go.aoe.com/flamingo/core/order/interfaces/controller"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// OrdersModule for Orders
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		Debug          bool             `inject:"config:debug.mode"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.CustomerOrderService)(nil)).To(orderservice.CustomerOrders{})
	m.RouterRegistry.Handle("customerorders", new(controller.DataControllerCustomerOrders))
}
