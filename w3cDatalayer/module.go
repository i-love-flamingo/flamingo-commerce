package w3cDatalayer

import (
	"flamingo.me/flamingo-commerce/w3cDatalayer/application"
	"flamingo.me/flamingo-commerce/w3cDatalayer/interfaces/templatefunctions"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/template"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	template.BindCtxFunc(injector, "w3cDatalayerService", new(templatefunctions.W3cDatalayerService))
	injector.BindMulti((*event.SubscriberWithContext)(nil)).To(application.EventReceiver{})
}
