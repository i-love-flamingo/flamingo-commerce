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

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("category.view", new(controller.View))
	m.RouterRegistry.Route("/category/:categorycode/:name.html", "category.view")
	m.RouterRegistry.Route("/category/:categorycode", "category.view")
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"core.category.view.template": "category/category",
	}
}
