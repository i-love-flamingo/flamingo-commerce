package dto

import (
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

// Product contains normalized product information regardless of being a variant or simple product
type (

	// Product interface - needs to be implemented by SimpleProducts and ConfigurableProducts
	Product interface {
		Type() string
	}

	// A simple product without variants
	SimpleProduct struct {
		MarketPlaceCode string
		Media           []productDomain.Media
		Price           productDomain.PriceInfo
		Title           string
		Categories      Categories
		Description     string
		Meta            Meta
		Loyalty			Loyalty
		Attributes		productDomain.Attributes
	}

	// A product variant that reflects one possibible configuration of a configurable
	ConfigurableProduct struct {
		MarketPlaceCode string
		Media           []productDomain.Media
		Price           productDomain.PriceInfo
		Title           string
		Categories      Categories
		Description     string
		Meta            Meta
		Attributes		productDomain.Attributes
	}

	// All loyalty related information
	Loyalty struct {
		Price 	productDomain.LoyaltyPriceInfo
		Earning productDomain.LoyaltyEarningInfo
	}

	// Categories
	Categories struct {
		Main productDomain.CategoryTeaser
		All  []productDomain.CategoryTeaser
	}

	// Retailer information
	//RetailerInfo struct {
	//	Id 		string
	//	Title 	string
	//}
	//
	//// Brand information
	//BrandInfo struct {
	//	Id 		string
	//	Title 	string
	//}

	// Normalized Meta data
	Meta struct {
		Description string
		Title		string
		Keywords	[]string
	}
)

func (sp SimpleProduct) Type() string {
	return "simple"
}

func (cp ConfigurableProduct) Type() string {
	return "configurable"
}


func MapProductToConfigurableProductDto(configurableProduct productDomain.ConfigurableProduct) *ConfigurableProduct {
	return &ConfigurableProduct{
		MarketPlaceCode: configurableProduct.BaseData().MarketPlaceCode,
		Media: configurableProduct.TeaserData().Media,
		Price: configurableProduct.TeaserData().TeaserPrice,
		Title: configurableProduct.BaseData().Title, // TODO: Needs to come from variant
		Categories: Categories{
			Main: configurableProduct.BaseData().MainCategory,
			All: configurableProduct.BaseData().Categories,
		},
		Meta: Meta{
			Keywords: configurableProduct.BaseData().Keywords,
		},
	}
}

func MapProductToSimpleProductDto(simpleProduct productDomain.SimpleProduct) *SimpleProduct {
	return &SimpleProduct{
		MarketPlaceCode: simpleProduct.BaseData().MarketPlaceCode,
		Media: simpleProduct.TeaserData().Media,
		Price: simpleProduct.TeaserData().TeaserPrice,
		Title: simpleProduct.BaseData().Title, // TODO: Needs to come from variant
		Categories: Categories{
			Main: simpleProduct.BaseData().MainCategory,
			All: simpleProduct.BaseData().Categories,
		},
		Meta: Meta{
			Keywords: simpleProduct.BaseData().Keywords,
		},
	}
}
