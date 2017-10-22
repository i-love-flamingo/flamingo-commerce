package templatefunctions

import (
	"net/url"

	"go.aoe.com/flamingo/core/product/interfaces/controller"
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
	return func(productdata interface{}) *url.URL {
		return tf.Router.URL(controller.URL("jjj", "kkk"))
	}
}
