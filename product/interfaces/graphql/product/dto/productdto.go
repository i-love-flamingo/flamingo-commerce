package graphqlProductDto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

const (
	VariationSelectionOptionStateActive  VariationSelectionOptionState = "ACTIVE"   // Option is currently active
	VariationSelectionOptionStateMatch   VariationSelectionOptionState = "MATCH"    // Product exists for this option
	VariationSelectionOptionStateNoMatch VariationSelectionOptionState = "NO_MATCH" // Product does not exists for this option
)

// Product contains normalized product information regardless of being a variant or simple product
type (

	// Product interface - needs to be implemented by SimpleProducts and ConfigurableProducts
	Product interface {
		Type() string
		MarketPlaceCode() string
		Media() ProductMedia
		Price() productDomain.PriceInfo
		Title() string
		Categories() ProductCategories
		Description() string
		Meta() ProductMeta
		Loyalty() ProductLoyalty
		Attributes() productDomain.Attributes
	}

	// A simple product without variants
	SimpleProduct struct {
		marketPlaceCode string
		media           ProductMedia
		price           productDomain.PriceInfo
		title           string
		categories      ProductCategories
		description     string
		meta            ProductMeta
		loyalty         ProductLoyalty
		attributes      productDomain.Attributes
	}

	// A configurable without active variant
	ConfigurableProduct struct {
		marketPlaceCode     string
		media               ProductMedia
		price               productDomain.PriceInfo
		title               string
		categories          ProductCategories
		description         string
		meta                ProductMeta
		loyalty             ProductLoyalty
		attributes          productDomain.Attributes
		VariationSelections []VariationSelection
	}

	// A product variant that reflects one possible configuration of a configurable
	ActiveVariantProduct struct {
		marketPlaceCode           string
		media                     ProductMedia
		price                     productDomain.PriceInfo
		title                     string
		categories                ProductCategories
		description               string
		meta                      ProductMeta
		loyalty                   ProductLoyalty
		attributes                productDomain.Attributes
		VariationSelections       []VariationSelection
		VariantMarketPlaceCode    string
		ActiveVariationSelections []ActiveVariationSelection
		BaseConfigurableProduct   ConfigurableProduct
	}

	// All loyalty related information
	ProductLoyalty struct {
		Price   productDomain.LoyaltyPriceInfo
		Earning productDomain.LoyaltyEarningInfo
	}

	ProductMedia struct {
		All []productDomain.Media
	}

	// ProductCategories
	ProductCategories struct {
		Main productDomain.CategoryTeaser
		All  []productDomain.CategoryTeaser
	}

	// Normalized ProductMeta data
	ProductMeta struct {
		Keywords []string
	}

	// A selection for a product variation
	VariationSelection struct {
		Code    string
		Label   string
		Options []VariationSelectionOption
	}

	// One possible variation for the product
	VariationSelectionOption struct {
		Code                   string
		Label                  string
		State                  VariationSelectionOptionState
		VariantMarketPlaceCode string
	}

	// Possible state of option depending on active variant
	VariationSelectionOptionState string

	// The variation for the currently active variant
	ActiveVariationSelection struct {
		AttributeLabel string
		OptionLabel    string
	}
)

var (
	_ Product = SimpleProduct{}
	_ Product = ConfigurableProduct{}
	_ Product = ActiveVariantProduct{}
)

//Type
func (sp SimpleProduct) Type() string {
	return "simple"
}

//MarketPlaceCode
func (sp SimpleProduct) MarketPlaceCode() string {
	return sp.marketPlaceCode
}

//Media
func (sp SimpleProduct) Media() ProductMedia {
	return sp.media
}

//Price
func (sp SimpleProduct) Price() productDomain.PriceInfo {
	return sp.price
}

//Title
func (sp SimpleProduct) Title() string {
	return sp.title
}

//ProductCategories
func (sp SimpleProduct) Categories() ProductCategories {
	return sp.categories
}

//Description
func (sp SimpleProduct) Description() string {
	return sp.description
}

//ProductMeta
func (sp SimpleProduct) Meta() ProductMeta {
	return sp.meta
}

//ProductLoyalty
func (sp SimpleProduct) Loyalty() ProductLoyalty {
	return sp.loyalty
}

//Attributes
func (sp SimpleProduct) Attributes() productDomain.Attributes {
	return sp.attributes
}

//Type
func (cp ConfigurableProduct) Type() string {
	return "configurable"
}

//MarketPlaceCode
func (cp ConfigurableProduct) MarketPlaceCode() string {
	return cp.marketPlaceCode
}

//Media
func (cp ConfigurableProduct) Media() ProductMedia {
	return cp.media
}

//Price
func (cp ConfigurableProduct) Price() productDomain.PriceInfo {
	return cp.price
}

//Title
func (cp ConfigurableProduct) Title() string {
	return cp.title
}

//ProductCategories
func (cp ConfigurableProduct) Categories() ProductCategories {
	return cp.categories
}

//Description
func (cp ConfigurableProduct) Description() string {
	return cp.description
}

//ProductMeta
func (cp ConfigurableProduct) Meta() ProductMeta {
	return cp.meta
}

//ProductLoyalty
func (cp ConfigurableProduct) Loyalty() ProductLoyalty {
	return cp.loyalty
}

//Attributes
func (cp ConfigurableProduct) Attributes() productDomain.Attributes {
	return cp.attributes
}

//Type
func (avp ActiveVariantProduct) Type() string {
	return "activeVariant"
}

//MarketPlaceCode
func (avp ActiveVariantProduct) MarketPlaceCode() string {
	return avp.marketPlaceCode
}

//Media
func (avp ActiveVariantProduct) Media() ProductMedia {
	return avp.media
}

//Price
func (avp ActiveVariantProduct) Price() productDomain.PriceInfo {
	return avp.price
}

//Title
func (avp ActiveVariantProduct) Title() string {
	return avp.title
}

//ProductCategories
func (avp ActiveVariantProduct) Categories() ProductCategories {
	return avp.categories
}

//Description
func (avp ActiveVariantProduct) Description() string {
	return avp.description
}

//ProductMeta
func (avp ActiveVariantProduct) Meta() ProductMeta {
	return avp.meta
}

//ProductLoyalty
func (avp ActiveVariantProduct) Loyalty() ProductLoyalty {
	return avp.loyalty
}

//Attributes
func (avp ActiveVariantProduct) Attributes() productDomain.Attributes {
	return avp.attributes
}

// GetMedia returns the FIRST found product media by usage
func (pm ProductMedia) GetMedia(usage string) *productDomain.Media {
	for _, media := range pm.All {
		if media.Usage == usage {
			return &media
		}
	}
	return nil
}

func mapProductToConfigurableProductDto(cp productDomain.ConfigurableProduct) *ConfigurableProduct {
	return &ConfigurableProduct{
		marketPlaceCode: cp.BaseData().MarketPlaceCode,
		media:           ProductMedia{All: cp.TeaserData().Media},
		price:           cp.TeaserData().TeaserPrice,
		title:           cp.BaseData().Title, // TODO: Needs to come from variant
		categories: ProductCategories{
			Main: cp.BaseData().MainCategory,
			All:  cp.BaseData().Categories,
		},
		meta: ProductMeta{
			Keywords: cp.BaseData().Keywords,
		},
	}
}

func mapProductToSimpleProductDto(sp productDomain.SimpleProduct) *SimpleProduct {
	return &SimpleProduct{
		marketPlaceCode: sp.BaseData().MarketPlaceCode,
		media:           ProductMedia{All: sp.TeaserData().Media},
		price:           sp.TeaserData().TeaserPrice,
		title:           sp.BaseData().Title, // TODO: Needs to come from variant
		categories: ProductCategories{
			Main: sp.BaseData().MainCategory,
			All:  sp.BaseData().Categories,
		},
		meta: ProductMeta{
			Keywords: sp.BaseData().Keywords,
		},
	}
}

func MapProductToGraphqlProductDto(product productDomain.BasicProduct) Product {
	if product.Type() == productDomain.TypeConfigurable {
		configurableProduct := product.(productDomain.ConfigurableProduct)
		return mapProductToConfigurableProductDto(configurableProduct)
	}

	simpleProduct := product.(productDomain.SimpleProduct)
	return mapProductToSimpleProductDto(simpleProduct)
}
