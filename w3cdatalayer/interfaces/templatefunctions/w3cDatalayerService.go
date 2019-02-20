package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/w3cdatalayer/application"
)

type (
	// W3cDatalayerService struct for the template function
	W3cDatalayerService struct {
		applicationServiceProvider application.ServiceProvider
	}
)

// Inject dependencies
func (w3cdl *W3cDatalayerService) Inject(provider application.ServiceProvider) {
	w3cdl.applicationServiceProvider = provider
}

// Func template function factory
func (w3cdl *W3cDatalayerService) Func(ctx context.Context) interface{} {
	// Usage
	// w3cDatalayerService().get()
	return func() *application.Service {
		service := w3cdl.applicationServiceProvider()
		service.Init(ctx)
		return service
	}
}
