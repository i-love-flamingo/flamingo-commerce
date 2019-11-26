package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
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
	var filters = dto.SearchRequestToFilters(request, int(r.defaultPageSize))
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
