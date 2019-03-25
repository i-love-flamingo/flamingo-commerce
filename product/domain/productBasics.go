package domain

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

// Media usage constants
const (
	MediaUsageList      = "list"
	MediaUsageDetail    = "detail"
	MediaUsageThumbnail = "thumbnail"
)

type (
	// BasicProduct interface - Need to be implements by all Product Types!
	BasicProduct interface {
		BaseData() BasicProductData
		TeaserData() TeaserData
		GetSpecifications() Specifications
		//IsSaleable - indicates if that product type can be purchased
		IsSaleable() bool
		SaleableData() Saleable
		Type() string
		GetIdentifier() string
		HasMedia(group string, usage string) bool
		GetMedia(group string, usage string) Media
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
		RetailerName    string

		CreatedAt   time.Time
		UpdatedAt   time.Time
		VisibleFrom time.Time
		VisibleTo   time.Time

		Categories   []CategoryTeaser
		MainCategory CategoryTeaser

		CategoryToCodeMapping []string

		StockLevel string

		Keywords []string
		IsNew    bool
	}

	//Represents some Teaser infos for Category
	CategoryTeaser struct {
		//Code the idendifier of the Category
		Code string
		//The Path (root to leaf) for this Category - seperated by "/"
		Path string
		//Name - speaking name of the category
		Name string
	}

	// Saleable are properties required for being selled
	Saleable struct {
		IsSaleable      bool
		SaleableFrom    time.Time
		SaleableTo      time.Time
		ActivePrice     PriceInfo
		AvailablePrices []PriceInfo
		//LoyaltyPrices - Optional infos for products that can be payed in a loyalty program
		LoyaltyPrices []LoyaltyPriceInfo
	}

	// PriceInfo holds product price information
	PriceInfo struct {
		Default           priceDomain.Price
		Discounted        priceDomain.Price
		DiscountText      string
		ActiveBase        big.Float
		ActiveBaseAmount  big.Float
		ActiveBaseUnit    string
		IsDiscounted      bool
		CampaignRules     []string
		DenyMoreDiscounts bool
		Context           PriceContext
		TaxClass          string
	}

	//LoyaltyPriceInfo - contains info used for product with
	LoyaltyPriceInfo struct {
		//Type - Name( or Type) of the Loyalty program
		Type             string
		PointPrice       priceDomain.Price
		MinPointsToSpent big.Float
		MaxPointsToSpent big.Float
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
		//TeaserPrice is the price that should be shown in teasers (listview)
		TeaserPrice PriceInfo
		//TeaserPriceIsFromPrice - is set to true in cases where a product might have different prices (e.g. configurable)
		TeaserPriceIsFromPrice bool
		//PreSelectedVariantSku - might be set for configurables to give a hint to link to a variant of a configurable (That might be the case if a user filters for an attribute and in the teaser the variant with that attribute is shown)
		PreSelectedVariantSku string
		//Media
		Media []Media
		//The sku that should be used to link from Teasers
		MarketPlaceCode       string
		TeaserAvailablePrices []PriceInfo
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
		UnitCode string
	}

	// Specifications of a product
	Specifications struct {
		Groups []SpecificationGroup
	}

	// SpecificationGroup groups specifications
	SpecificationGroup struct {
		Title   string
		Entries []SpecificationEntry
	}

	// SpecificationEntry data
	SpecificationEntry struct {
		Label  string
		Values []string
	}
)

// Value returns the raw value
func (at Attribute) Value() string {
	return strings.Trim(fmt.Sprintf("%v", at.RawValue), " ")
}

// IsEnabledValue returns true if the value can be seen as a toogle and is enabled
func (at Attribute) IsEnabledValue() bool {
	switch at.RawValue {
	case "Yes", "yes":
		return true
	case "true", true:
		return true
	case "1", 1:
		return true
	default:
		return false
	}
}

// IsDisabledValue returns true if the value can be seen as a disable toggle/swicth value
func (at Attribute) IsDisabledValue() bool {
	switch at.RawValue {
	case "No", "no":
		return true
	case "false", false:
		return true
	case "0", 0:
		return true
	default:
		return false
	}
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

// HasUnitCode checks if a unit code is set on the attribute
func (at Attribute) HasUnitCode() bool {
	return len(at.UnitCode) > 0
}

// GetUnit returns the unit on an attribute
func (at Attribute) GetUnit() Unit {
	unit, ok := Units[at.UnitCode]
	if !ok {
		return Unit{
			Code:   at.UnitCode,
			Symbol: "",
		}
	}
	return unit
}

// HasAttribute check
func (bpd BasicProductData) HasAttribute(key string) bool {
	if _, ok := bpd.Attributes[key]; ok {
		return true
	}
	return false
}

// GetFinalPrice getter for price that should be used in calculations (either discounted or default)
func (p PriceInfo) GetFinalPrice() priceDomain.Price {
	if p.IsDiscounted {
		return p.Discounted
	}
	return p.Default
}

// GetListMedia returns the product media for listing
func (bpd BasicProductData) GetListMedia() Media {
	return bpd.GetMedia(MediaUsageList)
}

// GetSpecifications getter
func (bpd BasicProductData) GetSpecifications() Specifications {
	if specs, ok := bpd.Attributes["specifications"].RawValue.(Specifications); ok {
		return specs
	}

	return Specifications{}
}

// GetMedia returns the FIRST found product media by usage
func (bpd BasicProductData) GetMedia(usage string) Media {
	var emptyMedia Media
	for _, media := range bpd.Media {
		if media.Usage == usage {
			return media
		}
	}
	return emptyMedia
}

// IsSaleableNow  checks flag and time
func (p Saleable) IsSaleableNow() bool {
	if p.IsSaleable == false {
		return false
	}

	//For some reasons IsZero does not always work - thats why we check for 1970
	if (p.SaleableFrom.IsZero() || p.SaleableFrom.Year() == 1970 || p.SaleableFrom.Before(time.Now())) &&
		(p.SaleableTo.IsZero() || p.SaleableTo.Year() == 1970 || p.SaleableTo.After(time.Now())) {
		return true
	}

	return false
}

// GetChargesToPay  Gets the Charges that need to be payed
func (p Saleable) GetChargesToPay(whishedCharges []priceDomain.Charge) []priceDomain.Charge {
	var requiredCharges []priceDomain.Charge
	requiredCharges = append(requiredCharges, priceDomain.Charge{
		Price: p.ActivePrice.GetFinalPrice(),
		Type:  priceDomain.ChargeTypeMain,
	})
	return requiredCharges
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

// IsInStock returns information if current product whether in stock or not
func (bpd BasicProductData) IsInStock() bool {
	if bpd.HasAttribute("alwaysInStock") && bpd.Attributes["alwaysInStock"].Value() == "true" {
		return true
	}

	if bpd.StockLevel == "" || bpd.StockLevel == "out" {
		return false
	}

	return true
}
