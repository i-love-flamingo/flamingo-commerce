package domain

import (
	"errors"
)

const (
	// TypeConfigurable denotes configurable products
	TypeConfigurable = "configurable"

	// TypeConfigurableWithActiveVariant denotes configurable products that has a variant selected
	TypeConfigurableWithActiveVariant = "configurable_with_activevariant"
)

type (

	// Configurable - interface that is implemented by ConfigurableProduct and ConfigurableProductWithActiveVariant
	Configurable interface {
		GetConfigurableWithActiveVariant(variantMarketplaceCode string) (ConfigurableProductWithActiveVariant, error)
		Variant(variantMarketplaceCode string) (*Variant, error)
		GetDefaultVariant() (*Variant, error)
		HasVariant(variantMarketplaceCode string) bool
	}

	// ConfigurableProduct - A product that can be teasered and that has Sellable Variants Aggregated
	ConfigurableProduct struct {
		Identifier string
		BasicProductData
		Teaser                            TeaserData
		VariantVariationAttributes        []string
		Variants                          []Variant
		VariantVariationAttributesSorting map[string][]string
	}

	// ConfigurableProductWithActiveVariant - A product that can be teasered and that has Sellable Variants Aggregated, One Variant is Active
	ConfigurableProductWithActiveVariant struct {
		Identifier string
		BasicProductData
		Teaser                            TeaserData
		VariantVariationAttributes        []string
		Variants                          []Variant
		VariantVariationAttributesSorting map[string][]string
		ActiveVariant                     Variant
	}

	// Variant is a concrete kind of a product
	Variant struct {
		BasicProductData
		Saleable
	}
)

var _ BasicProduct = ConfigurableProduct{}
var _ BasicProduct = ConfigurableProductWithActiveVariant{}
var _ Configurable = ConfigurableProduct{}
var _ Configurable = ConfigurableProductWithActiveVariant{}

// Type interface implementation for SimpleProduct
func (p ConfigurableProduct) Type() string {
	return TypeConfigurable
}

// IsSaleable defaults to false
func (p ConfigurableProduct) IsSaleable() bool {
	return false
}

// SaleableData getter for ConfigurableProduct - Configurable is NOT Salable
func (p ConfigurableProduct) SaleableData() Saleable {
	return Saleable{}
}

// GetConfigurableWithActiveVariant getter
func (p ConfigurableProduct) GetConfigurableWithActiveVariant(variantMarketplaceCode string) (ConfigurableProductWithActiveVariant, error) {
	variant, err := p.Variant(variantMarketplaceCode)
	if err != nil {
		return ConfigurableProductWithActiveVariant{}, err
	}
	return ConfigurableProductWithActiveVariant{
		Identifier:                        p.Identifier,
		BasicProductData:                  p.BasicProductData,
		Teaser:                            p.Teaser,
		VariantVariationAttributes:        p.VariantVariationAttributes,
		VariantVariationAttributesSorting: p.VariantVariationAttributesSorting,
		Variants:                          p.Variants,
		ActiveVariant:                     *variant,
	}, nil
}

// GetIdentifier interface implementation for SimpleProduct
func (p ConfigurableProduct) GetIdentifier() string {
	return p.Identifier
}

// BaseData interface implementation for ConfigurableProduct
func (p ConfigurableProduct) BaseData() BasicProductData {
	return p.BasicProductData
}

// TeaserData interface implementation for SimpleProduct
func (p ConfigurableProduct) TeaserData() TeaserData {
	return p.Teaser
}

// Variant getter for ConfigurableProduct
// Variant is retrieved via marketplaceCode of the variant
func (p ConfigurableProduct) Variant(variantMarketplaceCode string) (*Variant, error) {
	for _, variant := range p.Variants {
		if variant.BasicProductData.MarketPlaceCode == variantMarketplaceCode {
			return &variant, nil
		}
	}
	return nil, errors.New("No Variant with code " + variantMarketplaceCode + " found ")
}

// GetDefaultVariant getter
func (p ConfigurableProduct) GetDefaultVariant() (*Variant, error) {
	if len(p.Variants) > 0 {
		return &p.Variants[0], nil
	}
	return nil, errors.New("there is no variant. ")
}

// HasMedia for ConfigurableProduct
func (p ConfigurableProduct) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(p), group, usage)
	return media != nil
}

// GetMedia  for ConfigurableProduct
func (p ConfigurableProduct) GetMedia(group string, usage string) Media {
	return *findMediaInProduct(BasicProduct(p), group, usage)
}

// HasVariant  for ConfigurableProduct
func (p ConfigurableProduct) HasVariant(variantMarketplaceCode string) bool {
	for _, variant := range p.Variants {
		if variant.BasicProductData.MarketPlaceCode == variantMarketplaceCode {
			return true
		}
	}
	return false
}

// BaseData getter for BasicProductData
func (v Variant) BaseData() BasicProductData {
	return v.BasicProductData
}

// SaleableData getter for Saleable
func (v Variant) SaleableData() Saleable {
	return v.Saleable
}

// ********CONFIGURABLE WITH ACTIVE VARIANT

// Type getter
func (p ConfigurableProductWithActiveVariant) Type() string {
	return TypeConfigurableWithActiveVariant
}

// IsSaleable is true
func (p ConfigurableProductWithActiveVariant) IsSaleable() bool {
	return true
}

// GetIdentifier getter
func (p ConfigurableProductWithActiveVariant) GetIdentifier() string {
	return p.Identifier
}

// BaseData returns only BaseData for Active Variant. If you need the BaseData of the Configurable - use ConfigurableBaseData()
func (p ConfigurableProductWithActiveVariant) BaseData() BasicProductData {
	return p.ActiveVariant.BasicProductData
}

// ConfigurableBaseData getter
func (p ConfigurableProductWithActiveVariant) ConfigurableBaseData() BasicProductData {
	return p.BasicProductData
}

// TeaserData interface implementation for SimpleProduct
func (p ConfigurableProductWithActiveVariant) TeaserData() TeaserData {
	return p.Teaser
}

// Variant getter for ConfigurableProduct
// Variant is retrieved via marketplaceCode of the variant
func (p ConfigurableProductWithActiveVariant) Variant(variantMarketplaceCode string) (*Variant, error) {
	for _, variant := range p.Variants {
		if variant.BasicProductData.MarketPlaceCode == variantMarketplaceCode {
			return &variant, nil
		}
	}
	return nil, errors.New("No Variant with code " + variantMarketplaceCode + " found ")
}

// GetDefaultVariant getter
func (p ConfigurableProductWithActiveVariant) GetDefaultVariant() (*Variant, error) {
	if len(p.Variants) > 0 {
		return &p.Variants[0], nil
	}
	return nil, errors.New("there is no variant. ")
}

// SaleableData getter for ConfigurableProduct
// Gets either the first or the active variants saleableData
func (p ConfigurableProductWithActiveVariant) SaleableData() Saleable {
	return p.ActiveVariant.Saleable
}

// HasMedia  for ConfigurableProduct
func (p ConfigurableProductWithActiveVariant) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(p), group, usage)
	return media != nil
}

// GetMedia  for ConfigurableProduct
func (p ConfigurableProductWithActiveVariant) GetMedia(group string, usage string) Media {
	return *findMediaInProduct(BasicProduct(p), group, usage)
}

// HasVariant  for ConfigurableProduct
func (p ConfigurableProductWithActiveVariant) HasVariant(variantMarketplaceCode string) bool {
	for _, variant := range p.Variants {
		if variant.BasicProductData.MarketPlaceCode == variantMarketplaceCode {
			return true
		}
	}
	return false
}

// GetConfigurableWithActiveVariant getter
func (p ConfigurableProductWithActiveVariant) GetConfigurableWithActiveVariant(variantMarketplaceCode string) (ConfigurableProductWithActiveVariant, error) {
	variant, err := p.Variant(variantMarketplaceCode)
	if err != nil {
		return ConfigurableProductWithActiveVariant{}, err
	}
	return ConfigurableProductWithActiveVariant{
		Identifier:                 p.Identifier,
		BasicProductData:           p.BasicProductData,
		Teaser:                     p.Teaser,
		VariantVariationAttributes: p.VariantVariationAttributes,
		Variants:                   p.Variants,
		ActiveVariant:              *variant,
	}, nil
}
