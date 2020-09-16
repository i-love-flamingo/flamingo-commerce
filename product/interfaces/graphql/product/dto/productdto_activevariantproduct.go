package graphqlProductDto

import productDomain "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	// A Product variant that reflects one possible configuration of a configurable
	ActiveVariantProduct struct {
		product productDomain.ConfigurableProductWithActiveVariant
	}
)

//getActiveVariant
func (avp ActiveVariantProduct) getActiveVariant() productDomain.Variant {
	return avp.product.ActiveVariant
}

//Type
func (avp ActiveVariantProduct) Type() string {
	return productDomain.TypeConfigurableWithActiveVariant
}

//Product
func (avp ActiveVariantProduct) Product() productDomain.BasicProduct {
	return avp.product
}

//MarketPlaceCode
func (avp ActiveVariantProduct) MarketPlaceCode() string {
	return avp.getActiveVariant().BaseData().MarketPlaceCode
}

//Media
func (avp ActiveVariantProduct) Media() ProductMedia {
	return ProductMedia{All: avp.getActiveVariant().BaseData().Media}
}

//Price
func (avp ActiveVariantProduct) Price() productDomain.PriceInfo {
	return avp.product.TeaserData().TeaserPrice
}

//Title
func (avp ActiveVariantProduct) Title() string {
	return avp.getActiveVariant().BaseData().Title
}

//ProductCategories
func (avp ActiveVariantProduct) Categories() ProductCategories {
	return ProductCategories{
		Main: avp.getActiveVariant().BaseData().MainCategory,
		All:  avp.getActiveVariant().BaseData().Categories,
	}
}

//Description
func (avp ActiveVariantProduct) Description() string {
	return avp.getActiveVariant().BaseData().Description
}

//ProductMeta
func (avp ActiveVariantProduct) Meta() ProductMeta {
	return ProductMeta{
		Keywords: avp.getActiveVariant().BaseData().Keywords,
	}
}

//ProductLoyalty
func (avp ActiveVariantProduct) Loyalty() ProductLoyalty {
	return ProductLoyalty{
		Price:   avp.product.TeaserData().TeaserLoyaltyPriceInfo,
		Earning: avp.product.TeaserData().TeaserLoyaltyEarningInfo,
	}
}

//Attributes
func (avp ActiveVariantProduct) Attributes() productDomain.Attributes {
	return avp.getActiveVariant().BaseData().Attributes
}

func (avp ActiveVariantProduct) VariationSelections() []VariationSelection {
	panic("implement me")
}

func (avp ActiveVariantProduct) VariantMarketPlaceCode() string {
	return avp.getActiveVariant().MarketPlaceCode
}

func (avp ActiveVariantProduct) ActiveVariationSelections() []ActiveVariationSelection {
	panic("implement me")
}
