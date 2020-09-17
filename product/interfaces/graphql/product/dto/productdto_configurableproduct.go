package graphqlProductDto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// A configurable without active variant
	ConfigurableProduct struct {
		product productDomain.BasicProduct
	}
)

//Type
func (cp ConfigurableProduct) Type() string {
	return productDomain.TypeConfigurable
}

func (cp ConfigurableProduct) Product() productDomain.BasicProduct {
	return cp.product
}

//MarketPlaceCode
func (cp ConfigurableProduct) MarketPlaceCode() string {
	return cp.product.BaseData().MarketPlaceCode
}

//Media
func (cp ConfigurableProduct) Media() ProductMedia {
	return ProductMedia{All: cp.product.TeaserData().Media}
}

//Price
func (cp ConfigurableProduct) Price() productDomain.PriceInfo {
	return cp.product.TeaserData().TeaserPrice
}

//Title
func (cp ConfigurableProduct) Title() string {
	return cp.product.BaseData().Title
}

//ProductCategories
func (cp ConfigurableProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: cp.product.BaseData().MainCategory,
		All:  cp.product.BaseData().Categories,
	}
}

//Description
func (cp ConfigurableProduct) Description() string {
	return cp.product.BaseData().Description
}

//ProductMeta
func (cp ConfigurableProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: cp.product.BaseData().Keywords,
	}
}

//ProductLoyalty
func (cp ConfigurableProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   cp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: cp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

//Attributes
func (cp ConfigurableProduct) Attributes() productDomain.Attributes {
	return cp.product.BaseData().Attributes
}

func (cp ConfigurableProduct) VariationSelections() []VariationSelection {
	return NewVariantsToVariationSelections(cp.Product())
}
