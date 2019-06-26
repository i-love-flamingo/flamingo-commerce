package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

type Service struct{}

func (*Service) Schema() []byte {
	// language=graphql
	return []byte(`
interface Commerce_Product {
	baseData: Commerce_BasicProductData!
	teaserData: Commerce_ProductTeaserData!
	specifications: Commerce_ProductSpecifications!
	isSaleable: Boolean!
	saleableData: Commerce_ProductSaleable!
	type: String!
	getIdentifier: String!
	hasMedia(group: String!, usage: String!): Boolean!
	getMedia(group: String!, usage: String!): Commerce_ProductMedia!
}

type Commerce_SimpleProduct implements Commerce_Product {
	baseData: Commerce_BasicProductData!
	teaserData: Commerce_ProductTeaserData!
	specifications: Commerce_ProductSpecifications!
	isSaleable: Boolean!
	saleableData: Commerce_ProductSaleable!
	type: String!
	getIdentifier: String!
	hasMedia(group: String!, usage: String!): Boolean!
	getMedia(group: String!, usage: String!): Commerce_ProductMedia!
}

type Commerce_BasicProductData {
	title:            String!
	attributes:       Commerce_ProductAttributes
	shortDescription: String!
	description:      String!
	media:            [Commerce_ProductMedia!]

	marketPlaceCode: String!
	retailerCode:    String!
	retailerSku:     String!
	retailerName:    String!

	createdAt:   Time!
	updatedAt:   Time!
	visibleFrom: Time!
	visibleTo:   Time!

#	Categories:   [Commerce_CategoryTeaser!]
#	MainCategory: [Commerce_CategoryTeaser!]

	categoryToCodeMapping: [String!]

	stockLevel: String!

	keywords: [String!]
	isNew:    Boolean!
}

type Commerce_ProductTeaserData {
	shortTitle: String
	shortDescription: String
	# TeaserPrice is the price that should be shown in teasers (listview)
	# teaserPrice : Commerce_ProductPriceInfo
	# TeaserPriceIsFromPrice - is set to true in cases where a product might have different prices (e.g. configurable)
	teaserPriceIsFromPrice: Boolean
	# PreSelectedVariantSku - might be set for configurables to give a hint to link to a variant of a configurable (That might be the case if a user filters for an attribute and in the teaser the variant with that attribute is shown)
	preSelectedVariantSku: String
	# Media
 	media : [Commerce_ProductMedia!]
	# The sku that should be used to link from Teasers
	marketPlaceCode: String
	#teaserAvailablePrices : [PriceInfo!]
	# TeaserLoyaltyPriceInfo - optional the Loyaltyprice that can be used for teaser (e.g. on listing views)
	#teaserLoyaltyPriceInfo *LoyaltyPriceInfo
}

type Commerce_ProductSpecifications {
	groups: [Commerce_ProductSpecificationGroup!]
}

type Commerce_ProductSpecificationGroup {
	title: String!
	entries: [Commerce_ProductSpecificationEntry!]
}

type Commerce_ProductSpecificationEntry {
	label: String!
	values: [String!]
}

type Commerce_ProductSaleable {
	isSaleable: Boolean!
	saleableFrom: Time
	saleableTo: Time
#	activePrice: Commerce_ProductPriceInfo
#	availablePrices: [Commerce_ProductPriceInfo!]
	# loyaltyPrices - Optional infos for products that can be payed in a loyalty program
#	loyaltyPrices: [Commerce_ProductLoyaltyPriceInfo!]
}

type Commerce_ProductMedia {
	type:      String!
	mimeType:  String!
	usage:     String!
	title:     String!
	reference: String!
}

type Commerce_ProductAttributes {
	getAttributeKeys: [String!]
	getAttributes: [Commerce_ProductAttribute!]
	hasAttribute(key: String!): Boolean
	getAttribute(key: String!): Commerce_ProductAttribute
}

type Commerce_ProductAttribute {
	code: String!
	label: String!
#	rawValue: String!
	unitCode: String!
}

extend type Query {
	Commerce_Product(marketplaceCode: String!): Commerce_Product
}
`)
}

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
		"Commerce_BasicProductData":          domain.BasicProductData{},
		"Commerce_ProductTeaserData":         domain.TeaserData{},
		"Commerce_ProductSpecifications":     domain.Specifications{},
		"Commerce_ProductSpecificationGroup": domain.SpecificationGroup{},
		"Commerce_ProductSpecificationEntry": domain.SpecificationEntry{},
		"Commerce_ProductSaleable":           domain.Saleable{},
		"Commerce_ProductMedia":              domain.Media{},
		"Commerce_ProductAttributes":         domain.Attributes{},
		"Commerce_ProductAttribute":          domain.Attribute{},
	}.Models()
}

type CommerceProductQueryResolver struct {
	productService domain.ProductService
}

func (r *CommerceProductQueryResolver) Inject(productService domain.ProductService) {
	r.productService = productService
}

func (r *CommerceProductQueryResolver) CommerceProduct(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	return r.productService.Get(ctx, marketplaceCode)
}
