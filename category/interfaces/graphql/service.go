package graphql

import (
	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o schema.go -pkg graphql schema.graphql

// Service describes the Commerce/Category GraphQL Service
type Service struct{}

// Schema for category, delivery and addresses
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models mapping for Commerce_Cart types
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Tree":         new(domain.Tree),
		"Commerce_CategoryTree": domain.TreeData{},
	}.Models()
}
