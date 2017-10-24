package templatefunctions

import (
	"go.aoe.com/flamingo/core/product/interfaces/controller"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/core/product/domain"
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
	return func(p domain.BasicProduct) string {
		if p == nil {
			return "-"
		}
		return tf.Router.URL(controller.URL(p.BaseData().MarketPlaceCode, p.BaseData().Title)).String()
	}
}
