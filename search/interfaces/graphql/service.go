package graphql

import (
	"flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
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
		"Commerce_Search_Request":        searchdto.CommerceSearchRequest{},
		"Commerce_Search_KeyValueFilter": searchdto.CommerceSearchKeyValueFilter{},
		"Commerce_Search_Suggestion":     domain.Suggestion{},
		"Commerce_Search_Result":         application.SearchResult{},
		"Commerce_Search_SortOption":     searchdto.CommerceSearchSortOption{},
		"Commerce_Search_Facet":          new(searchdto.CommerceSearchFacet),
		"Commerce_Search_ListFacet":      searchdto.CommerceSearchListFacet{},
		"Commerce_Search_TreeFacet":      searchdto.CommerceSearchTreeFacet{},
		"Commerce_Search_RangeFacet":     searchdto.CommerceSearchRangeFacet{},
		"Commerce_Search_FacetItem":      new(searchdto.CommerceSearchFacetItem),
		"Commerce_Search_ListFacetItem":  searchdto.CommerceSearchListFacetItem{},
		"Commerce_Search_TreeFacetItem":  searchdto.CommerceSearchTreeFacetItem{},
		"Commerce_Search_RangeFacetItem": searchdto.CommerceSearchRangeFacetItem{},
	}.Models()
}
