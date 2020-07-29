package graphql

import (
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/graphql/categorydto"
	"flamingo.me/graphql"
)

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -nometadata -o schema.go -pkg graphql schema.graphql

// Service describes the Commerce/Category GraphQL Service
type Service struct{}

var _ graphql.Service = new(Service)

// Schema for category, delivery and addresses
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Tree", new(domain.Tree))
	types.Map("Commerce_CategoryTree", domain.TreeData{})
	types.Map("Commerce_Category", new(domain.Category))
	types.Map("Commerce_CategoryData", domain.CategoryData{})
	types.Map("Commerce_Category_SearchResult", categorydto.CategorySearchResult{})
	types.Map("Commerce_Category_Attributes", domain.Attributes{})
	types.Map("Commerce_Category_Attribute", domain.Attribute{})
	types.Map("Commerce_Category_AttributeValue", domain.AttributeValue{})

	types.Resolve("Query", "Commerce_CategoryTree", CommerceCategoryQueryResolver{}, "CommerceCategoryTree")
	types.Resolve("Query", "Commerce_Category", CommerceCategoryQueryResolver{}, "CommerceCategory")
}
