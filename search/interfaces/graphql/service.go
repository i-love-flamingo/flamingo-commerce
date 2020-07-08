package graphql

import (
	"flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o graphql.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models return the 'Schema name' => 'Go model' mapping of this module
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Search_Meta":           domain.SearchMeta{},
		"Commerce_Search_Request":        dto.CommerceSearchRequest{},
		"Commerce_Search_KeyValueFilter": dto.CommerceSearchKeyValueFilter{},
		"Commerce_Search_Suggestion":     domain.Suggestion{},
		"Commerce_Search_Result":         application.SearchResult{},
		"Commerce_Search_SortOption":     dto.CommerceSearchSortOption{},
		"Commerce_Search_Facet":          new(dto.Commerce_Search_Facet),
		"Commerce_Search_ListFacet":      dto.Commerce_Search_ListFacet{},
		"Commerce_Search_TreeFacet":      dto.Commerce_Search_TreeFacet{},
		"Commerce_Search_RangeFacet":     dto.Commerce_Search_RangeFacet{},
		"Commerce_Search_FacetItem":      new(dto.Commerce_Search_FacetItem),
		"Commerce_Search_ListFacetItem":  dto.Commerce_Search_ListFacetItem{},
		"Commerce_Search_TreeFacetItem":  dto.Commerce_Search_TreeFacetItem{},
		"Commerce_Search_RangeFacetItem": dto.Commerce_Search_RangeFacetItem{},
	}.Models()
}
