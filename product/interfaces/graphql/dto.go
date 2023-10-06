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

// Promotion returns possible promotion data
func (obj *SearchResultDTO) Promotion() *searchdto.PromotionDTO {
	if len(obj.result.Promotions) > 0 {
		return searchdto.WrapPromotion(&obj.result.Promotions[0])
	}

	return nil
}

func mapFacet(facet searchdomain.Facet, logger flamingo.Logger) searchdto.CommerceSearchFacet {
	switch searchdomain.FacetType(facet.Type) {
	case searchdomain.ListFacet:
		return searchdto.WrapListFacet(facet)

	case searchdomain.TreeFacet:
		return searchdto.WrapTreeFacet(facet)

	case searchdomain.RangeFacet:
		return searchdto.WrapRangeFacet(facet)

	default:
		logger.Warn("Trying to map unknown facet type: ", facet.Type)
		return nil
	}
}

// HasSelectedFacet check if there are any selected facets
func (obj *SearchResultDTO) HasSelectedFacet() bool {
	return len(obj.result.SearchMeta.SelectedFacets) > 0
}
