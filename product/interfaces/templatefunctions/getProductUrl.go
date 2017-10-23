package templatefunctions

import (
	"go.aoe.com/flamingo/core/product/interfaces/controller"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// GetProductUrl is exported as a template function
	GetProductUrl struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (tf GetProductUrl) Name() string {
	return "getProductUrl"
}

// Func returns the JSON object
func (tf GetProductUrl) Func() interface{} {
	return func(p *pugjs.Map) string {
		sku := p.Field("baseData").Field("marketPlaceCode").String()
		return tf.Router.URL(controller.URL(sku, sku)).String()
	}
}
