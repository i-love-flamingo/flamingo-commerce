package graphql

import (
	"sort"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/product/application"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

// SearchResultDTOFactory creates SearchResultDTO instances with injected dependencies.
type SearchResultDTOFactory struct {
	logger       flamingo.Logger
	facetMappers []searchdto.FacetMapper
}

// Inject dependencies
func (f *SearchResultDTOFactory) Inject(logger flamingo.Logger, facetMappers []searchdto.FacetMapper) {
	f.logger = logger
	f.facetMappers = facetMappers
}

// NewSearchResultDTO creates a new SearchResultDTO with the factory's dependencies.
func (f *SearchResultDTOFactory) NewSearchResultDTO(res *application.SearchResult) *SearchResultDTO {
	return &SearchResultDTO{
		result:       res,
		logger:       f.logger,
		facetMappers: f.facetMappers,
	}
}

// WrapSearchResult wraps the search result into the graphql dto.
//
// Deprecated: Use SearchResultDTOFactory.NewSearchResultDTO instead to ensure
// proper dependency injection for facet mapping.
func WrapSearchResult(res *application.SearchResult) *SearchResultDTO {
	return &SearchResultDTO{
		result: res,
	}
}

// SearchResultDTO search result dto for graphql
type SearchResultDTO struct {
	result       *application.SearchResult
	logger       flamingo.Logger
	facetMappers []searchdto.FacetMapper
}

// Suggestions get suggestions
func (obj *SearchResultDTO) Suggestions() []searchdomain.Suggestion {
	return obj.result.Suggestions
}

func (obj *SearchResultDTO) Actions() []searchdomain.Action {
	return obj.result.Actions
}

// Products get products
func (obj *SearchResultDTO) Products() []graphqlProductDto.Product {
	products := make([]graphqlProductDto.Product, 0, len(obj.result.Products))
	for _, p := range obj.result.Products {
		products = append(products, graphqlProductDto.NewGraphqlProductDto(p, nil, nil))
	}

	return products
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
func (obj *SearchResultDTO) Facets() []searchdto.CommerceSearchFacet {
	var res = []searchdto.CommerceSearchFacet{}

	for _, facet := range obj.result.Facets {
		mappedFacet := mapFacet(facet, obj.logger, obj.facetMappers)
		if mappedFacet != nil {
			res = append(res, mappedFacet)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Position() < res[j].Position()
	})

	return res
}

// Promotion returns possible promotion data
func (obj *SearchResultDTO) Promotion() *searchdto.PromotionDTO {
	if len(obj.result.Promotions) > 0 {
		return searchdto.WrapPromotion(&obj.result.Promotions[0])
	}

	return nil
}

func mapFacet(facet searchdomain.Facet, logger flamingo.Logger, facetMappers []searchdto.FacetMapper) searchdto.CommerceSearchFacet {
	for _, mapper := range facetMappers {
		if mapped, ok := mapper.MapFacet(facet); ok {
			return mapped
		}
	}

	logger.Warn("Trying to map unknown facet type: ", facet.Type)

	return nil
}

// HasSelectedFacet check if there are any selected facets
func (obj *SearchResultDTO) HasSelectedFacet() bool {
	return len(obj.result.SearchMeta.SelectedFacets) > 0
}
