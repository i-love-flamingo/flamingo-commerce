package customer

import (
	"flamingo.me/dingo"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
)

type (
	// Module is a fake customer module
	Module struct{}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.CustomerIdentityService)).To(new(FakeService))
}
