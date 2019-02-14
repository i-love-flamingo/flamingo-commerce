package category

import (
	"flamingo.me/flamingo-commerce/v3/category/application"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module registers our profiler
type Module struct{}

// URL to category
func URL(code string) (string, map[string]string) {
	return application.URL(code)
}

// URLWithName to category
func URLWithName(code, name string) (string, map[string]string) {
	return application.URLWithName(code, web.URLTitle(name))
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
	injector.Bind(new(application.RouterRouter)).To(new(router.Router))
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"category.view.template":       "category/category",
		"category.view.teaserTemplate": "category/teaser",
	}
}

type routes struct {
	view   *controller.View
	entity *controller.Entity
	tree   *controller.Tree
}

func (r *routes) Inject(view *controller.View, entity *controller.Entity, tree *controller.Tree) {
	r.view = view
	r.entity = entity
	r.tree = tree
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("category.view", r.view.Get)
	registry.Route("/category/:code/:name.html", "category.view(code, name, *)").Normalize("name")
	registry.Route("/category/:code", "category.view(code, *)")

	registry.HandleData("category.tree", r.tree.Data)
	registry.HandleData("category", r.entity.Data)
}
