package product

import (
	"flamingo/core/product/interfaces/controller"
	"flamingo/framework/config"
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
	m.RouterRegistry.Handle("product.view", new(controller.View))
	m.RouterRegistry.Route("/product/:marketplacecode/:name.html", "product.view")
	m.RouterRegistry.Route("/product/:marketplacecode/:variantcode/:name.html", "product.view")
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.product.view.template": "product/product",
	}
}
