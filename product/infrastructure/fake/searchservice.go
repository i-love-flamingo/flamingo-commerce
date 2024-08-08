package fake

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"flamingo.me/flamingo/v3/framework/config"

	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// SearchService is just mocking stuff
	SearchService struct {
		productService         *ProductService
		liveSearchJSON         string
		categoryFacetItemsJSON string
		sortConfig             []sortConfig
	}
	liveSearchData struct {
		Marketplacecodes []string                  `json:"marketplacecodes"`
		Sugestions       []searchDomain.Suggestion `json:"sugestions"`
		Promotions       []searchDomain.Promotion  `json:"promotions"`
		Actions          []searchDomain.Action     `json:"actions"`
	}

	// sortConfig contains sorting configuration for search results
	sortConfig struct {
		Key   string
		Label string
		Asc   string
		Desc  string
	}
)

// Inject dependencies
func (s *SearchService) Inject(
	productService *ProductService,
	cfg *struct {
		LiveSearchJSON         string       `inject:"config:commerce.product.fakeservice.jsonTestDataLiveSearch,optional"`
		CategoryFacetItemsJSON string       `inject:"config:commerce.product.fakeservice.jsonTestDataCategoryFacetItems,optional"`
		SortConfig             config.Slice `inject:"config:commerce.product.fakeservice.sorting"`
	},
) *SearchService {
	s.productService = productService
	if cfg != nil {
		s.liveSearchJSON = cfg.LiveSearchJSON
		s.categoryFacetItemsJSON = cfg.CategoryFacetItemsJSON

		err := cfg.SortConfig.MapInto(&s.sortConfig)
		if err != nil {
			panic(err)
		}
	}

	return s
}

// Search returns Products based on given Filters
func (s *SearchService) Search(ctx context.Context, filters ...searchDomain.Filter) (*domain.SearchResult, error) {
	hits := s.findProducts(ctx, filters, s.productService.GetMarketPlaceCodes())
	currentPage := s.findCurrentPage(filters)
	facets, selectedFacets := s.createFacets(filters)

	var promotions []searchDomain.Promotion
	if len(hits) > 0 {
		promotions = append(promotions, searchDomain.Promotion{
			Title:                "Promotion",
			Content:              "",
			URL:                  "http://promotion.example.com/",
			Media:                []searchDomain.Media{{}},
			AdditionalAttributes: nil,
		})
	}

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
				SortOptions:    mapSortOptions(s.sortConfig, filters),
			},
			Hits:       documents,
			Suggestion: []searchDomain.Suggestion{},
			Facets:     facets,
			Promotions: promotions,
		},
		Hits: hits,
	}, nil
}

func (s *SearchService) livesearch(ctx context.Context, query string) (*domain.SearchResult, error) {
	data := make(map[string]liveSearchData)

	fileContent, err := os.ReadFile(s.liveSearchJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return nil, err
	}

	liveSearchData := data[query]

	hits := s.findProducts(ctx, nil, liveSearchData.Marketplacecodes)
	documents := make([]searchDomain.Document, len(hits))
	for i, hit := range hits {
		documents[i] = hit
	}

	return &domain.SearchResult{
		Result: searchDomain.Result{
			SearchMeta: searchDomain.SearchMeta{
				Page:       1,
				NumPages:   1,
				NumResults: len(hits),
			},
			Hits:       documents,
			Suggestion: liveSearchData.Sugestions,
			Promotions: liveSearchData.Promotions,
			Actions:    liveSearchData.Actions,
		},
		Hits: hits,
	}, nil
}

// SearchBy returns Products prefiltered by the given attribute (also based on additional given Filters)
func (s *SearchService) SearchBy(ctx context.Context, attribute string, _ []string, filters ...searchDomain.Filter) (*domain.SearchResult, error) {
	if attribute == "livesearch" && s.liveSearchJSON != "" {
		var query string
		for _, f := range filters {
			if qf, ok := f.(*searchDomain.QueryFilter); ok {
				_, q := qf.Value()
				query = q[0]
				break
			}
		}
		return s.livesearch(ctx, query)
	}

	return s.Search(ctx, filters...)
}

func (s *SearchService) findProducts(ctx context.Context, filters []searchDomain.Filter, marketPlaceCodes []string) []domain.BasicProduct {
	products := make([]domain.BasicProduct, 0)

	// - try finding product by marketPlaceCode given in query or return nothing if query is no-results
	if query, found := s.filterValue(filters, "q"); found {
		if len(query) > 0 {
			if query[0] == "no-results" {
				return products
			}
			product, _ := s.productService.Get(ctx, query[0])
			if product != nil {
				products = append(products, product)
			}
		}
	}

	// - get default products
	if len(products) == 0 {
		for _, marketPlaceCode := range marketPlaceCodes {
			product, _ := s.productService.Get(ctx, marketPlaceCode)
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

// mapSortOptions maps configured sort options to the ones potentially provided via filter
func mapSortOptions(sortConfigs []sortConfig, filters []searchDomain.Filter) []searchDomain.SortOption {
	lookup := make(map[string]bool, 1) // only one field expected

	for _, filter := range filters {
		if sortFiler, ok := filter.(*searchDomain.SortFilter); ok {
			lookup[sortFiler.Field()] = true // direction always true for that filter name
		}
	}

	result := make([]searchDomain.SortOption, len(sortConfigs))

	for i, sortConfig := range sortConfigs {
		result[i] = searchDomain.SortOption{
			Label:        sortConfig.Label,
			Field:        sortConfig.Key,
			Asc:          sortConfig.Asc,
			Desc:         sortConfig.Desc,
			SelectedAsc:  lookup[sortConfig.Asc],
			SelectedDesc: lookup[sortConfig.Desc],
		}
	}

	// if no filters come, we should decide on default, so let it be first in config
	if len(lookup) == 0 && len(result) > 0 {
		result[0].SelectedDesc = true
	}

	return result
}

func (s *SearchService) createFacets(filters []searchDomain.Filter) (map[string]searchDomain.Facet, []searchDomain.Facet) {
	selectedFacets := make([]searchDomain.Facet, 0)

	categoryFilterValue := s.categoryFilterValue(filters)

	facets := map[string]searchDomain.Facet{
		"brandCode": {
			Type:  searchDomain.ListFacet,
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
			Type:  searchDomain.ListFacet,
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

		"categoryCodes": s.createCategoryFacet(categoryFilterValue),
	}

	if s.hasFilterWithValue(filters, "brandCode", "apple") {
		facets["brandCode"].Items[0].Active = true
		facets["brandCode"].Items[0].Selected = true
		selectedFacets = append(selectedFacets, facets["brandCode"])
	}

	if s.hasFilterWithValue(filters, "retailerCode", "retailer") {
		facets["retailerCode"].Items[0].Active = true
		facets["retailerCode"].Items[0].Selected = true
		selectedFacets = append(selectedFacets, facets["retailerCode"])
	}

	if categoryFilterValue != "" {
		selectedFacets = append(selectedFacets, facets["categoryCodes"])
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

func (s *SearchService) categoryFilterValue(filters []searchDomain.Filter) string {
	for _, filter := range filters {
		filterKey, filterValues := filter.Value()
		switch filterKey {
		case "category", "categoryCodes":
			if len(filterValues) > 0 {
				return filterValues[0]
			}
		}
	}

	return ""
}

func (s *SearchService) createCategoryFacet(selectedCategory string) searchDomain.Facet {
	return searchDomain.Facet{
		Type:     searchDomain.TreeFacet,
		Name:     "categoryCodes",
		Label:    "Category",
		Items:    s.createCategoryFacetItems(selectedCategory),
		Position: 0,
	}
}

func (s *SearchService) createCategoryFacetItems(selectedCategory string) []*searchDomain.FacetItem {
	items, err := loadCategoryFacetItems(s.categoryFacetItemsJSON)

	if err != nil {
		return nil
	}

	selectFacetItems(selectedCategory, items)
	return items
}

func selectFacetItems(selectedCategory string, items []*searchDomain.FacetItem) bool {
	for _, item := range items {
		if item.Value == selectedCategory {
			item.Active = true
			item.Selected = true
			return true
		}
		childSelectedOrActive := selectFacetItems(selectedCategory, item.Items)
		if childSelectedOrActive {
			item.Active = true
			return true
		}
	}
	return false
}
