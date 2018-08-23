package templatefunctions

import (
	"context"
	"log"

	productDomain "flamingo.me/flamingo-commerce/product/domain"
)

type (
	// FindProducts is exported as a template function
	FindProducts struct {
		SearchService productDomain.SearchService `inject:""`
	}
)

func (tf *FindProducts) Func(ctx context.Context) interface{} {
	return func() []productDomain.BasicProduct {
		products, e := tf.SearchService.Search(ctx, nil)
		if e != nil {
			log.Printf("Error: product.interfaces.templatefunc %v", e)
			return nil
		}
		return products.Hits
	}
}
