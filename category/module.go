package category

import (
	"flamingo/core/category/interfaces/controller"
	"flamingo/framework/config"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
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

// URL with name to category
func URLWithName(code, name string) (string, map[string]string) {
	return controller.URLWithName(code, name)
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("category.view", new(controller.View))
	m.RouterRegistry.Route("/category/:code/:name.html", "category.view")
	m.RouterRegistry.Route("/category/:code", "category.view")
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.category.view.template": "category/category",
	}
}
