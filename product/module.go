package product

import (
	"flamingo/core/product/interfaces"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("product.view", new(interfaces.ViewController))
	m.RouterRegistry.Route("/product/{uid}/{name}.html", "product.view")
}
