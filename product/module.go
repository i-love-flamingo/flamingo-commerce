package product

import (
	"flamingo/core/breadcrumbs"
	"flamingo/core/product/interfaces/controller"
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

	m.RouterRegistry.Handle("product.debug.data", new(controller.DebugData))
	m.RouterRegistry.Route("/product-debug/:marketplacecode", "product.debug.data")
	m.RouterRegistry.Route("/product-debug/:marketplacecode/:variantcode", "product.debug.data")
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"core.product.view.template": "product/product",
	}
}

func (m *Module) Dependencies() []dingo.Module {
	return []dingo.Module{
		new(breadcrumbs.Module),
	}
}
