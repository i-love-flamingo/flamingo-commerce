package domain

import (
	"github.com/pkg/errors"
)

const (

	// TYPECONFIGURABLE denotes configurable products
	TYPECONFIGURABLE = "configurable"

	// TYPECONFIGURABLE denotes configurable products that has a variant selected
	TYPECONFIGURABLE_WITH_ACTIVE_VARIANT = "configurable_with_activevariant"
)

type (

	// ConfigurableProduct - A product that can be teasered and that has Sellable Variants Aggregated
	ConfigurableProduct struct {
		Identifier string
		BasicProductData
		Teaser                     TeaserData
		VariantVariationAttributes []string
		Variants                   []Variant
	}

	// ConfigurableProductWithActiveVariant - A product that can be teasered and that has Sellable Variants Aggregated, One Variant is Active
	ConfigurableProductWithActiveVariant struct {
		Identifier string
		BasicProductData
		Teaser                     TeaserData
		VariantVariationAttributes []string
		Variants                   []Variant
		ActiveVariant              Variant
	}

	// Variant is a concrete kind of a product
	Variant struct {
		BasicProductData
		Saleable
	}
)

var _ BasicProduct = ConfigurableProduct{}
var _ BasicProduct = ConfigurableProductWithActiveVariant{}

// Type interface implementation for SimpleProduct
func (p ConfigurableProduct) Type() string {
	return TYPECONFIGURABLE
}

func (p ConfigurableProduct) IsSaleable() bool {
	return false
}

//SaleableData getter for ConfigurableProduct - Configurable is NOT Salable
func (p ConfigurableProduct) SaleableData() Saleable {
	return Saleable{}
}

func (p ConfigurableProduct) GetConfigurableWithActiveVariant(variantMarketplaceCode string) (ConfigurableProductWithActiveVariant, error) {

	variant, err := p.Variant(variantMarketplaceCode)
	if err != nil {
		return ConfigurableProductWithActiveVariant{}, err
	}
	return ConfigurableProductWithActiveVariant{
		Identifier:       p.Identifier,
		BasicProductData: p.BasicProductData,
		Teaser:           p.Teaser,
		VariantVariationAttributes: p.VariantVariationAttributes,
		Variants:                   p.Variants,
		ActiveVariant:              *variant,
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

// GetDefaultVariant
func (p ConfigurableProduct) GetDefaultVariant() (*Variant, error) {
	if len(p.Variants) > 0 {
		return &p.Variants[0], nil
	}
	return nil, errors.New("There is no Variant")
}

// HasMedia  for ConfigurableProduct
func (p ConfigurableProduct) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(p), group, usage)
	if media == nil {
		return false
	}
	return true
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

//********CONFIGURABLE WITH ACTIVE VARIANT

func (p ConfigurableProductWithActiveVariant) Type() string {
	return TYPECONFIGURABLE_WITH_ACTIVE_VARIANT
}

func (p ConfigurableProductWithActiveVariant) IsSaleable() bool {
	return true
}

func (p ConfigurableProductWithActiveVariant) GetIdentifier() string {
	return p.Identifier
}

// Returns only BaseData for Active Variant. If you need the BaseData of the Configurable - use ConfigurableBaseData()
func (p ConfigurableProductWithActiveVariant) BaseData() BasicProductData {
	return p.ActiveVariant.BasicProductData
}

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

// GetDefaultVariant
func (p ConfigurableProductWithActiveVariant) GetDefaultVariant() (*Variant, error) {
	if len(p.Variants) > 0 {
		return &p.Variants[0], nil
	}
	return nil, errors.New("There is no Variant")
}

/*
	SaleableData getter for ConfigurableProduct
	Gets either the first or the active variants saleableData
*/
func (p ConfigurableProductWithActiveVariant) SaleableData() Saleable {
	return p.ActiveVariant.Saleable
}

// HasMedia  for ConfigurableProduct
func (p ConfigurableProductWithActiveVariant) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(p), group, usage)
	if media == nil {
		return false
	}
	return true
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
