package graphqlproductdto

import productDomain "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	// SimpleProduct A simple Product without variants
	SimpleProduct struct {
		product productDomain.SimpleProduct
	}
)

// Type the product type
func (sp SimpleProduct) Type() string {
	return productDomain.TypeSimple
}

// Product the basic product information
func (sp SimpleProduct) Product() productDomain.BasicProduct {
	return sp.product
}

// MarketPlaceCode of the product
func (sp SimpleProduct) MarketPlaceCode() string {
	return sp.product.BaseData().MarketPlaceCode
}

// Identifier of the product
func (sp SimpleProduct) Identifier() string {
	return sp.product.GetIdentifier()
}

// Media of the product
func (sp SimpleProduct) Media() ProductMedia {
	return ProductMedia{All: sp.product.TeaserData().Media}
}

// Price of the product
func (sp SimpleProduct) Price() productDomain.PriceInfo {
	return sp.product.Saleable.ActivePrice
}

// AvailablePrices of the product
func (sp SimpleProduct) AvailablePrices() []productDomain.PriceInfo {
	return sp.product.Saleable.AvailablePrices
}

// Title of the product
func (sp SimpleProduct) Title() string {
	return sp.product.BaseData().Title
}

// Categories of the product
func (sp SimpleProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: sp.product.BaseData().MainCategory,
		All:  sp.product.BaseData().Categories,
	}
}

// Description of the product
func (sp SimpleProduct) Description() string {
	return sp.product.BaseData().Description
}

// ShortDescription of the product
func (sp SimpleProduct) ShortDescription() string {
	return sp.product.BaseData().ShortDescription
}

// Meta of the product
func (sp SimpleProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: sp.product.BaseData().Keywords,
	}
}

// Loyalty of the product
func (sp SimpleProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   sp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: sp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

// Attributes of the product
func (sp SimpleProduct) Attributes() productDomain.Attributes {
	return sp.product.BaseData().Attributes
}

// Badges of the product
func (sp SimpleProduct) Badges() ProductBadges {
	return ProductBadges{
		All: sp.product.BaseData().Badges,
	}
}
