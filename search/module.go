package search

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/search/interfaces"
	searchgraphql "flamingo.me/flamingo-commerce/v3/search/interfaces/graphql"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/graphql"
)

// Module registers our search package
type Module struct{}

// Configure the search URL
func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))

	injector.BindMulti(new(graphql.Service)).To(new(searchgraphql.Service))
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
