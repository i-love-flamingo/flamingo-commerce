package graphql

import (
	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/product/application"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
)

// MapFacetForTest exports mapFacet for use in black-box tests.
func MapFacetForTest(facet searchdomain.Facet, facetMappers []searchdto.FacetMapper, logger flamingo.Logger) searchdto.CommerceSearchFacet {
	return mapFacet(facet, facetMappers, logger)
}

// NewSearchResultDTOForTest creates a SearchResultDTO for use in black-box tests.
func NewSearchResultDTOForTest(result *application.SearchResult, logger flamingo.Logger, facetMappers []searchdto.FacetMapper) *SearchResultDTO {
	return &SearchResultDTO{
		result:       result,
		logger:       logger,
		facetMappers: facetMappers,
	}
}
