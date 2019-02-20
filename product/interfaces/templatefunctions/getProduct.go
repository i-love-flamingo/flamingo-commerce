package templatefunctions

import (
	"context"
	"log"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// GetProduct is exported as a template function
	GetProduct struct {
		ProductService domain.ProductService `inject:""`
	}
)

// Func factory
func (tf *GetProduct) Func(ctx context.Context) interface{} {
	return func(marketplaceCode string) domain.BasicProduct {
		product, e := tf.ProductService.Get(ctx, marketplaceCode)
		if e != nil {
			log.Printf("Error: product.interfaces.templatefunc %v", e)
		}
		return product
	}
}
