package templatefunctions

import (
	"flamingo.me/flamingo-commerce/w3cDatalayer/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	W3cDatalayerService struct {
		applicationServiceProvider application.ServiceProvider
	}
)

// Inject dependencies
func (w3cdl *W3cDatalayerService) Inject(provider application.ServiceProvider) {
	w3cdl.applicationServiceProvider = provider
}

// Name alias for use in template
func (w3cdl W3cDatalayerService) Name() string {
	return "w3cDatalayerService"
}

// Func template function factory
func (w3cdl W3cDatalayerService) Func(ctx web.Context) interface{} {
	// Usage
	// w3cDatalayerService().get()
	return func() *application.Service {
		service := w3cdl.applicationServiceProvider()
		service.Init(ctx)
		return service
	}
}
