package w3cdatalayer

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/application"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/interfaces/templatefunctions"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "w3cDatalayerService", new(templatefunctions.W3cDatalayerService))
	flamingo.BindEventSubscriber(injector).To(application.EventReceiver{})
}

// CueConfig schema and configuration
func (m *Module) CueConfig() string {
	return `w3cDatalayer: {	hashEncoding: string | *"base64url" }`
}
