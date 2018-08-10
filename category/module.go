package category

import (
	"flamingo.me/flamingo-commerce/category/interfaces/controller"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

// Module registers our profiler
type Module struct{}

// URL to category
func URL(code string) (string, map[string]string) {
	return controller.URL(code)
}

// URLWithName to category
func URLWithName(code, name string) (string, map[string]string) {
	return controller.URLWithName(code, web.URLTitle(name))
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, new(routes))
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.category.view.template":       "category/category",
		"core.category.view.teaserTemplate": "category/teaser",
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

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleGet("category.view", r.view.Get)
	registry.Route("/category/:code/:name.html", "category.view(code, name, *)").Normalize("name")
	registry.Route("/category/:code", "category.view(code, *)")

	registry.HandleData("category.tree", r.tree.Data)
	registry.HandleData("category", r.entity.Data)
}
