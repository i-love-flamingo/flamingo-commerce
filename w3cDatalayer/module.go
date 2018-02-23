package w3cDatalayer

import (
	"go.aoe.com/flamingo/core/w3cDatalayer/application"
	"go.aoe.com/flamingo/core/w3cDatalayer/interfaces/templatefunctions"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/template"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*template.ContextFunction)(nil)).To(templatefunctions.W3cDatalayerService{})
	injector.BindMulti((*event.SubscriberWithContext)(nil)).To(application.EventReceiver{})
}
