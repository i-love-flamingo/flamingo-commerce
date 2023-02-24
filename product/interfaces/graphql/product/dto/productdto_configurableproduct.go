package graphqlproductdto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// ConfigurableProduct is a configurable without active variant
	ConfigurableProduct struct {
		product productDomain.ConfigurableProduct
	}
)

// Type of the configurable
func (cp ConfigurableProduct) Type() string {
	return productDomain.TypeConfigurable
}

// Product the base product
func (cp ConfigurableProduct) Product() productDomain.BasicProduct {
	return cp.product
}

// MarketPlaceCode of the configurable
func (cp ConfigurableProduct) MarketPlaceCode() string {
	return cp.product.BaseData().MarketPlaceCode
}

// Identifier of the configurable
func (cp ConfigurableProduct) Identifier() string {
	return cp.product.GetIdentifier()
}

// Media of the configurable
func (cp ConfigurableProduct) Media() ProductMedia {
	return ProductMedia{All: cp.product.TeaserData().Media}
}

// Price of the configurable
func (cp ConfigurableProduct) Price() productDomain.PriceInfo {
	return productDomain.PriceInfo{} // Price info is always empty for configurable products because they are not saleable
}

// AvailablePrices of the configurable
func (cp ConfigurableProduct) AvailablePrices() []productDomain.PriceInfo {
	return nil // AvailablePrices is always empty for configurable products because they are not saleable
}

// Title of the configurable
func (cp ConfigurableProduct) Title() string {
	return cp.product.BaseData().Title
}

// Categories of the configurable
func (cp ConfigurableProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: cp.product.BaseData().MainCategory,
		All:  cp.product.BaseData().Categories,
	}
}

// Description of the configurable
func (cp ConfigurableProduct) Description() string {
	return cp.product.BaseData().Description
}

// ShortDescription of the product
func (cp ConfigurableProduct) ShortDescription() string {
	return cp.product.BaseData().ShortDescription
}

// Meta metadata of the configurable
func (cp ConfigurableProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: cp.product.BaseData().Keywords,
	}
}

// Loyalty information about the configurable
func (cp ConfigurableProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   cp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: cp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

// Attributes of the configurable
func (cp ConfigurableProduct) Attributes() productDomain.Attributes {
	return cp.product.BaseData().Attributes
}

// VariantSelection contains possible combinations of variation attributes
func (cp ConfigurableProduct) VariantSelection() VariantSelection {
	return MapVariantSelections(cp.product)
}

// Badges of the configurable product
func (cp ConfigurableProduct) Badges() ProductBadges {
	return ProductBadges{
		All: cp.product.BaseData().Badges,
	}
}
