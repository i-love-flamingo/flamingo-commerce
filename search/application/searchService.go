/*
Package application in the search module provides a more explicit access to the domain searchservice, and combines the result with Paginations etc.
The structs defined here can for example be used in the interface layers
*/
package application

import (
	"context"
	"net/url"

	categoryDomain "flamingo.me/flamingo-commerce/category/domain"
	"flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (
	// SearchService - Application service that offers a more explicit way to search for results - on top of the domain.SearchService
	SearchService struct {
		SearchService         domain.SearchService         `inject:""`
		PaginationInfoFactory *utils.PaginationInfoFactory `inject:""`
		DefaultPageSize       float64                      `inject:"config:pagination.defaultPageSize,optional"`
		Logger                flamingo.Logger              `inject:""`
	}

	// SearchRequest is a simple DTO for the search query data
	SearchRequest struct {
		FilterBy         map[string]interface{}
		PageSize         int
		Page             int
		SortBy           string
		SortDirection    string
		Query            string
		PaginationConfig *utils.PaginationConfig
	}

	// SearchResult is the DTO for the search result
	SearchResult struct {
		Hits           []domain.Document
		SearchMeta     domain.SearchMeta
		Facets         domain.FacetCollection
		Suggestions    []domain.Suggestion
		PaginationInfo utils.PaginationInfo
	}
)

// FindBy returns a SearchResult for the given Request
func (s *SearchService) FindBy(ctx context.Context, documentType string, searchRequest SearchRequest) (*SearchResult, error) {
	var currentURL *url.URL
	request, found := web.FromContext(ctx)
	if !found {
		currentURL = nil
	} else {
		currentURL = request.Request().URL
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.PaginationInfoFactory.DefaultConfig
	}

	// pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the SearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.DefaultPageSize)
	}

	result, err := s.SearchService.SearchFor(ctx, documentType, BuildFilters(searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	// do a logical pageSize check - and log warning
	//  10 pageSize * (3 pages* -1 ) + lastPageSize = 35 results*
	if pageSize != 0 {
		if err := result.SearchMeta.ValidatePageSize(pageSize); err != nil {
			s.Logger.WithField("category", "application.SearchService").Warn("The Searchservice seems to ignore pageSize Filter")
		}
	}

	paginationInfo := utils.BuildWith(utils.CurrentResultInfos{
		LastPage:   result.SearchMeta.NumPages,
		TotalHits:  result.SearchMeta.NumResults,
		PageSize:   searchRequest.PageSize,
		ActivePage: result.SearchMeta.Page,
	}, *searchRequest.PaginationConfig, currentURL)

	return &SearchResult{
		SearchMeta:     result.SearchMeta,
		Facets:         result.Facets,
		Suggestions:    result.Suggestion,
		Hits:           result.Hits,
		PaginationInfo: paginationInfo,
	}, nil
}

// Find returns a Searchresult for all document types for the given Request
func (s *SearchService) Find(ctx context.Context, searchRequest SearchRequest) (map[string]*SearchResult, error) {
	var currentURL *url.URL
	request, found := web.FromContext(ctx)
	if !found {
		currentURL = nil
	} else {
		currentURL = request.Request().URL
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.PaginationInfoFactory.DefaultConfig
	}

	// pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the SearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.DefaultPageSize)
	}

	result, err := s.SearchService.Search(ctx, BuildFilters(searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	// do a logical pageSize check - and log warning
	//  10 pageSize * (3 pages* -1 ) + lastPageSize = 35 results*
	if pageSize != 0 {
		for k, r := range result {
			if err := r.SearchMeta.ValidatePageSize(pageSize); err != nil {
				s.Logger.WithField("category", "application.SearchService").Warn("The Searchservice seems to ignore pageSize Filter for document type ", k)
			}
		}
	}

	searchResult := make(map[string]*SearchResult)

	for k, r := range result {
		paginationInfo := utils.BuildWith(utils.CurrentResultInfos{
			LastPage:   r.SearchMeta.NumPages,
			TotalHits:  r.SearchMeta.NumResults,
			PageSize:   searchRequest.PageSize,
			ActivePage: r.SearchMeta.Page,
		}, *searchRequest.PaginationConfig, currentURL)

		searchResult[k] = &SearchResult{
			SearchMeta:     r.SearchMeta,
			Facets:         r.Facets,
			Suggestions:    r.Suggestion,
			Hits:           r.Hits,
			PaginationInfo: paginationInfo,
		}
	}

	return searchResult, nil
}

// BuildFilters creates a slice of search filters from the request data
func BuildFilters(request SearchRequest, defaultPageSize int) []domain.Filter {
	var filters []domain.Filter
	if request.Query != "" {
		filters = append(filters, domain.NewQueryFilter(request.Query))
	}

	if request.Page != 0 {
		filters = append(filters, domain.NewPaginationPageFilter(request.Page))
	}

	if request.PageSize != 0 {
		filters = append(filters, domain.NewPaginationPageSizeFilter(request.PageSize))
	} else if defaultPageSize != 0 {
		filters = append(filters, domain.NewPaginationPageSizeFilter(defaultPageSize))
	}

	if request.SortBy != "" {
		filters = append(filters, domain.NewSortFilter(request.SortBy, request.SortDirection))
	}

	for k, v := range request.FilterBy {
		switch k {
		case string(categoryDomain.CategoryKey):
			if _, ok := v.(categoryDomain.CategoryFacet); ok {
				filters = append(filters, v.(categoryDomain.CategoryFacet))
			} else {
				filters = append(
					filters,
					categoryDomain.NewCategoryFacet(
						categoryDomain.CategoryData{
							CategoryCode: v.([]string)[0],
						},
					),
				)
			}
		default:
			filters = append(filters, domain.NewKeyValueFilter(k, v.([]string)))
		}
	}

	return filters
}
