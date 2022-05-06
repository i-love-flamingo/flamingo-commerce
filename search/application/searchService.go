/*
Package application in the search module provides a more explicit access to the domain searchservice, and combines the result with Paginations etc.
The structs defined here can for example be used in the interface layers
*/
package application

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// SearchService - Application service that offers a more explicit way to search for results - on top of the domain.ProductSearchService
	SearchService struct {
		searchService         domain.SearchService
		paginationInfoFactory *utils.PaginationInfoFactory
		defaultPageSize       float64
		logger                flamingo.Logger
	}

	// SearchRequest is a simple DTO for the search query data
	SearchRequest struct {
		AdditionalFilter []domain.Filter
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

// Inject dependencies
func (s *SearchService) Inject(
	paginationInfoFactory *utils.PaginationInfoFactory,
	logger flamingo.Logger,
	optionals *struct {
		SearchService   domain.SearchService `inject:",optional"`
		DefaultPageSize float64              `inject:"config:commerce.pagination.defaultPageSize,optional"`
	}) *SearchService {
	s.paginationInfoFactory = paginationInfoFactory
	s.logger = logger.
		WithField(flamingo.LogKeyModule, "search").
		WithField(flamingo.LogKeyCategory, "application.ProductSearchService")
	if optionals != nil {
		s.searchService = optionals.SearchService
		s.defaultPageSize = optionals.DefaultPageSize
	}
	return s
}

// FindBy returns a SearchResult for the given Request
func (s *SearchService) FindBy(ctx context.Context, documentType string, searchRequest SearchRequest) (*SearchResult, error) {
	if s.searchService == nil {
		return nil, errors.New("no searchservice available")
	}
	var currentURL *url.URL
	request := web.RequestFromContext(ctx)
	if request == nil {
		currentURL = nil
	} else {
		currentURL = request.Request().URL
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.paginationInfoFactory.DefaultConfig
	}

	// pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the ProductSearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.defaultPageSize)
	}

	result, err := s.searchService.SearchFor(ctx, documentType, BuildFilters(searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	// do a logical pageSize check - and log warning
	//  10 pageSize * (3 pages* -1 ) + lastPageSize = 35 results*
	if pageSize != 0 {
		if err := result.SearchMeta.ValidatePageSize(pageSize); err != nil {
			err = fmt.Errorf("the Searchservice seems to ignore pageSize filter, %w", err)
			s.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "FindBy").Warn(err)
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
	if s.searchService == nil {
		return nil, errors.New("no searchservice available")
	}
	var currentURL *url.URL
	request := web.RequestFromContext(ctx)
	if request == nil {
		currentURL = nil
	} else {
		currentURL = request.Request().URL
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.paginationInfoFactory.DefaultConfig
	}

	// pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the ProductSearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.defaultPageSize)
	}

	result, err := s.searchService.Search(ctx, BuildFilters(searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	// do a logical pageSize check - and log warning
	//  10 pageSize * (3 pages* -1 ) + lastPageSize = 35 results*
	if pageSize != 0 {
		for k, r := range result {
			if err := r.SearchMeta.ValidatePageSize(pageSize); err != nil {
				err = fmt.Errorf("the Searchservice seems to ignore pageSize filter for document type %q, %w", k, err)
				s.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "Find").Warn(err)
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

	filters = append(filters, request.AdditionalFilter...)

	return filters
}

// AddAdditionalFilter adds an additional filter
func (r *SearchRequest) AddAdditionalFilter(filter domain.Filter) {
	r.AdditionalFilter = append(r.AdditionalFilter, filter)
}

// SetAdditionalFilter - adds or replaces the given filter
func (r *SearchRequest) SetAdditionalFilter(filterToAdd domain.Filter) {
	for k, existingFilter := range r.AdditionalFilter {
		existingFilterKey, _ := existingFilter.Value()
		filterToAddKey, _ := filterToAdd.Value()
		if existingFilterKey == filterToAddKey {
			r.AdditionalFilter[k] = filterToAdd
			return
		}
	}
	r.AddAdditionalFilter(filterToAdd)
}

// AddAdditionalFilters adds multiple additional filters
func (r *SearchRequest) AddAdditionalFilters(filters ...domain.Filter) {
	for _, filter := range filters {
		r.AddAdditionalFilter(filter)
	}
}
