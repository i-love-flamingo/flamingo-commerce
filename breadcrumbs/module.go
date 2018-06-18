package breadcrumbs

import (
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

// Module for breadcrumbs
type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

// Configure module
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("breadcrumbs", new(Controller))
}
