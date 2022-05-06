package order

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/order/domain"
	"flamingo.me/flamingo-commerce/v3/order/infrastructure/fake"
	"flamingo.me/flamingo-commerce/v3/order/interfaces/controller"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module definition of the order module
	Module struct {
		useFakeAdapter bool
	}
)

const (
	// LogKey defines the module log category key
	LogKey = "order"
)

// Inject dependencies
func (m *Module) Inject(
	config *struct {
		UseFakeAdapter bool `inject:"config:commerce.order.useFakeAdapter,optional"`
	},
) {
	if config != nil {
		m.useFakeAdapter = config.UseFakeAdapter
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {

	if m.useFakeAdapter {
		injector.Bind((*domain.CustomerIdentityOrderService)(nil)).To(fake.CustomerOrders{})
	}

	injector.Bind((*domain.OrderDecoratorInterface)(nil)).To(domain.OrderDecorator{})
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	controller *controller.DataControllerCustomerOrders
}

func (r *routes) Inject(controller *controller.DataControllerCustomerOrders) {
	r.controller = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleData("customerorders", r.controller.Data)
}

// FlamingoLegacyConfigAlias maps legacy config entries to new ones
func (m *Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"order.useFakeAdapters": "commerce.order.useFakeAdapter",
	}
}
