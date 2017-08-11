package domain

import (
	"time"

	"github.com/pkg/errors"
)

type (

	// BasicProduct interface - shared by Simple and Configurable
	BasicProduct interface {
		GetBaseData() BasicProductData
		GetTeaserData() TeaserData
		GetType() string
	}

	// SimpleProduct - A product without Variants that can be teasered and selled
	SimpleProduct struct {
		Identifier string
		BasicProductData
		SaleableData
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
		SaleableData
	}

	// basicProductData is the basic product model
	BasicProductData struct {
		MarketPlaceCode  string
		Title            string
		Attributes       Attributes
		ShortDescription string
		Description      string
		Media            []Media
		RetailerCode     string
		CreatedAt        time.Time
		UpdatedAt        time.Time
		VisibleFrom      time.Time
		VisibleTo        time.Time

		CategoryPath  []string
		CategoryCodes []string

		Keywords []string
	}

	// SaleableData are properties required for beeing selled
	SaleableData struct {
		IsSaleable      bool
		SaleableFrom    time.Time
		SaleableTo      time.Time
		ActivePrice     PriceInfo
		AvailablePrices []PriceInfo
		RetailerSku     string
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
		Title            string
		ShortTitle       string
		Teaser           string
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

const mediaTypeExternalImage string = "image-external"
const mediaTypeImageService string = "image-service"

// Verify Interfaces
var _ BasicProduct = SimpleProduct{}
var _ BasicProduct = ConfigurableProduct{}

// GetType interface implementation for SimpleProduct
func (ps SimpleProduct) GetType() string {
	return "simple"
}

// GetBaseData interface implementation for SimpleProduct
func (ps SimpleProduct) GetBaseData() BasicProductData {
	bp := ps.BasicProductData
	return bp
}

// GetTeaserData interface implementation for SimpleProduct
func (ps SimpleProduct) GetTeaserData() TeaserData {
	return ps.Teaser
}

// GetType interface implementation for SimpleProduct
func (p ConfigurableProduct) GetType() string {
	return "configurable"
}

// GetBaseData interface implementation for SimpleProduct
func (p ConfigurableProduct) GetBaseData() BasicProductData {
	bp := p.BasicProductData
	return bp
}

// GetTeaserData interface implementation for SimpleProduct
func (p ConfigurableProduct) GetTeaserData() TeaserData {
	return p.Teaser
}

// GetBaseData interface implementation for SimpleProduct
func (p ConfigurableProduct) GetVariant(marketplaceCode string) (*Variant, error) {

	for _, variant := range p.Variants {
		if variant.MarketPlaceCode == marketplaceCode {
			return &variant, nil
		}
	}
	return nil, errors.New("No Variant with code " + marketplaceCode + " found ")
}
