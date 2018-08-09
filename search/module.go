package search

import (
	"flamingo.me/flamingo-commerce/search/interfaces"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module registers our search package
type Module struct{}

// Configure the search URL
func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, new(routes))
}

type routes struct {
	controller *interfaces.ViewController
}

func (r *routes) Inject(controller *interfaces.ViewController) {
	r.controller = controller
}

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleGet("search.search", r.controller.Get)
	registry.Route("/search/:type", `search.search(type, *)`)
	registry.Route("/search", `search.search`)
}
