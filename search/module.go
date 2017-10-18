package search

import (
	"go.aoe.com/flamingo/core/search/interfaces"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// Module registers our search package
	Module struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure the search URL
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("search.search", new(interfaces.ViewController))
	m.RouterRegistry.Route("/search", `search.search(type="product")`)
	m.RouterRegistry.Route("/search/:type", `search.search(type)`)
	m.RouterRegistry.Route("/search", `search.search`)
}
