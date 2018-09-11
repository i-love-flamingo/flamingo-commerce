package templatefunctions

import (
	"context"
	"log"
	"flamingo.me/flamingo/core/pugtemplate/pugjs"

	productDomain "flamingo.me/flamingo-commerce/product/domain"
	searchDomain "flamingo.me/flamingo-commerce/search/domain"
)

type (
	// FindProducts is exported as a template function
	FindProducts struct {
		SearchService productDomain.SearchService `inject:""`
	}
)

func (tf *FindProducts) Func(ctx context.Context) interface{} {
	return func(filtersPug pugjs.Map, sortPug pugjs.Map) []productDomain.BasicProduct {
		var filter []searchDomain.Filter

		// use filtersPug as KeyValueFilter
		for k, v := range filtersPug.Items {
			if v, ok := v.(*pugjs.Array); ok {
				var filterList []string
				for _, item:= range v.Items() {
					filterList = append(filterList, item.String())
				}
				filter = append(filter, searchDomain.NewKeyValueFilter(k.String(), filterList))
			}
			if v, ok := v.(pugjs.String); ok {
				filter = append(filter, searchDomain.NewKeyValueFilter(k.String(), []string{v.String()}))
			}
		}

		// use sortPug as SortFilter
		for k, v := range sortPug.Items {
			if v, ok := v.(pugjs.String); ok {
				filter = append(filter, searchDomain.NewSortFilter(k.String(), v.String()))
			}
		}

		products, e := tf.SearchService.Search(ctx, filter...)
		if e != nil {
			log.Printf("Error: product.interfaces.templatefunc %v", e)
			return nil
		}

		return products.Hits
	}
}