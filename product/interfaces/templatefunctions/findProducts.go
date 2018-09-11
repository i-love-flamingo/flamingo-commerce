package templatefunctions

import (
	"context"
	"log"

	"fmt"

	"flamingo.me/flamingo-commerce/product/application"
	searchApplication "flamingo.me/flamingo-commerce/search/application"
	"flamingo.me/flamingo/core/pugtemplate/pugjs"
)

type (
	// FindProducts is exported as a template function
	FindProducts struct {
		SearchService *application.ProductSearchService `inject:""`
	}
)

func (tf *FindProducts) Func(ctx context.Context) interface{} {

	return func(widgetName string, config interface{}) *application.SearchResult {
		var searchRequest searchApplication.SearchRequest
		if pugjsMap, ok := config.(pugjs.Map); ok {
			configValues := pugjsMap.AsStringMap()
			fmt.Printf("%#v", configValues)
			//TODO - fill all the searchRequest
			filters := make(map[string][]string)
			for k, v := range configValues {
				filters[k] = []string{v}
			}
			searchRequest = searchApplication.SearchRequest{
				FilterBy: filters,
			}
		}
		result, e := tf.SearchService.Find(ctx, &searchRequest)
		if e != nil {
			log.Printf("Error: product.interfaces.templatefunc %v", e)
			return &application.SearchResult{}
		}

		return result
	}
}
