package graphql

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
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
	types.Map("Commerce_Product", new(graphqlProductDto.Product))
	types.Map("Commerce_Product_SimpleProduct", graphqlProductDto.SimpleProduct{})
	types.Map("Commerce_Product_ConfigurableProduct", graphqlProductDto.ConfigurableProduct{})
	types.Map("Commerce_Product_ActiveVariantProduct", graphqlProductDto.ActiveVariantProduct{})
	types.Map("Commerce_Product_VariationSelection", graphqlProductDto.VariationSelection{})
	types.Map("Commerce_Product_ActiveVariationSelection", graphqlProductDto.ActiveVariationSelection{})
	types.Map("Commerce_Product_VariationSelection_Option", graphqlProductDto.VariationSelectionOption{})
	types.Map("Commerce_Product_VariationSelection_OptionState", new(graphqlProductDto.VariationSelectionOptionState))
	types.Map("Commerce_Product_VariationSelection_OptionVariant", graphqlProductDto.VariationSelectionOptionVariant{})
	types.Map("Commerce_Product_Categories", graphqlProductDto.ProductCategories{})
	types.Map("Commerce_Product_Meta", graphqlProductDto.ProductMeta{})
	types.Map("Commerce_Product_Loyalty", graphqlProductDto.ProductLoyalty{})
	types.Map("Commerce_Product_Loyalty_PriceInfo", domain.LoyaltyPriceInfo{})
	types.Map("Commerce_Product_Loyalty_EarningInfo", domain.LoyaltyEarningInfo{})
	types.Map("Commerce_Product_PriceContext", domain.PriceContext{})
	types.Map("Commerce_Product_Media", graphqlProductDto.ProductMedia{})
	types.Map("Commerce_Product_MediaItem", domain.Media{})
	types.Map("Commerce_Product_Attributes", domain.Attributes{})
	types.GoField("Commerce_Product_Attributes", "getAttribute", "Attribute")
	types.GoField("Commerce_Product_Attributes", "getAttributes", "Attributes")
	types.GoField("Commerce_Product_Attributes", "getAttributeKeys", "AttributeKeys")
	types.GoField("Commerce_Product_Attributes", "getAttributesByKey", "AttributesByKey")
	types.Map("Commerce_Product_Attribute", domain.Attribute{})
	types.Map("Commerce_Product_CategoryTeaser", domain.CategoryTeaser{})
	types.Map("Commerce_Product_PriceInfo", domain.PriceInfo{})
	types.Resolve("Commerce_Product_PriceInfo", "activeBase", CommerceProductQueryResolver{}, "ActiveBase")
	types.Map("Commerce_Product_SearchResult", SearchResultDTO{})
	types.Map("Commerce_Product_Badges", graphqlProductDto.ProductBadges{})
	types.Map("Commerce_Product_Badge", domain.Badge{})

	types.Resolve("Query", "Commerce_Product", CommerceProductQueryResolver{}, "CommerceProduct")
	types.Resolve("Query", "Commerce_Product_Search", CommerceProductQueryResolver{}, "CommerceProductSearch")
}
