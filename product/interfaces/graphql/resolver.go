package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	applicationSearchService "flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
)

// CommerceProductQueryResolver resolves graphql product queries
type CommerceProductQueryResolver struct {
	productService       domain.ProductService
	productSearchService domain.SearchService
	defaultPageSize      float64
}

// Inject dependencies
func (r *CommerceProductQueryResolver) Inject(
	productService domain.ProductService,
	productSearchService domain.SearchService,
	optionals *struct {
		DefaultPageSize float64 `inject:"config:pagination.defaultPageSize,optional"`
	},
) *CommerceProductQueryResolver {
	r.productService = productService
	r.productSearchService = productSearchService
	if optionals != nil {
		r.defaultPageSize = optionals.DefaultPageSize
	}
	return r
}

// CommerceProduct returns a product with the given marketplaceCode from productService
func (r *CommerceProductQueryResolver) CommerceProduct(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	return r.productService.Get(ctx, marketplaceCode)
}

// CommerceProductSearch returns a search result of products based on the given search request
func (r *CommerceProductQueryResolver) CommerceProductSearch(ctx context.Context, request *dto.CommerceSearchRequest) (*application.SearchResult, error) {
	var filters = r.searchRequestToFilters(request)

	result, err := r.productSearchService.Search(ctx, filters...)

	if err != nil {
		return nil, err
	}

	return &application.SearchResult{
		Suggestions: result.Suggestion,
		Products:    result.Hits,
		SearchMeta:  result.SearchMeta,
		Facets:      result.Facets,
	}, nil
}

// searchRequestToFilters maps CommerceSearchRequest to Filter
func (r *CommerceProductQueryResolver) searchRequestToFilters(searchRequest *dto.CommerceSearchRequest) []searchDomain.Filter {
	var filters []searchDomain.Filter

	if searchRequest != nil {
		filters = applicationSearchService.BuildFilters(applicationSearchService.SearchRequest{
			AdditionalFilter: nil,
			PageSize:         searchRequest.PageSize,
			Page:             searchRequest.Page,
			SortBy:           searchRequest.SortBy,
			SortDirection:    searchRequest.SortDirection,
			Query:            searchRequest.Query,
			PaginationConfig: nil,
		}, int(r.defaultPageSize))

		for _, filter := range searchRequest.KeyValueFilters {
			filters = append(filters, searchDomain.NewKeyValueFilter(filter.K, filter.V))
		}
	}

	return filters
}
