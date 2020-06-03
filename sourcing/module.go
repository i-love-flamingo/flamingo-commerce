package sourcing

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/sourcing/application"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain"
)

type (
	// Module registers sourcing module
	Module struct {
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.SourcingService)).To(domain.DefaultSourcingService{})
	injector.Bind(new(application.SourcingApplication)).To(application.Service{})
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(cart.Module),
	}
}
