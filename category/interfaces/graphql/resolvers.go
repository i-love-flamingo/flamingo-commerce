package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	searchGraphQlDto "flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
)

// CommerceCategoryQueryResolver resolves graphql category queries
type CommerceCategoryQueryResolver struct {
	categoryService      domain.CategoryService
	productSearchService *productApplication.ProductSearchService
}

// Inject dependencies
func (r *CommerceCategoryQueryResolver) Inject(service domain.CategoryService, productSearchService *productApplication.ProductSearchService) {
	r.categoryService = service
	r.productSearchService = productSearchService
}

// CommerceCategoryTree returns a Tree with the given activeCategoryCode from categoryService
func (r *CommerceCategoryQueryResolver) CommerceCategoryTree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	return r.categoryService.Tree(ctx, activeCategoryCode)
}

// CommerceCategory returns product search result with the given categoryCode from searchService
func (r *CommerceCategoryQueryResolver) CommerceCategory(
	ctx context.Context,
	categoryCode string,
	categorySearchRequest *searchGraphQlDto.CommerceSearchRequest) (*controller.ViewData, error) {

	category, err := r.categoryService.Get(ctx, categoryCode)

	if err != nil {
		return &controller.ViewData{Category: category, ProductSearchResult: nil}, err
	}

	searchRequest := new(searchApplication.SearchRequest)

	if categorySearchRequest != nil {
		for _, filter := range categorySearchRequest.KeyValueFilters {
			searchRequest.AdditionalFilter = append(searchRequest.AdditionalFilter, searchDomain.NewKeyValueFilter(filter.K, filter.V))
		}
	}

	// - use categoryDomain.NewCategoryFacet as filter to use product/category endpoint from searchperience
	searchRequest.SetAdditionalFilter(domain.NewCategoryFacet(categoryCode))
	result, err := r.productSearchService.Find(ctx, searchRequest)

	if err != nil {
		return &controller.ViewData{Category: category, ProductSearchResult: nil}, err
	}

	return &controller.ViewData{Category: category, ProductSearchResult: result}, err
}
