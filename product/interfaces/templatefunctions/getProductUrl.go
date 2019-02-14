package templatefunctions

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// GetProductUrl is exported as a template function
	GetProductUrl struct {
		UrlService *application.UrlService `inject:""`
		Logger     flamingo.Logger         `inject:""`
	}
)

// Func returns the JSON object
func (tf *GetProductUrl) Func(context.Context) interface{} {
	return func(p domain.BasicProduct) string {
		if p == nil {
			tf.Logger.WithField("category", "product").Warn("Called getPrpductUrl templatefunc without a product")
			return ""
		}
		url, err := tf.UrlService.Get(p, "")
		if err != nil {
			tf.Logger.WithField("category", "product").Error(err)
			return ""
		}
		return url
	}
}
