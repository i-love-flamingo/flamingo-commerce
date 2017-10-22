package templatefunctions

import (
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// GetProduct is exported as a template function
	GetProduct struct {
		ProductService domain.ProductService `inject:""`
	}
)

// Name alias for use in template
func (tf GetProduct) Name() string {
	return "getProduct"
}

func (tf GetProduct) Func(ctx web.Context) interface{} {
	return func(marketplaceCode string) domain.BasicProduct {
		product, _ := tf.ProductService.Get(ctx, marketplaceCode)
		return product
	}
}
