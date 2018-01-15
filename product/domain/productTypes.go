package domain

import (
	"github.com/pkg/errors"
)

const (
	// TYPESIMPLE denotes simple products
	TYPESIMPLE = "simple"
	// TYPECONFIGURABLE denotes configurable products
	TYPECONFIGURABLE = "configurable"
)

type (
	// BasicProduct interface - shared by TypeSimple and TypeConfigurable
	BasicProduct interface {
		BaseData() BasicProductData
		TeaserData() TeaserData
		SaleableData() Saleable
		Type() string
		GetIdentifier() string
		HasMedia(group string, usage string) bool
		GetMedia(group string, usage string) Media
		IsNew() bool
	}

	// SimpleProduct - A product without Variants that can be teasered and being sold
	SimpleProduct struct {
		Identifier string
		BasicProductData
		Saleable
		Teaser TeaserData
	}

	// ConfigurableProduct - A product that can be teasered and that has Sellable Variants Aggregated
	ConfigurableProduct struct {
		Identifier string
		BasicProductData
		Teaser                     TeaserData
		VariantVariationAttributes []string
		Variants                   []Variant
		ActiveVariant              *Variant
	}

	// Variant is a concrete kind of a product
	Variant struct {
		BasicProductData
		Saleable
	}
)

// Verify Interfaces
var _ BasicProduct = SimpleProduct{}
var _ BasicProduct = ConfigurableProduct{}

// Type interface implementation for SimpleProduct
func (p SimpleProduct) Type() string {
	return TYPESIMPLE
}

// BaseData interface implementation for SimpleProduct
func (p SimpleProduct) BaseData() BasicProductData {
	bp := p.BasicProductData
	return bp
}

// TeaserData interface implementation for SimpleProduct
func (p SimpleProduct) TeaserData() TeaserData {
	return p.Teaser
}

// SaleableData getter for SimpleProduct
func (p SimpleProduct) SaleableData() Saleable {
	return p.Saleable
}

// GetIdentifier interface implementation for SimpleProduct
func (p SimpleProduct) GetIdentifier() string {
	return p.Identifier
}

// HasMedia  for SimpleProduct
func (p SimpleProduct) HasMedia(group string, usage string) bool {
	media := findMediaInProduct(BasicProduct(p), group, usage)
	if media == nil {
		return false
	}
	return true
}

// GetMedia  for SimpleProduct
func (p SimpleProduct) GetMedia(group string, usage string) Media {
	return *findMediaInProduct(BasicProduct(p), group, usage)
}

func (p SimpleProduct) IsNew() bool {
	newFromDate := ""
	if p.HasAttribute("newFromDate") {
		newFromDate = p.Attributes["newFromDate"].Value()
	}

	newToDate := ""
	if p.HasAttribute("newToDate") {
		newToDate = p.Attributes["newToDate"].Value()
	}

	return isNew(newFromDate, newToDate)
}

// Type interface implementation for SimpleProduct
func (p ConfigurableProduct) Type() string {
	return TYPECONFIGURABLE
}

// GetIdentifier interface implementation for SimpleProduct
func (p ConfigurableProduct) GetIdentifier() string {
	return p.Identifier
}

// BaseData interface implementation for SimpleProduct
func (p ConfigurableProduct) BaseData() BasicProductData {
	bp := p.BasicProductData
	return bp
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

/*
	SaleableData getter for ConfigurableProduct
	Gets either the first or the active variants saleableData
*/
func (p ConfigurableProduct) SaleableData() Saleable {
	if p.HasActiveVariant() {
		return p.ActiveVariant.Saleable
	}
	if len(p.Variants) > 0 {
		return p.Variants[0].Saleable
	}
	return Saleable{}
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

// HasActiveVariant  for ConfigurableProduct
func (p ConfigurableProduct) HasActiveVariant() bool {
	return p.ActiveVariant != nil
}

func (p ConfigurableProduct) IsNew() bool {
	for _, variant := range p.Variants {
		if variant.IsNew() {
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

func (v Variant) IsNew() bool {
	newFromDate := ""
	if v.HasAttribute("newFromDate") {
		newFromDate = v.Attributes["newFromDate"].Value()
	}

	newToDate := ""
	if v.HasAttribute("newToDate") {
		newToDate = v.Attributes["newToDate"].Value()
	}

	return isNew(newFromDate, newToDate)
}
