package fake

import (
	"context"
	"strconv"

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

// Search returns Products based on given Filters
func (s *SearchService) Search(ctx context.Context, filters ...searchDomain.Filter) (*domain.SearchResult, error) {
	hits := s.findProducts(ctx, filters)
	currentPage := s.findCurrentPage(filters)
	facets, selectedFacets := s.createFacets(filters)

	documents := make([]searchDomain.Document, len(hits))
	for i, hit := range hits {
		documents[i] = hit
	}

	return &domain.SearchResult{
		Result: searchDomain.Result{
			SearchMeta: searchDomain.SearchMeta{
				Query:          "",
				OriginalQuery:  "",
				Page:           currentPage,
				NumPages:       10,
				NumResults:     len(hits),
				SelectedFacets: selectedFacets,
				SortOptions:    nil,
			},
			Hits:       documents,
			Suggestion: []searchDomain.Suggestion{},
			Facets:     facets,
		},
		Hits: hits,
	}, nil
}

// SearchBy returns Products prefiltered by the given attribute (also based on additional given Filters)
func (s *SearchService) SearchBy(ctx context.Context, _ string, _ []string, _ ...searchDomain.Filter) (*domain.SearchResult, error) {
	return s.Search(ctx)
}

func (s *SearchService) findProducts(ctx context.Context, filters []searchDomain.Filter) []domain.BasicProduct {
	products := make([]domain.BasicProduct, 0)

	// - try finding product by marketPlaceCode given in query
	if query, found := s.filterValue(filters, "q"); found {
		if len(query) > 0 {
			product, _ := s.productService.Get(ctx, query[0])
			if product != nil {
				products = append(products, product)
			}
		}
	}

	return products
}

func (s *SearchService) findCurrentPage(filters []searchDomain.Filter) int {
	currentPage := 1

	if page, found := s.filterValue(filters, "page"); found {
		if page, err := strconv.Atoi(page[0]); err == nil {
			currentPage = page
		}
	}

	return currentPage
}

func (s *SearchService) createFacets(filters []searchDomain.Filter) (map[string]searchDomain.Facet, []searchDomain.Facet) {
	selectedFacets := make([]searchDomain.Facet, 0)

	facets := map[string]searchDomain.Facet{
		"brandCode": {
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
		},
	}

	if s.hasFilterWithValue(filters, "brandCode", "apple") != false {
		facets["brandCode"].Items[0].Active = true
		facets["brandCode"].Items[0].Selected = true
		selectedFacets = append(selectedFacets, facets["brandCode"])
	}

	if s.hasFilterWithValue(filters, "retailerCode", "retailer") != false {
		facets["retailerCode"].Items[0].Active = true
		facets["retailerCode"].Items[0].Selected = true
		selectedFacets = append(selectedFacets, facets["retailerCode"])
	}

	return facets, selectedFacets
}

func (s *SearchService) filterValue(filters []searchDomain.Filter, key string) ([]string, bool) {
	for _, filter := range filters {
		filterKey, filterValues := filter.Value()
		if filterKey == key {
			return filterValues, true
		}
	}

	return []string{}, false
}

func (s *SearchService) hasFilterWithValue(filters []searchDomain.Filter, key string, value string) bool {
	if filterValues, found := s.filterValue(filters, key); found {
		for _, filterValue := range filterValues {
			if value == filterValue {
				return true
			}
		}
	}

	return false
}
