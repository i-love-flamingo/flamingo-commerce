package product

import (
	"flamingo/core/product/interfaces"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("product.view", new(interfaces.ViewController))
	m.RouterRegistry.Route("/product/:marketplacecode/:name.html", "product.view")
	m.RouterRegistry.Route("/product/:marketplacecode/:variantcode/:name.html", "product.view")
}
