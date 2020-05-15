package price

import (
	"flamingo.me/dingo"
	pricegraphql "flamingo.me/flamingo-commerce/v3/price/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/price/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/graphql"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m Module) Configure(injector *dingo.Injector) {
	flamingo.BindTemplateFunc(injector, "commercePriceFormat", new(templatefunctions.CommercePriceFormatFunc))
	injector.BindMulti(new(graphql.Service)).To(pricegraphql.Service{})
}

// Depends adds our dependencies
func (*Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(locale.Module),
	}
}
