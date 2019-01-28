package breadcrumbs

import (
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module for breadcrumbs
type Module struct{}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, new(routes))
}

type routes struct {
	controller *Controller
}

// Inject required dependencies
func (r *routes) Inject(controller *Controller) {
	r.controller = controller
}

// Routes defining the name for the data controller
func (r *routes) Routes(registry *router.Registry) {
	registry.HandleData("breadcrumbs", r.controller.Data)
}
