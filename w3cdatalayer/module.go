package w3cdatalayer

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/application"
	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/interfaces/templatefunctions"
)

type (
	// Module represents our w3cdatalayer module
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "w3cDatalayerService", new(templatefunctions.W3cDatalayerService))
	flamingo.BindEventSubscriber(injector).To(application.EventReceiver{})
}

// CueConfig schema and configuration
func (m *Module) CueConfig() string {
	return `
commerce: w3cDatalayer: {	
	pageInstanceIDPrefix?: string
	pageInstanceIDStage?: string
	pageNamePrefix?: string
	siteName?: string
	defaultCurrency?: string
	version?: string
	hashUserValues: bool | *false
	hashEncoding: string | *"base64url" 
	productMediaBaseUrl?: string
	productMediaThumbnailUrlPrefix?: string
	productMediaUrlPrefix?: string
}
`
}

// FlamingoLegacyConfigAlias mapping
func (m *Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"w3cDatalayer.pageInstanceIDPrefix":           "commerce.w3cDatalayer.pageInstanceIDPrefix",
		"w3cDatalayer.pageInstanceIDStage":            "commerce.w3cDatalayer.pageInstanceIDStage",
		"w3cDatalayer.pageNamePrefix":                 "commerce.w3cDatalayer.pageNamePrefix",
		"w3cDatalayer.siteName":                       "commerce.w3cDatalayer.siteName",
		"w3cDatalayer.defaultCurrency":                "commerce.w3cDatalayer.defaultCurrency",
		"w3cDatalayer.version":                        "commerce.w3cDatalayer.version",
		"w3cDatalayer.hashUserValues":                 "commerce.w3cDatalayer.hashUserValues",
		"w3cDatalayer.hashEncoding":                   "commerce.w3cDatalayer.hashEncoding",
		"w3cDatalayer.productMediaBaseUrl":            "commerce.w3cDatalayer.productMediaBaseUrl",
		"w3cDatalayer.productMediaThumbnailUrlPrefix": "commerce.w3cDatalayer.productMediaThumbnailUrlPrefix",
		"w3cDatalayer.productMediaUrlPrefix":          "commerce.w3cDatalayer.productMediaUrlPrefix",
	}
}
