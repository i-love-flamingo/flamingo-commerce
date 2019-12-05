package product

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/controller"
	productgraphql "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/templatefunctions"
	"flamingo.me/flamingo-commerce/v3/search"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/graphql"
)

// Module registers our profiler
type Module struct {
	FakeService bool `inject:"config:commerce.product.fakeservice.enabled,optional"`
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "getProduct", new(templatefunctions.GetProduct))
	flamingo.BindTemplateFunc(injector, "getProductUrl", new(templatefunctions.GetProductURL))
	flamingo.BindTemplateFunc(injector, "findProducts", new(templatefunctions.FindProducts))

	web.BindRoutes(injector, new(routes))

	injector.BindMulti(new(graphql.Service)).To(new(productgraphql.Service))
	if m.FakeService {
		injector.Bind((*domain.ProductService)(nil)).To(fake.ProductService{})
		injector.Bind((*domain.SearchService)(nil)).To(fake.SearchService{})
	}

}

// Depends adds our dependencies
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		price.Module{},
		new(search.Module),
	}
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"commerce": config.Map{
			"product": config.Map{
				"view": config.Map{
					"template": "product/product",
				},
				"priceIsGross":      true,
				"generateSlug":      true,
				"slugAttributeCode": "urlSlug",
			},
		},
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
	h, _ := registry.Route("/product/:marketplacecode/:name.html", `product.view(marketplacecode, name, backurl?="")`)
	h.Normalize("name")
	h, _ = registry.Route("/product/:marketplacecode/:variantcode/:name.html", `product.view(marketplacecode, variantcode, name, backurl?="")`)
	h.Normalize("name")
}
