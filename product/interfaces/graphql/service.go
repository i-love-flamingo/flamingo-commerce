package graphql

import (
	// embed schema.graphql
	_ "embed"

	"flamingo.me/graphql"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	graphqlProductDto "flamingo.me/flamingo-commerce/v3/product/interfaces/graphql/product/dto"
)

// Service is the Graphql-Service of this module
type Service struct{}

var _ graphql.Service = new(Service)

//go:embed schema.graphql
var schema []byte

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return schema
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Product", new(graphqlProductDto.Product))
	types.Map("Commerce_Product_SimpleProduct", graphqlProductDto.SimpleProduct{})
	types.Map("Commerce_Product_ConfigurableProduct", graphqlProductDto.ConfigurableProduct{})
	types.Map("Commerce_Product_ActiveVariantProduct", graphqlProductDto.ActiveVariantProduct{})
	types.Map("Commerce_Product_BundleProduct", graphqlProductDto.BundleProduct{})
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
	types.Map("Commerce_Product_Choice", graphqlProductDto.Choice{})
	types.Map("Commerce_Product_Option", graphqlProductDto.Option{})
	types.Map("Commerce_Product_ChoiceConfigurationInput", graphqlProductDto.ChoiceConfiguration{})
	types.Map("Commerce_Product_VariantSelection", graphqlProductDto.VariantSelection{})
	types.Map("Commerce_Product_VariantSelection_Attribute", graphqlProductDto.VariantSelectionAttribute{})
	types.Map("Commerce_Product_VariantSelection_Attribute_Option", graphqlProductDto.VariantSelectionAttributeOption{})
	types.Map("Commerce_Product_VariantSelection_Option_OtherAttributesRestriction", graphqlProductDto.OtherAttributesRestriction{})
	types.Map("Commerce_Product_VariantSelection_Match", graphqlProductDto.VariantSelectionMatch{})
	types.Map("Commerce_Product_VariantSelection_Match_Attributes", graphqlProductDto.VariantSelectionMatchAttributes{})
	types.Map("Commerce_Product_VariantSelection_Match_Variant", graphqlProductDto.VariantSelectionMatchVariant{})

	types.Resolve("Query", "Commerce_Product", CommerceProductQueryResolver{}, "CommerceProduct")
	types.Resolve("Query", "Commerce_Product_Search", CommerceProductQueryResolver{}, "CommerceProductSearch")
}
