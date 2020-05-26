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
		useFakeAdapters    bool
		useInMemoryService bool
	}
)

const (
	// LogKey defines the module log category key
	LogKey = "order"
)

// Inject dependencies
func (m *Module) Inject(
	config *struct {
		UseFakeAdapters bool `inject:"config:order.useFakeAdapters,optional"`
	},
) {
	if config != nil {
		m.useFakeAdapters = config.UseFakeAdapters
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {

	if m.useFakeAdapters {
		injector.Bind((*domain.CustomerOrderService)(nil)).To(fake.CustomerOrders{})
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
