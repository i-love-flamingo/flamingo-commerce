package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	graphqlDto "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql/dto"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
)

// CommerceCategoryQueryResolver resolves graphql category queries
type CommerceCategoryQueryResolver struct {
	categoryService      domain.CategoryService
	productSearchService productDomain.SearchService
	defaultPageSize      float64
}

// Inject dependencies
func (r *CommerceCategoryQueryResolver) Inject(
	service domain.CategoryService,
	productSearchService productDomain.SearchService,
	optionals *struct {
		DefaultPageSize float64 `inject:"config:pagination.defaultPageSize,optional"`
	},
) *CommerceCategoryQueryResolver {
	r.categoryService = service
	r.productSearchService = productSearchService
	if optionals != nil {
		r.defaultPageSize = optionals.DefaultPageSize
	}
	return r
}

// CommerceCategoryTree returns a Tree with the given activeCategoryCode from categoryService
func (r *CommerceCategoryQueryResolver) CommerceCategoryTree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	return r.categoryService.Tree(ctx, activeCategoryCode)
}

// CommerceCategory returns product search result with the given categoryCode from searchService
func (r *CommerceCategoryQueryResolver) CommerceCategory(
	ctx context.Context,
	categoryCode string,
	request *dto.CommerceSearchRequest) (*graphqlDto.CategorySearchResult, error) {
	category, err := r.categoryService.Get(ctx, categoryCode)

	if err != nil {
		return nil, err
	}

	var filters = dto.SearchRequestToFilters(request, int(r.defaultPageSize))

	filters = append(filters, domain.NewCategoryFacet(categoryCode))

	result, err := r.productSearchService.Search(ctx, filters...)

	if err != nil {
		return nil, err
	}

	productSearchResult := &productApplication.SearchResult{
		Suggestions: result.Suggestion,
		Products:    result.Hits,
		SearchMeta:  result.SearchMeta,
		Facets:      result.Facets,
	}

	return &graphqlDto.CategorySearchResult{Category: category, ProductSearchResult: productSearchResult}, nil
}
