package product

import (
	"flamingo.me/flamingo-commerce/v3/product/interfaces/controller"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// Module registers our profiler
type Module struct{}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMap(new(template.CtxFunc), "getProduct").To(templatefunctions.GetProduct{})
	injector.BindMap(new(template.Func), "getProductUrl").To(templatefunctions.GetProductUrl{})
	injector.BindMap(new(template.CtxFunc), "findProducts").To(templatefunctions.FindProducts{})

	web.BindRoutes(injector, new(routes))
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

type routes struct {
	controller *controller.View
}

func (r *routes) Inject(controller *controller.View) {
	r.controller = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("product.view", r.controller.Get)
	registry.Route("/product/:marketplacecode/:name.html", `product.view(marketplacecode, name, backurl?="")`).Normalize("name")
	registry.Route("/product/:marketplacecode/:variantcode/:name.html", `product.view(marketplacecode, variantcode, name, backurl?="")`).Normalize("name")
}
