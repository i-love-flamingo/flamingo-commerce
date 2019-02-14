package templatefunctions

import (
	"context"
	"log"

	"strconv"

	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/pugtemplate/pugjs"
)

type (
	// FindProducts is exported as a template function
	FindProducts struct {
		SearchService *application.ProductSearchService `inject:""`
	}
)

// Func defines the find products function
func (tf *FindProducts) Func(ctx context.Context) interface{} {

	/*
		widgetName - used to namespace widget - in case we need pagination
		config - map with certain keys - used to specifiy th searchRequest better
	*/
	return func(widgetName string, searchConfig interface{}, additionalFilters interface{}) *application.SearchResult {
		var searchRequest searchApplication.SearchRequest
		//fmt.Printf("%#v", searchConfig)

		if pugjsMap, ok := searchConfig.(*pugjs.Map); ok {
			searchConfigValues := pugjsMap.AsStringMap()
			//fmt.Printf("%#v", searchConfigValues)

			searchRequest = searchApplication.SearchRequest{
				SortDirection: searchConfigValues["sortDirection"],
				SortBy:        searchConfigValues["sortBy"],
				Query:         searchConfigValues["query"],
			}
			pageSize, err := strconv.Atoi(searchConfigValues["pageSize"])
			if err == nil {
				searchRequest.PageSize = pageSize
			}
		}

		searchRequest.FilterBy = asFilterMap(additionalFilters)
		//fmt.Printf("%#v", searchRequest)
		result, e := tf.SearchService.Find(ctx, &searchRequest)
		if e != nil {
			log.Printf("Error: product.interfaces.templatefunc %v", e)
			return &application.SearchResult{}
		}

		return result
	}
}

func asFilterMap(additionalFilters interface{}) map[string]interface{} {
	filters := make(map[string]interface{})
	// use filtersPug as KeyValueFilter
	if filtersPug, ok := additionalFilters.(*pugjs.Map); ok {
		for k, v := range filtersPug.Items {
			if v, ok := v.(*pugjs.Array); ok {
				var filterList []string
				for _, item := range v.Items() {
					filterList = append(filterList, item.String())
				}
				filters[k.String()] = filterList
			}
			if v, ok := v.(pugjs.String); ok {
				filters[k.String()] = []string{v.String()}
			}
		}
	}
	return filters
}
