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
	documents := make([]searchDomain.Document, 0)
	hits := make([]domain.BasicProduct, 0)

	for _, marketPlaceCode := range s.productService.GetMarketPlaceCodes() {
		product, _ := s.productService.Get(ctx, marketPlaceCode)
		if product == nil {
			continue
		}
		documents = append(documents, product)
		hits = append(hits, product)
	}

	return &domain.SearchResult{
		Result: searchDomain.Result{
			SearchMeta: searchDomain.SearchMeta{
				Query:          "",
				OriginalQuery:  "",
				Page:           1,
				NumPages:       1,
				NumResults:     len(hits),
				SelectedFacets: nil,
				SortOptions:    nil,
			},
			Hits:       documents,
			Suggestion: []searchDomain.Suggestion{},
			Facets: map[string]searchDomain.Facet{"brandCode": {
				Type:  string(searchDomain.ListFacet),
				Name:  "brandCode",
				Label: "Brand",
				Items: []*searchDomain.FacetItem{{
					Label:    "Apple",
					Value:    "apple",
					Active:   false,
					Selected: false,
					Count:    2,
				}},
				Position: 0,
			},
				"retailerCode": {
					Type:  string(searchDomain.ListFacet),
					Name:  "retailerCode",
					Label: "Retailer",
					Items: []*searchDomain.FacetItem{{
						Label:    "Test Retailer",
						Value:    "retailer",
						Active:   false,
						Selected: false,
						Count:    2,
					}},
					Position: 0,
				}},
		},
		Hits: hits,
	}, nil

}

// SearchBy returns Products prefiltered by the given attribute (also based on additional given Filters)
func (s *SearchService) SearchBy(ctx context.Context, attribute string, values []string, filter ...searchDomain.Filter) (*domain.SearchResult, error) {
	return s.Search(ctx)
}
