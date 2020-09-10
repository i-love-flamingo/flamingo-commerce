package graphql

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
	productDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
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
	types.Map("Commerce_Product", new(productDto.Product))
	types.Map("Commerce_SimpleProduct", productDto.SimpleProduct{})
	types.Map("Commerce_ConfigurableProduct", productDto.ConfigurableProduct{})
	types.Map("Commerce_ProductCategories", productDto.Categories{})
	types.Map("Commerce_ProductMeta", productDto.Meta{})
	types.Map("Commerce_ProductLoyalty", productDto.Loyalty{})
	//types.Map("Commerce_BasicProductData", domain.BasicProductData{})
	//types.Map("Commerce_ProductTeaserData", domain.TeaserData{})
	types.Map("Commerce_Product_Variant", domain.Variant{})

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
