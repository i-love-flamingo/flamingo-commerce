package order

import (
	"flamingo.me/flamingo-commerce/order/application"
	"flamingo.me/flamingo-commerce/order/domain"
	"flamingo.me/flamingo-commerce/order/infrastructure/fake"
	"flamingo.me/flamingo-commerce/order/infrastructure/inmemory"
	"flamingo.me/flamingo-commerce/order/interfaces/controller"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
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
		UseInMemoryService bool `inject:"config:order.useInMemoryService,optional"`
	},
) {
	m.logger = logger
	if config != nil {
		m.useFakeAdapters = config.UseFakeAdapters
		m.useInMemoryService = config.UseInMemoryService
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	// check if the configuration is used with care and sense
	if m.useFakeAdapters && m.useInMemoryService {
		m.logger.WithField(flamingo.LogKeyCategory, LogKey).Panic("fake adapters _and_ inmemory service are both activated - please choose none or only one of them")
		panic("fake adapters _and_ inmemory service are both activated - please choose none or only one of them")
	}

	if m.useFakeAdapters {
		injector.Bind((*domain.GuestOrderService)(nil)).To(fake.GuestOrders{})
		injector.Bind((*domain.CustomerOrderService)(nil)).To(fake.CustomerOrders{})
	}
	if m.useInMemoryService {
		injector.Bind((*inmemory.Storager)(nil)).To(inmemory.Storage{}).AsEagerSingleton()
		injector.Bind((*domain.GuestOrderService)(nil)).To(inmemory.GuestOrderService{})
		injector.Bind((*domain.CustomerOrderService)(nil)).To(inmemory.CustomerOrderService{})

	}
	injector.Bind((*application.EventPublisher)(nil)).To(application.DefaultEventPublisher{})
	injector.Bind((*domain.OrderDecoratorInterface)(nil)).To(domain.OrderDecorator{})
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
