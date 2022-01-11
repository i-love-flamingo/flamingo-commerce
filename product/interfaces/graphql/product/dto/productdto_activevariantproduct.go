package graphqlproductdto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// ActiveVariantProduct is a Product variant that reflects one possible configuration of a configurable
	ActiveVariantProduct struct {
		product productDomain.ConfigurableProductWithActiveVariant
	}
)

// getActiveVariant returns the active product variant
func (avp ActiveVariantProduct) getActiveVariant() productDomain.Variant {
	return avp.product.ActiveVariant
}

// Type of the product
func (avp ActiveVariantProduct) Type() string {
	return productDomain.TypeConfigurableWithActiveVariant
}

// Product the basic product domain object
func (avp ActiveVariantProduct) Product() productDomain.BasicProduct {
	return avp.product
}

// MarketPlaceCode of the active variant
func (avp ActiveVariantProduct) MarketPlaceCode() string {
	return avp.product.BasicProductData.MarketPlaceCode
}

// Identifier of the active variant
func (avp ActiveVariantProduct) Identifier() string {
	return avp.product.GetIdentifier()
}

// Media of the active variant
func (avp ActiveVariantProduct) Media() ProductMedia {
	return ProductMedia{All: avp.getActiveVariant().BaseData().Media}
}

// Price of the active variant
func (avp ActiveVariantProduct) Price() productDomain.PriceInfo {
	return avp.product.ActiveVariant.ActivePrice
}

// AvailablePrices of the active variant
func (avp ActiveVariantProduct) AvailablePrices() []productDomain.PriceInfo {
	return avp.product.ActiveVariant.AvailablePrices
}

// Title of the active variant
func (avp ActiveVariantProduct) Title() string {
	return avp.getActiveVariant().BaseData().Title
}

// Categories of the active variant
func (avp ActiveVariantProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: avp.getActiveVariant().BaseData().MainCategory,
		All:  avp.getActiveVariant().BaseData().Categories,
	}
}

// Description of the active variant
func (avp ActiveVariantProduct) Description() string {
	return avp.getActiveVariant().BaseData().Description
}

// ShortDescription of the product
func (avp ActiveVariantProduct) ShortDescription() string {
	return avp.product.BaseData().ShortDescription
}

// Meta contains meta information from the active variant
func (avp ActiveVariantProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: avp.getActiveVariant().BaseData().Keywords,
	}
}

// Loyalty contains loyalty information of the active variant
func (avp ActiveVariantProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   avp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: avp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

// Attributes of the active variant
func (avp ActiveVariantProduct) Attributes() productDomain.Attributes {
	return avp.getActiveVariant().BaseData().Attributes
}

// VariationSelections contains information about the available variations for the product
func (avp ActiveVariantProduct) VariationSelections() []VariationSelection {
	return NewVariantsToVariationSelections(avp.product)
}

// VariantMarketPlaceCode of the active variant
func (avp ActiveVariantProduct) VariantMarketPlaceCode() string {
	return avp.getActiveVariant().MarketPlaceCode
}

// ActiveVariationSelections helper to easily access active variant attributes
func (avp ActiveVariantProduct) ActiveVariationSelections() []ActiveVariationSelection {
	variationSelections := avp.VariationSelections()
	var activeVariationSelections []ActiveVariationSelection
	for _, variationSelection := range variationSelections {
		for _, option := range variationSelection.Options {
			if option.State == VariationSelectionOptionStateActive {
				activeVariationSelections = append(activeVariationSelections, ActiveVariationSelection{
					Code:     variationSelection.Code,
					Label:    variationSelection.Label,
					Value:    option.Label,
					UnitCode: option.Variant.BaseData().Attribute(variationSelection.Code).UnitCode,
				})
			}
		}
	}

	return activeVariationSelections
}

// Badges of the active variant
func (avp ActiveVariantProduct) Badges() ProductBadges {
	return ProductBadges{
		All: avp.getActiveVariant().BaseData().Badges,
	}
}
