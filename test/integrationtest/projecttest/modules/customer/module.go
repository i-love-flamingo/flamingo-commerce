package customer

import (
	"flamingo.me/dingo"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
)

type (
	// Module registers our fake customer order module
	Module struct{}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.CustomerService)).To(new(FakeService))
	injector.Bind(new(domain.CustomerIdentityService)).To(new(FakeService))
}
