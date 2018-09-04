package templatefunctions

import (
	"flamingo.me/flamingo-commerce/product/application"
	"flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// GetProductUrl is exported as a template function
	GetProductUrl struct {
		UrlService *application.UrlService `inject:""`
		Logger     flamingo.Logger         `inject:""`
	}
)

// Func returns the JSON object
func (tf *GetProductUrl) Func() interface{} {
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
