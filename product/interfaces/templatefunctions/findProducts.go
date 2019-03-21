package templatefunctions

import (
	"context"
	"log"
	"strconv"
	"strings"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/search/domain"

	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
)

// FindProducts is exported as a template function
type (
	FindProducts struct {
		ProductSearchService *application.ProductSearchService `inject:""`
	}

	filterProcessing struct {
		buildSearchRequest searchApplication.SearchRequest
		whiteList          []string
		blackList          []string
	}
)

// Func defines the find products function
func (tf *FindProducts) Func(ctx context.Context) interface{} {

	/*
		widgetName - used to namespace widget - in case we need pagination
		config - map with certain keys - used to specifiy th searchRequest better
	*/
	return func(namespace string, configs ...map[string]string) *application.SearchResult {
		var searchConfig, keyValueFilters, filterConstrains map[string]string

		if len(configs) > 0 {
			searchConfig = configs[0]
		} else {
			searchConfig = make(map[string]string)
		}

		if len(configs) > 1 {
			keyValueFilters = configs[1]
		} else {
			keyValueFilters = make(map[string]string)
		}

		if len(configs) > 2 {
			filterConstrains = configs[2]
		} else {
			filterConstrains = make(map[string]string)
		}

		filterProcessing := newFilterProcessing(web.RequestFromContext(ctx), namespace, searchConfig, keyValueFilters, filterConstrains)

		//searchRequest.FilterBy = asFilterMap(keyValueFilters)
		//fmt.Printf("%#v", searchRequest)
		result, e := tf.ProductSearchService.Find(ctx, &filterProcessing.buildSearchRequest)

		if e != nil {
			log.Printf("Error: product.interfaces.templatefunc %v", e)
			return &application.SearchResult{}
		}
		return filterProcessing.modifyResult(result)
	}
}

func newFilterProcessing(request *web.Request, namespace string, searchConfig, keyValueFilters, filterConstrains map[string]string) filterProcessing {
	var filterProcessing filterProcessing
	var searchRequest searchApplication.SearchRequest

	// 1- set the originalSearchRequest from given searchConfig and keyValueFilters
	searchRequest = searchApplication.SearchRequest{
		SortDirection: searchConfig["sortDirection"],
		SortBy:        searchConfig["sortBy"],
		Query:         searchConfig["query"],
	}
	pageSize, err := strconv.Atoi(searchConfig["pageSize"])
	if err == nil {
		searchRequest.PageSize = pageSize
	}

	for k, v := range keyValueFilters {
		searchRequest.AddAdditionalFilter(domain.NewKeyValueFilter(k, []string{v}))
	}

	// Set blackList and whiteList
	filterProcessing.blackList = strings.Split(filterConstrains["blackList"], ",")
	filterProcessing.whiteList = strings.Split(filterConstrains["whiteList"], ",")

	//2 - Use the url parameters to modify the filters:
	for k, v := range request.QueryAll() {
		if !strings.HasPrefix(k, namespace) {
			continue
		}
		splitted := strings.SplitN(k, ".", 2)
		if len(splitted) < 2 {
			continue
		}
		filterKey := splitted[1]
		if filterProcessing.isAllowed(filterKey) {
			searchRequest.SetAdditionalFilter(domain.NewKeyValueFilter(filterKey, v))
		}
	}
	filterProcessing.buildSearchRequest = searchRequest
	return filterProcessing
}

func (f filterProcessing) modifyResult(result *application.SearchResult) *application.SearchResult {
	var newFacetCollection domain.FacetCollection
	newFacetCollection = make(map[string]domain.Facet)
	for k, facet := range result.Facets {
		if f.isAllowed(k) {
			newFacetCollection[k] = facet
		}
	}
	result.Facets = newFacetCollection
	return result
}

//isAllowed - checks the given key against the defined whitelist and blacklist (whitelist prefered)
func (f filterProcessing) isAllowed(key string) bool {
	if len(f.whiteList) > 0 {
		for _, wl := range f.whiteList {
			if wl == key {
				return true
			}
		}
		return false
	} else {
		for _, wl := range f.blackList {
			if wl == key {
				return false
			}
		}
	}
	return true
}
