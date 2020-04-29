package customer

import (
	"flamingo.me/dingo"
	customerDomain "flamingo.me/flamingo-commerce/v3/customer/domain"
	customerInfrastructure "flamingo.me/flamingo-commerce/v3/customer/infrastructure"
	customerGraphql "flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql"
	flamingoGraphql "flamingo.me/graphql"
)

type (
	// Module registers our customer module
	Module struct {
		useNilCustomerAdapter bool
	}
)

// Inject  module
func (m *Module) Inject(config *struct {
	UseNilCustomerAdapter bool `inject:"config:commerce.customer.useNilCustomerAdapter,optional"`
}) {
	if config != nil {
		m.useNilCustomerAdapter = config.UseNilCustomerAdapter
	}
}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	if m.useNilCustomerAdapter {
		injector.Bind((*customerDomain.CustomerService)(nil)).To(customerInfrastructure.NilCustomerServiceAdapter{})
		injector.Bind((*customerDomain.CustomerIdentityService)(nil)).To(customerInfrastructure.NilCustomerServiceAdapter{})
	}
	injector.BindMulti(new(flamingoGraphql.Service)).To(customerGraphql.Service{})
}
