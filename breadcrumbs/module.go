package breadcrumbs

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module for breadcrumbs
type Module struct{}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	controller *Controller
}

// Inject required dependencies
func (r *routes) Inject(controller *Controller) {
	r.controller = controller
}

// Routes defining the name for the data controller
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleData("breadcrumbs", r.controller.Data)
}
