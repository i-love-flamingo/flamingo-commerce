package order

import (
	"flamingo.me/flamingo-commerce/v3/order/domain"
	"flamingo.me/flamingo-commerce/v3/order/infrastructure/fake"
	"flamingo.me/flamingo-commerce/v3/order/interfaces/controller"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module definition of the order module
	Module struct {
		logger             flamingo.Logger
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
	logger flamingo.Logger,
	config *struct {
		UseFakeAdapters    bool `inject:"config:order.useFakeAdapters,optional"`
	},
) {
	m.logger = logger
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
