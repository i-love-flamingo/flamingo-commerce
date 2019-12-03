package graphql

import (
	"context"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
)

// CommerceProductQueryResolver resolves graphql product queries
type CommerceProductQueryResolver struct {
	productService domain.ProductService
	searchService  *productApplication.ProductSearchService
}

// Inject dependencies
func (r *CommerceProductQueryResolver) Inject(
	productService domain.ProductService,
	searchService *productApplication.ProductSearchService,
) *CommerceProductQueryResolver {
	r.productService = productService
	r.searchService = searchService
	return r
}

// CommerceProduct returns a product with the given marketplaceCode from productService
func (r *CommerceProductQueryResolver) CommerceProduct(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	return r.productService.Get(ctx, marketplaceCode)
}

// CommerceProductSearch returns a search result of products based on the given search request
func (r *CommerceProductQueryResolver) CommerceProductSearch(ctx context.Context, request *dto.CommerceSearchRequest) (*productApplication.SearchResult, error) {

	var filters []searchDomain.Filter
	for _, filter := range request.KeyValueFilters {
		filters = append(filters, searchDomain.NewKeyValueFilter(filter.K, filter.V))
	}

	result, err := r.searchService.Find(ctx, &application.SearchRequest{
		AdditionalFilter: filters,
		PageSize:         request.PageSize,
		Page:             request.Page,
		SortBy:           request.SortBy,
		SortDirection:    request.SortDirection,
		Query:            request.Query,
		PaginationConfig: nil,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
