package category

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/category/application"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/infrastructure"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	categoryGraphql "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/search"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	flamingographql "flamingo.me/graphql"
)

// Module registers our profiler
type Module struct {
	useCategoryFixedAdapter bool
}

// URL to category
func URL(code string) (string, map[string]string) {
	return application.URL(code)
}

// URLWithName to category
func URLWithName(code, name string) (string, map[string]string) {
	return application.URLWithName(code, web.URLTitle(name))
}

// Inject dependencies
func (m *Module) Inject(
	routerRegistry *web.RouterRegistry,
	config *struct {
		UseCategoryFixedAdapter bool `inject:"config:commerce.category.useCategoryFixedAdapter,optional"`
	},
) {
	if config != nil {
		m.useCategoryFixedAdapter = config.UseCategoryFixedAdapter
	}
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	if m.useCategoryFixedAdapter {
		injector.Bind((*domain.CategoryService)(nil)).To(infrastructure.CategoryServiceFixed{})

	}
	web.BindRoutes(injector, new(routes))
	injector.Bind(new(application.RouterRouter)).To(new(web.Router))
	injector.BindMulti(new(flamingographql.Service)).To(categoryGraphql.Service{})
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"commerce.category.view.template":       "category/category",
		"commerce.category.view.teaserTemplate": "category/teaser",
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(product.Module),
		new(search.Module),
	}
}

type routes struct {
	view   *controller.ViewController
	entity *controller.Entity
	tree   *controller.Tree
}

func (r *routes) Inject(view *controller.ViewController, entity *controller.Entity, tree *controller.Tree) {
	r.view = view
	r.entity = entity
	r.tree = tree
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("category.view", r.view.Get)
	handler, _ := registry.Route("/category/:code/:name.html", "category.view(code, name, *)")
	handler.Normalize("name")
	registry.Route("/category/:code", "category.view(code, *)")

	registry.HandleData("category.tree", r.tree.Data)
	registry.HandleData("category", r.entity.Data)
}
