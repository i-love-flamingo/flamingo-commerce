package graphql

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/interfaces"
	graphqlDto "flamingo.me/flamingo-commerce/v3/category/interfaces/graphql/categorydto"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	productGraphql "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql"
	"flamingo.me/flamingo-commerce/v3/search/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
)

// CommerceCategoryQueryResolver resolves graphql category queries
type CommerceCategoryQueryResolver struct {
	categoryService     domain.CategoryService
	searchService       *productApplication.ProductSearchService
	searchResultFactory *productGraphql.SearchResultDTOFactory
	logger              flamingo.Logger
}

// Inject dependencies
func (r *CommerceCategoryQueryResolver) Inject(
	service domain.CategoryService,
	searchService *productApplication.ProductSearchService,
	searchResultFactory *productGraphql.SearchResultDTOFactory,
	logger flamingo.Logger,
) *CommerceCategoryQueryResolver {
	r.categoryService = service
	r.searchService = searchService
	r.searchResultFactory = searchResultFactory
	r.logger = logger.WithField(flamingo.LogKeyModule, "category").WithField(flamingo.LogKeyCategory, "graphql")
	return r
}

// CommerceCategoryTree returns a Tree with the given activeCategoryCode from categoryService
func (r *CommerceCategoryQueryResolver) CommerceCategoryTree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	tree, err := r.categoryService.Tree(ctx, activeCategoryCode)
	if err != nil {
		r.logger.Error("Failed to get category tree", err)
		return nil, interfaces.ErrCategoryGeneral
	}

	return tree, nil
}

// CommerceCategory returns product search result with the given categoryCode from searchService
func (r *CommerceCategoryQueryResolver) CommerceCategory(
	ctx context.Context,
	categoryCode string,
	request *searchdto.CommerceSearchRequest) (*graphqlDto.CategorySearchResult, error) {
	category, err := r.categoryService.Get(ctx, categoryCode)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, interfaces.ErrCategoryNotFound
		}

		r.logger.Error("Failed to get category", err)
		return nil, interfaces.ErrCategoryGeneral
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
		r.logger.Error("Failed to search products for category", err)
		return nil, interfaces.ErrCategoryGeneral
	}

	return &graphqlDto.CategorySearchResult{Category: category, ProductSearchResult: r.searchResultFactory.NewSearchResultDTO(result)}, nil
}
