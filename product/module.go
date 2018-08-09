package product

import (
	"flamingo.me/flamingo-commerce/product/interfaces/controller"
	"flamingo.me/flamingo-commerce/product/interfaces/templatefunctions"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/template"
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
	m.RouterRegistry.Route("/product/:marketplacecode/:name.html", `product.view(marketplacecode, name, backurl?="")`).Normalize("marketplacecode", "name")
	m.RouterRegistry.Route("/product/:marketplacecode/:variantcode/:name.html", `product.view(marketplacecode, variantcode, name, backurl?="")`).Normalize("marketplacecode", "name", "variantcode")

	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.GetProduct{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.GetProductUrl{})
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.product.view.template": "product/product",
		"templating": config.Map{
			"product": config.Map{
				"attributeRenderer": config.Map{},
			},
		},
	}
}
