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

	//Basti - backward generation (e.g. for redirect) does not work with alternative route - only with new handler
	m.RouterRegistry.Handle("product.view.variant", new(interfaces.ViewController))
	m.RouterRegistry.Route("/product/:marketplacecode/:variantcode/:name.html", "product.view.variant")
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"core.product.view.template": "product/product",
	}
}
