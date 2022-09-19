package product

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/graphql"

	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/product/infrastructure/fake"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/controller"
	productgraphql "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/templatefunctions"
	"flamingo.me/flamingo-commerce/v3/search"
)

// Module represents the product module
type Module struct {
	fakeService bool
	api         bool
}

// Inject module configuration
func (m *Module) Inject(
	cfg *struct {
		FakeService bool `inject:"config:commerce.product.fakeservice.enabled,optional"`
		API         bool `inject:"config:commerce.product.api.enabled,optional"`
	},
) *Module {
	if cfg != nil {
		m.api = cfg.API
		m.fakeService = cfg.FakeService
	}

	return m
}

// Configure the product module
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "getProduct", new(templatefunctions.GetProduct))
	flamingo.BindTemplateFunc(injector, "getProductUrl", new(templatefunctions.GetProductURL))
	flamingo.BindTemplateFunc(injector, "findProducts", new(templatefunctions.FindProducts))

	web.BindRoutes(injector, new(routes))
	if m.api {
		web.BindRoutes(injector, new(apiRoutes))
	}

	injector.BindMulti(new(graphql.Service)).To(new(productgraphql.Service))
	if m.fakeService {
		injector.Override((*domain.ProductService)(nil), "").To(fake.ProductService{}).In(dingo.ChildSingleton)
		injector.Override((*domain.SearchService)(nil), "").To(fake.SearchService{}).In(dingo.ChildSingleton)
	}

}

// Depends adds our dependencies
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(price.Module),
		new(search.Module),
	}
}

// CueConfig defines the product module configuration
func (*Module) CueConfig() string {
	// language=cue
	return `
commerce: {
	product: {
		view:  {
			template: *"product/product" | !=""
		}
		priceIsGross: bool | *true
		generateSlug: bool | *true
		slugAttributeCode: string | *"urlSlug"
		fakeservice: {
			enabled: bool | *false
			currency: *"€" | !=""
			defaultProducts: bool | *true
			if enabled {
			  jsonTestDataFolder?: string & !=""
			  jsonTestDataLiveSearch?: string & !=""
			}
			deliveryCodes: [...string] | *["testCode1", "testCode2"]
		}
		api: {
			enabled: bool | *true
		}
		pagination: defaultPageSize: number | *commerce.pagination.defaultPageSize
	}
}`
}

type routes struct {
	controller *controller.View
}

func (r *routes) Inject(controller *controller.View) {
	r.controller = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("product.view", r.controller.Get)
	h := registry.MustRoute("/product/:marketplacecode/:name.html", `product.view(marketplacecode, name, backurl?="")`)
	h.Normalize("name")
	h = registry.MustRoute("/product/:marketplacecode/:variantcode/:name.html", `product.view(marketplacecode, variantcode, name, backurl?="")`)
	h.Normalize("name")
}

type apiRoutes struct {
	apiController *controller.APIController
}

func (r *apiRoutes) Inject(apiController *controller.APIController) {
	r.apiController = apiController
}

func (r *apiRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/api/v1/products/:marketplacecode", "products.api.get")
	registry.HandleGet("products.api.get", r.apiController.Get)
}
