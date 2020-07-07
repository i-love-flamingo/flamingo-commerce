package sourcing

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	restrictors "flamingo.me/flamingo-commerce/v3/sourcing/domain/restrictor"

	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/sourcing/application"
	"flamingo.me/flamingo-commerce/v3/sourcing/domain"
)

type (
	// Module registers sourcing module
	Module struct {
		useDefaultSourcingService bool
		enableQtyRestrictor       bool
	}
)

// Inject dependencies
func (m *Module) Inject(
	config *struct {
		UseDefaultSourcingService bool `inject:"config:commerce.sourcing.useDefaultSourcingService,optional"`
		EnableQtyRestrictor       bool `inject:"config:commerce.sourcing.enableQtyRestrictor,optional"`
	},
) {

	if config != nil {
		m.useDefaultSourcingService = config.UseDefaultSourcingService
		m.enableQtyRestrictor = config.EnableQtyRestrictor
	}

}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.useDefaultSourcingService {
		injector.Bind(new(domain.SourcingService)).To(domain.DefaultSourcingService{})
	}

	if m.enableQtyRestrictor {
		injector.Bind(new(validation.MaxQuantityRestrictor)).To(restrictors.Restrictor{})
	}

	injector.Bind(new(application.SourcingApplication)).To(application.Service{})
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(cart.Module),
	}
}

// CueConfig defines the sourcing module configuration
func (m *Module) CueConfig() string {
	return `
commerce: {
	sourcing: {
		useDefaultSourcingService: bool | *true
		enableQtyRestrictor: bool | *false
	}
}
`
}
