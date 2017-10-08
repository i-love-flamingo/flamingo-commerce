package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	// TypeSimple denotes simple products
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
		Context           struct {
			CustomerGroup interface{} `json:"customerGroup"`
			ChannelCode   string      `json:"channelCode"`
			Locale        string      `json:"locale"`
		}
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

	// Attributes is a generic map[string]interface{}
	Attributes map[string]interface{}
)

// Verify Interfaces
var _ BasicProduct = SimpleProduct{}
var _ BasicProduct = ConfigurableProduct{}

// Type interface implementation for SimpleProduct
func (ps SimpleProduct) Type() string {
	return TYPESIMPLE
}

// BaseData interface implementation for SimpleProduct
func (ps SimpleProduct) BaseData() BasicProductData {
	bp := ps.BasicProductData
	return bp
}

// TeaserData interface implementation for SimpleProduct
func (ps SimpleProduct) TeaserData() TeaserData {
	return ps.Teaser
}

// SaleableData getter for SimpleProduct
func (ps SimpleProduct) SaleableData() Saleable {
	return ps.Saleable
}

// GetIdentifier interface implementation for SimpleProduct
func (p SimpleProduct) GetIdentifier() string {
	return p.Identifier
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

// GetListImage
func (b BasicProductData) GetListMedia() Media {
	var emptyMedia Media
	for _, media := range b.Media {
		if media.Usage == "list" {
			return media
		}
	}
	return emptyMedia
}
