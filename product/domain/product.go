package domain

import (
	"time"

	"github.com/pkg/errors"
)

const (
	// TypeSimple denotes simple products
	TypeSimple = "simple"
	// TypeConfigurable denotes configurable products
	TypeConfigurable = "configurable"
)

type (
	// BasicProduct interface - shared by TypeSimple and TypeConfigurable
	BasicProduct interface {
		BaseData() BasicProductData
		TeaserData() TeaserData
		SaleableData() Saleable
		Type() string
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

	Variant struct {
		BasicProductData
		Saleable
	}

	// baseData is the basic product model
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

	// SaleableData are properties required for beeing selled
	Saleable struct {
		IsSaleable      bool
		SaleableFrom    time.Time
		SaleableTo      time.Time
		ActivePrice     PriceInfo
		AvailablePrices []PriceInfo
	}

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
	return TypeSimple
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

func (ps SimpleProduct) SaleableData() Saleable {
	return ps.Saleable
}

// Type interface implementation for SimpleProduct
func (p ConfigurableProduct) Type() string {
	return TypeConfigurable
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

// BaseData interface implementation for SimpleProduct
func (p ConfigurableProduct) Variant(marketplaceCode string) (*Variant, error) {
	for _, variant := range p.Variants {
		if variant.BasicProductData.MarketPlaceCode == marketplaceCode {
			return &variant, nil
		}
	}
	return nil, errors.New("No Variant with code " + marketplaceCode + " found ")
}

func (p ConfigurableProduct) SaleableData() Saleable {
	return p.Variants[0].Saleable
}

func (v Variant) BaseData() BasicProductData {
	return v.BasicProductData
}

func (v Variant) SaleableData() Saleable {
	return v.Saleable
}
