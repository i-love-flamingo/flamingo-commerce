package application

import (
	"context"
	"fmt"
	"net/url"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/search/application"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

type (
	// ProductSearchService - Application service that offers a more explicit way to search for  product results - on top of the domain.ProductSearchService
	ProductSearchService struct {
		searchService         domain.SearchService
		paginationInfoFactory *utils.PaginationInfoFactory
		defaultPageSize       float64
		logger                flamingo.Logger
	}

	// SearchResult - much like the corresponding struct in search package, just that instead "Hits" we have a list of matching Products
	SearchResult struct {
		Suggestions    []searchdomain.Suggestion
		Products       []domain.BasicProduct
		SearchMeta     searchdomain.SearchMeta
		Facets         searchdomain.FacetCollection
		PaginationInfo utils.PaginationInfo
		Promotions     []searchdomain.Promotion
		Actions        []searchdomain.Action
	}
)

// Inject dependencies
func (s *ProductSearchService) Inject(
	searchService domain.SearchService,
	paginationInfoFactory *utils.PaginationInfoFactory,
	logger flamingo.Logger,
	cfg *struct {
		DefaultPageSize float64 `inject:"config:commerce.product.pagination.defaultPageSize,optional"`
	},
) *ProductSearchService {
	s.searchService = searchService
	s.paginationInfoFactory = paginationInfoFactory
	s.logger = logger.
		WithField(flamingo.LogKeyModule, "product").
		WithField(flamingo.LogKeyCategory, "application.ProductSearchService")
	if cfg != nil {
		s.defaultPageSize = cfg.DefaultPageSize
	}

	return s
}

// Find return SearchResult with matched products - based on given input
func (s *ProductSearchService) Find(ctx context.Context, searchRequest *application.SearchRequest) (*SearchResult, error) {
	var currentURL *url.URL
	request := web.RequestFromContext(ctx)
	if request == nil {
		currentURL = nil
	} else {
		currentURL = request.Request().URL
	}

	if searchRequest == nil {
		searchRequest = &application.SearchRequest{}
	}
	// pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the ProductSearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.defaultPageSize)
	}

	result, err := s.searchService.Search(ctx, application.BuildFilters(*searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.paginationInfoFactory.DefaultConfig
	}

	if pageSize != 0 {
		if err := result.SearchMeta.ValidatePageSize(pageSize); err != nil {
			err = fmt.Errorf("the Searchservice seems to ignore pageSize filter, %w", err)
			s.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "Find").Warn(err)
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
		Products:       result.Hits,
		PaginationInfo: paginationInfo,
		Promotions:     result.Promotions,
		Actions:        result.Actions,
	}, nil
}

// FindBy return SearchResult with matched products filtered by the given attribute - based on given input
func (s *ProductSearchService) FindBy(ctx context.Context, attributeCode string, values []string, searchRequest *application.SearchRequest) (*SearchResult, error) {
	var currentURL *url.URL
	request := web.RequestFromContext(ctx)
	if request == nil {
		currentURL = nil
	} else {
		currentURL = request.Request().URL
	}

	if searchRequest == nil {
		searchRequest = &application.SearchRequest{}
	}
	// pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the ProductSearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.defaultPageSize)
	}

	result, err := s.searchService.SearchBy(ctx, attributeCode, values, application.BuildFilters(*searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.paginationInfoFactory.DefaultConfig
	}

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
		Products:       result.Hits,
		PaginationInfo: paginationInfo,
		Promotions:     result.Result.Promotions,
		Actions:        result.Actions,
	}, nil
}
