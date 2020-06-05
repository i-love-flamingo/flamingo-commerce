package graphql

import (
	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o fs.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models return the 'Schema name' => 'Go model' mapping of this module
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Product": graphql.ModelMapEntry{
			Type: new(domain.BasicProduct),
			Fields: map[string]string{
				"specifications": "GetSpecifications",
			},
		},
		"Commerce_SimpleProduct": graphql.ModelMapEntry{
			Type: domain.SimpleProduct{},
			Fields: map[string]string{
				"specifications": "GetSpecifications",
			},
		},
		"Commerce_ConfigurableProduct": graphql.ModelMapEntry{
			Type: domain.ConfigurableProduct{},
			Fields: map[string]string{
				"specifications": "GetSpecifications",
			},
		},
		"Commerce_Product_Variant":           domain.Variant{},
		"Commerce_BasicProductData":          domain.BasicProductData{},
		"Commerce_ProductTeaserData":         domain.TeaserData{},
		"Commerce_ProductSpecifications":     domain.Specifications{},
		"Commerce_ProductSpecificationGroup": domain.SpecificationGroup{},
		"Commerce_ProductSpecificationEntry": domain.SpecificationEntry{},
		"Commerce_ProductSaleable":           domain.Saleable{},
		"Commerce_ProductMedia":              domain.Media{},
		"Commerce_ProductAttributes": graphql.ModelMapEntry{
			Type: domain.Attributes{},
			Fields: map[string]string{
				"getAttribute":       "Attribute",
				"getAttributes":      "Attributes",
				"getAttributeKeys":   "AttributeKeys",
				"getAttributesByKey": "AttributesByKey",
			},
		},
		"Commerce_ProductAttribute":          domain.Attribute{},
		"Commerce_CategoryTeaser":            domain.CategoryTeaser{},
		"Commerce_ProductPriceInfo":          domain.PriceInfo{},
		"Commerce_ProductLoyaltyPriceInfo":   domain.LoyaltyPriceInfo{},
		"Commerce_ProductLoyaltyEarningInfo": domain.LoyaltyEarningInfo{},
		"Commerce_PriceContext":              domain.PriceContext{},
		"Commerce_Product_SearchResult":      application.SearchResult{},
	}.Models()
}
