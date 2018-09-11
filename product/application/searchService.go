package application

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo-commerce/search/application"
	searchdomain "flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (

	// ProductSearchService
	ProductSearchService struct {
		SearchService         domain.SearchService         `inject:""`
		PaginationInfoFactory *utils.PaginationInfoFactory `inject:""`
		DefaultPageSize       float64                      `inject:"pagination.defaultPageSize,optional"`
		Logger                flamingo.Logger              `inject:""`
	}

	//SearchResult - much like the corresponding struct in search package, just that instead "Hits" we have a list of matching Products
	SearchResult struct {
		Suggestions    []searchdomain.Suggestion
		Products       []domain.BasicProduct
		SearchMeta     searchdomain.SearchMeta
		Facets         searchdomain.FacetCollection
		PaginationInfo utils.PaginationInfo
	}
)

//Find return SearchResult with matched products - based on given input
func (s *ProductSearchService) Find(ctx context.Context, searchRequest *application.SearchRequest) (*SearchResult, error) {
	var currentUrl *url.URL
	request, found := web.FromContext(ctx)
	if !found {
		currentUrl = nil
	} else {
		currentUrl = request.Request().URL
	}

	if searchRequest == nil {
		searchRequest = &application.SearchRequest{}
	}
	//pageSize can either be set in the request, or we use the configured default or if nothing set we rely on the SearchService later
	pageSize := searchRequest.PageSize
	if pageSize == 0 {
		pageSize = int(s.DefaultPageSize)
	}

	result, err := s.SearchService.Search(ctx, application.BuildFilters(*searchRequest, pageSize)...)
	if err != nil {
		return nil, err
	}

	if searchRequest.PaginationConfig == nil {
		searchRequest.PaginationConfig = s.PaginationInfoFactory.DefaultConfig
	}

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
	}, *searchRequest.PaginationConfig, currentUrl)

	return &SearchResult{
		SearchMeta:     result.SearchMeta,
		Facets:         result.Facets,
		Suggestions:    result.Suggestion,
		Products:       result.Hits,
		PaginationInfo: paginationInfo,
	}, nil
}
