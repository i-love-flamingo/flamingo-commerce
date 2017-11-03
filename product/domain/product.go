package domain

import (
	"time"

	"fmt"

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
	}

	// SimpleProduct - A product without Variants that can be teasered and selled
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
	}

	// Variant is a concrete kind of a product
	Variant struct {
		BasicProductData
		Saleable
	}

	// BasicProductData is the basic product model
	BasicProductData struct {
		Title            string
		Attributes       Attributes
		ShortDescription string
		Description      string
		Media            []Media

		MarketPlaceCode string
		RetailerCode    string
		RetailerSku     string

		CreatedAt   time.Time
		UpdatedAt   time.Time
		VisibleFrom time.Time
		VisibleTo   time.Time

		CategoryPath          []string
		CategoryCodes         []string
		CategoryToCodeMapping []string

		Keywords []string
	}

	// Saleable are properties required for beeing selled
	Saleable struct {
		IsSaleable      bool
		SaleableFrom    time.Time
		SaleableTo      time.Time
		ActivePrice     PriceInfo
		AvailablePrices []PriceInfo
	}

	// PriceInfo holds product price information
	PriceInfo struct {
		Default           float64
		Discounted        float64
		DiscountText      string
		Currency          string
		ActiveBase        float64
		ActiveBaseAmount  float64
		ActiveBaseUnit    string
		IsDiscounted      bool
		CampaignRules     []string
		DenyMoreDiscounts bool
		Context           PriceContext
	}

	// PriceContext defines the scope in which the price was calculated
	PriceContext struct {
		CustomerGroup string
		ChannelCode   string
		Locale        string
	}

	// TeaserData is the teaser-information for product previews
	TeaserData struct {
		ShortTitle       string
		ShortDescription string
		Media            []Media
	}

	// Media holds product media information
	Media struct {
		Type      string
		MimeType  string
		Usage     string
		Title     string
		Reference string
	}

	// Attributes describe a product attributes map
	Attributes map[string]Attribute

	// Attribute for product attributes
	Attribute struct {
		Code     string
		Label    string
		RawValue interface{}
	}
)

// Verify Interfaces
var _ BasicProduct = SimpleProduct{}
var _ BasicProduct = ConfigurableProduct{}

// Value returns the raw value
func (at Attribute) Value() string {
	return fmt.Sprintf("%v", at.RawValue)
}

// HasMultipleValues checks for multiple raw values
func (at Attribute) HasMultipleValues() bool {
	_, ok := at.RawValue.([]interface{})
	return ok
}

// Values builds a list of product attribute values
func (at Attribute) Values() []string {
	var result []string
	list, ok := at.RawValue.([]interface{})
	if ok {
		for _, entry := range list {
			result = append(result, fmt.Sprintf("%v", entry))
		}
	}
	return result
}

// HasAttribute check
func (bpd BasicProductData) HasAttribute(key string) bool {
	if _, ok := bpd.Attributes[key]; ok {
		return true
	}
	return false
}

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

// SaleableData getter for ConfigurableProduct
func (p ConfigurableProduct) SaleableData() Saleable {
	return p.Variants[0].Saleable
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

// BaseData getter for BasicProductData
func (v Variant) BaseData() BasicProductData {
	return v.BasicProductData
}

// SaleableData getter for Saleable
func (v Variant) SaleableData() Saleable {
	return v.Saleable
}

// GetFinalPrice getter for price that should be used in calculations (either discounted or default)
func (p PriceInfo) GetFinalPrice() float64 {
	if p.IsDiscounted {
		return p.Discounted
	}
	return p.Default
}

// GetListMedia returns the product media for listing
func (bpd BasicProductData) GetListMedia() Media {
	var emptyMedia Media
	for _, media := range bpd.Media {
		if media.Usage == "list" {
			return media
		}
	}
	return emptyMedia
}

func findMediaInProduct(p BasicProduct, group string, usage string) *Media {
	var mediaList []Media
	if group == "teaser" {
		mediaList = p.TeaserData().Media
		for _, media := range mediaList {
			if media.Usage == usage {
				return &media
			}
		}
	}

	mediaList = p.BaseData().Media
	for _, media := range mediaList {
		if media.Usage == usage {
			return &media
		}
	}
	return nil
}
