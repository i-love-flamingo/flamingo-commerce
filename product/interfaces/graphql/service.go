package graphql

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/graphql"
)

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -nometadata -o fs.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

var _ graphql.Service = new(Service)

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Product", new(domain.BasicProduct))
	types.GoField("Commerce_Product", "specifications", "GetSpecifications")
	types.Map("Commerce_SimpleProduct", domain.SimpleProduct{})
	types.GoField("Commerce_SimpleProduct", "specifications", "GetSpecifications")
	types.Map("Commerce_ConfigurableProduct", domain.ConfigurableProduct{})
	types.GoField("Commerce_ConfigurableProduct", "specifications", "GetSpecifications")
	types.Map("Commerce_Product_Variant", domain.Variant{})
	types.Map("Commerce_BasicProductData", domain.BasicProductData{})
	types.Map("Commerce_ProductTeaserData", domain.TeaserData{})
	types.Map("Commerce_ProductSpecifications", domain.Specifications{})
	types.Map("Commerce_ProductSpecificationGroup", domain.SpecificationGroup{})
	types.Map("Commerce_ProductSpecificationEntry", domain.SpecificationEntry{})
	types.Map("Commerce_ProductSaleable", domain.Saleable{})
	types.Map("Commerce_ProductMedia", domain.Media{})
	types.Map("Commerce_ProductAttributes", domain.Attributes{})
	types.GoField("Commerce_ProductAttributes", "getAttribute", "Attribute")
	types.GoField("Commerce_ProductAttributes", "getAttributes", "Attributes")
	types.GoField("Commerce_ProductAttributes", "getAttributeKeys", "AttributeKeys")
	types.GoField("Commerce_ProductAttributes", "getAttributesByKey", "AttributesByKey")
	types.Map("Commerce_ProductAttribute", domain.Attribute{})
	types.Map("Commerce_CategoryTeaser", domain.CategoryTeaser{})
	types.Map("Commerce_ProductPriceInfo", domain.PriceInfo{})
	types.Map("Commerce_ProductLoyaltyPriceInfo", domain.LoyaltyPriceInfo{})
	types.Map("Commerce_ProductLoyaltyEarningInfo", domain.LoyaltyEarningInfo{})
	types.Map("Commerce_PriceContext", domain.PriceContext{})
	types.Map("Commerce_Product_SearchResult", SearchResultDTO{})

	types.Resolve("Query", "Commerce_Product", CommerceProductQueryResolver{}, "CommerceProduct")
	types.Resolve("Query", "Commerce_Product_Search", CommerceProductQueryResolver{}, "CommerceProductSearch")
}
