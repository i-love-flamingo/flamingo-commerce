package search

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/search/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module registers our search package
type Module struct{}

// Configure the search URL
func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	controller *interfaces.ViewController
}

func (r *routes) Inject(controller *interfaces.ViewController) {
	r.controller = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("search.search", r.controller.Get)
	registry.Route("/search/:type", `search.search(type, *)`)
	registry.Route("/search", `search.search`)
}
