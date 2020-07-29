package graphql

import (
	"flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/searchdto"
	"flamingo.me/graphql"
)

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -nometadata -o graphql.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

var _ graphql.Service = new(Service)

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Search_Meta", domain.SearchMeta{})
	types.Resolve("Commerce_Search_Meta", "sortOptions", CommerceSearchQueryResolver{}, "SortOptions")
	types.Map("Commerce_Search_Request", searchdto.CommerceSearchRequest{})
	types.Map("Commerce_Search_KeyValueFilter", searchdto.CommerceSearchKeyValueFilter{})
	types.Map("Commerce_Search_Suggestion", domain.Suggestion{})
	types.Map("Commerce_Search_Result", application.SearchResult{})
	types.Map("Commerce_Search_SortOption", searchdto.CommerceSearchSortOption{})
	types.Map("Commerce_Search_Facet", new(searchdto.CommerceSearchFacet))
	types.Map("Commerce_Search_ListFacet", searchdto.CommerceSearchListFacet{})
	types.Map("Commerce_Search_TreeFacet", searchdto.CommerceSearchTreeFacet{})
	types.Map("Commerce_Search_RangeFacet", searchdto.CommerceSearchRangeFacet{})
	types.Map("Commerce_Search_FacetItem", new(searchdto.CommerceSearchFacetItem))
	types.Map("Commerce_Search_ListFacetItem", searchdto.CommerceSearchListFacetItem{})
	types.Map("Commerce_Search_TreeFacetItem", searchdto.CommerceSearchTreeFacetItem{})
	types.Map("Commerce_Search_RangeFacetItem", searchdto.CommerceSearchRangeFacetItem{})
}
