package fake

import (
	"context"

	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

// SearchService is just mocking stuff
type SearchService struct {
	productService *ProductService
}

// Inject dependencies
func (s *SearchService) Inject(
	productService *ProductService,
) *SearchService {
	s.productService = productService

	return s
}

//Search returns Products based on given Filters
func (s *SearchService) Search(ctx context.Context, filter ...searchDomain.Filter) (*domain.SearchResult, error) {
	p1, err := s.productService.Get(ctx, "fake_configurable")
	if err != nil {
		return nil, err
	}
	p2, err := s.productService.Get(ctx, "fake_simple")
	if err != nil {
		return nil, err
	}

	return &domain.SearchResult{
		Result: searchDomain.Result{
			SearchMeta: searchDomain.SearchMeta{
				Query:          "",
				OriginalQuery:  "",
				Page:           1,
				NumPages:       1,
				NumResults:     2,
				SelectedFacets: nil,
				SortOptions:    nil,
			},
			Hits:       []searchDomain.Document{p1, p2},
			Suggestion: nil,
			Facets:     nil,
		},
		Hits: []domain.BasicProduct{p1, p2},
	}, nil

}

// SearchBy returns Products prefiltered by the given attribute (also based on additional given Filters)
func (s *SearchService) SearchBy(ctx context.Context, attribute string, values []string, filter ...searchDomain.Filter) (*domain.SearchResult, error) {
	return s.Search(ctx)
}
