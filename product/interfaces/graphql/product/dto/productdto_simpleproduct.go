package graphqlProductDto

import productDomain "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	// A simple Product without variants
	SimpleProduct struct {
		product productDomain.BasicProduct
	}
)

//Type
func (sp SimpleProduct) Type() string {
	return productDomain.TypeSimple
}

func (sp SimpleProduct) Product() productDomain.BasicProduct {
	return sp.product
}

//MarketPlaceCode
func (sp SimpleProduct) MarketPlaceCode() string {
	return sp.product.BaseData().MarketPlaceCode
}

//Media
func (sp SimpleProduct) Media() ProductMedia {
	return ProductMedia{All: sp.product.TeaserData().Media}
}

//Price
func (sp SimpleProduct) Price() productDomain.PriceInfo {
	return sp.product.TeaserData().TeaserPrice
}

//Title
func (sp SimpleProduct) Title() string {
	return sp.product.BaseData().Title
}

//ProductCategories
func (sp SimpleProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: sp.product.BaseData().MainCategory,
		All:  sp.product.BaseData().Categories,
	}
}

//Description
func (sp SimpleProduct) Description() string {
	return sp.product.BaseData().Description
}

//ProductMeta
func (sp SimpleProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: sp.product.BaseData().Keywords,
	}
}

//ProductLoyalty
func (sp SimpleProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   sp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: sp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

//Attributes
func (sp SimpleProduct) Attributes() productDomain.Attributes {
	return sp.product.BaseData().Attributes
}
