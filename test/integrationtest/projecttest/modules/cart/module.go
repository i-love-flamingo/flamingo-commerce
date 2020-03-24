package cart

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
)

type (
	// Module for integration testing
	Module struct{}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Override((*infrastructure.VoucherHandler)(nil), "").To(&FakeVoucherHandler{})
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		&cart.Module{},
	}
}
