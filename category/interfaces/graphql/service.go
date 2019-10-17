package graphql

import (
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/category/interfaces/controller"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o schema.go -pkg graphql schema.graphql

// Service describes the Commerce/Category GraphQL Service
type Service struct{}

// Schema for category, delivery and addresses
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models mapping for Commerce_Category types
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Tree":                    new(domain.Tree),
		"Commerce_CategoryTree":            domain.TreeData{},
		"Commerce_Category":                new(domain.Category),
		"Commerce_CategoryData":            domain.CategoryData{},
		"Commerce_Category_SearchResult":   controller.ViewData{},
		"Commerce_Category_Attributes":     domain.Attributes{},
		"Commerce_Category_Attribute":      domain.Attribute{},
		"Commerce_Category_AttributeValue": domain.Attributevalue{},
	}.Models()
}
