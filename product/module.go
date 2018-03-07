package product

import (
	"go.aoe.com/flamingo/core/product/interfaces/controller"
	"go.aoe.com/flamingo/core/product/interfaces/templatefunctions"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/template"
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

	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.GetProduct{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.GetProductUrl{})
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.product.view.template": "product/product",
	}
}
