package category

import (
	"go.aoe.com/flamingo/core/category/interfaces/controller"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

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
	m.RouterRegistry.Handle("category.view", new(controller.View))
	m.RouterRegistry.Route("/category/:code/:name.html", "category.view(code, name, *)")
	m.RouterRegistry.Route("/category/:code", "category.view(code, *)")

	m.RouterRegistry.Handle("category.tree", new(controller.Tree))
	m.RouterRegistry.Handle("category", new(controller.Entity))
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.category.view.template": "category/category",
	}
}
