package graphql

import (
	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
	"flamingo.me/flamingo-commerce/v3/search/utils"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"sort"
)

// WrapSearchResult wraps the search result into the graphql dto
func WrapSearchResult(res *application.SearchResult) *SearchResultDTO {
	return &SearchResultDTO{
		result: res,
	}
}

// SearchResultDTO search result dto for graphql
type SearchResultDTO struct {
	result *application.SearchResult
	logger flamingo.Logger
}

// Inject dependencies
func (obj *SearchResultDTO) Inject(logger flamingo.Logger) {
	obj.logger = logger
}

// Suggestions get suggestions
func (obj *SearchResultDTO) Suggestions() []searchdomain.Suggestion {
	return obj.result.Suggestions
}

// Products get products
func (obj *SearchResultDTO) Products() []domain.BasicProduct {
	return obj.result.Products
}

// SearchMeta get search meta
func (obj *SearchResultDTO) SearchMeta() searchdomain.SearchMeta {
	return obj.result.SearchMeta
}

// PaginationInfo get pagination info
func (obj *SearchResultDTO) PaginationInfo() utils.PaginationInfo {
	return obj.result.PaginationInfo
}

// Facets get facets
func (obj *SearchResultDTO) Facets() []dto.CommerceSearchFacet {
	var res = []dto.CommerceSearchFacet{}

	for _, facet := range obj.result.Facets {
		mappedFacet := mapFacet(facet, obj.logger)
		if mappedFacet != nil {
			res = append(res, mappedFacet)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Position() < res[j].Position()
	})

	return res
}

func mapFacet(facet searchdomain.Facet, logger flamingo.Logger) dto.CommerceSearchFacet {
	switch searchdomain.FacetType(facet.Type) {
	case searchdomain.ListFacet:
		return dto.WrapListFacet(facet)

	case searchdomain.TreeFacet:
		return dto.WrapTreeFacet(facet)

	case searchdomain.RangeFacet:
		return dto.WrapRangeFacet(facet)

	default:
		logger.Warn("Trying to map unknown facet type: ", facet.Type)
		return nil
	}
}

// HasSelectedFacet check if there are any selected facets
func (obj *SearchResultDTO) HasSelectedFacet() bool {
	return len(obj.result.SearchMeta.SelectedFacets) > 0
}
