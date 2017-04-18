package product

import (
	"flamingo/framework/dingo"
	"flamingo/core/product/controller"
	"flamingo/framework/router"
)

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("product.view", new(controller.ViewController))
	m.RouterRegistry.Route("/product/{uid}", "product.view")
}
