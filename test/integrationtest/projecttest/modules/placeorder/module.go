package placeorder

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
)

type (
	// Module registers our fake place order module
	Module struct {
	}
)

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*placeorder.Service)(nil)).To(&FakeAdapter{}).In(dingo.Singleton)
}
