package price

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/price/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "commercePriceFormat", new(templatefunctions.CommercePriceFormatFunc))
}