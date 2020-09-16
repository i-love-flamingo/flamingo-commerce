package graphqlProductDto

import productDomain "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	// A Product variant that reflects one possible configuration of a configurable
	ActiveVariantProduct struct {
		product productDomain.BasicProduct
	}
)

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
	panic("implement me")
}

//Media
func (avp ActiveVariantProduct) Media() ProductMedia {
	panic("implement me")
}

//Price
func (avp ActiveVariantProduct) Price() productDomain.PriceInfo {
	panic("implement me")
}

//Title
func (avp ActiveVariantProduct) Title() string {
	panic("implement me")
}

//ProductCategories
func (avp ActiveVariantProduct) Categories() ProductCategories {
	panic("implement me")
}

//Description
func (avp ActiveVariantProduct) Description() string {
	panic("implement me")
}

//ProductMeta
func (avp ActiveVariantProduct) Meta() ProductMeta {
	panic("implement me")
}

//ProductLoyalty
func (avp ActiveVariantProduct) Loyalty() ProductLoyalty {
	panic("implement me")
}

//Attributes
func (avp ActiveVariantProduct) Attributes() productDomain.Attributes {
	panic("implement me")
}

func (avp ActiveVariantProduct) VariationSelections() []VariationSelection {
	panic("implement me")
}

func (avp ActiveVariantProduct) VariantMarketPlaceCode() string {
	panic("implement me")
}

func (avp ActiveVariantProduct) ActiveVariationSelections() []ActiveVariationSelection {
	panic("implement me")
}

func (avp ActiveVariantProduct) BaseConfigurableProduct() ConfigurableProduct {
	panic("implement me")
}
