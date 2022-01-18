package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	graphqlDto "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql/categorydto"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
)

// CommerceCategoryQueryResolver resolves graphql category queries
type CommerceCategoryQueryResolver struct {
	categoryService domain.CategoryService
	searchService   *productApplication.ProductSearchService
}

// Inject dependencies
func (r *CommerceCategoryQueryResolver) Inject(
	service domain.CategoryService,
	searchService *productApplication.ProductSearchService,
) *CommerceCategoryQueryResolver {
	r.categoryService = service
	r.searchService = searchService
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
	request *searchdto.CommerceSearchRequest) (*graphqlDto.CategorySearchResult, error) {
	category, err := r.categoryService.Get(ctx, categoryCode)

	if err != nil {
		return nil, err
	}

	var filters []searchDomain.Filter
	filters = append(filters, domain.NewCategoryFacet(categoryCode))
	searchRequest := &application.SearchRequest{AdditionalFilter: filters}

	if request != nil {
		for _, filter := range request.KeyValueFilters {
			filters = append(filters, searchDomain.NewKeyValueFilter(filter.K, filter.V))
		}
		searchRequest = &application.SearchRequest{
			AdditionalFilter: filters,
			PageSize:         request.PageSize,
			Page:             request.Page,
			SortBy:           request.SortBy,
			Query:            request.Query,
			PaginationConfig: nil,
		}
	}
	result, err := r.searchService.Find(ctx, searchRequest)

	if err != nil {
		return nil, err
	}

	return &graphqlDto.CategorySearchResult{Category: category, ProductSearchResult: graphql.WrapSearchResult(result)}, nil
}
