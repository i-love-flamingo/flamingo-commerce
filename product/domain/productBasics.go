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

	//CategoryTeaser - Represents some Teaser infos for Category
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
		Default          priceDomain.Price
		IsDiscounted     bool
		Discounted       priceDomain.Price
		DiscountText     string
		MinPointsToSpent big.Float
		MaxPointsToSpent big.Float
		Context          PriceContext
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
		//TeaserLoyaltyPriceInfo - optional the Loyaltyprice that can be used for teaser (e.g. on listing views)
		TeaserLoyaltyPriceInfo *LoyaltyPriceInfo
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

	//Charges - Represents the Charges the product need to be payed with
	Charges struct {
		chargesByType map[string]priceDomain.Charge
	}

	//WishedToPay - list of prices by type
	WishedToPay struct {
		priceByType map[string]priceDomain.Price
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


// GetLoyaltyPriceByType - returns the loyaltyentry that matches the type
func (p Saleable) GetLoyaltyPriceByType(ltype string) (*LoyaltyPriceInfo, bool) {
	for _, lp := range p.LoyaltyPrices {
		if lp.Type == ltype {
			return &lp,true
		}
	}
	return nil, false
}

// GetChargesToPay  Gets the Charges that need to be payed by type
func (p Saleable) GetChargesToPay(wishedToPay *WishedToPay) Charges {
	requiredCharges := make(map[string]priceDomain.Charge)
	valuedPrice := p.ActivePrice.GetFinalPrice()
	remainingMainChargeValue := valuedPrice.Amount()

	for _, loyaltyPrice := range p.LoyaltyPrices {
		chargeType := "loyalty." + loyaltyPrice.Type

		if !loyaltyPrice.GetFinalPrice().IsPositive() {
			continue
		}
		if loyaltyPrice.MinPointsToSpent.Cmp(big.NewFloat(0)) < 1 {
			continue
		}
		rateLoyaltyDefaultPrices := loyaltyPrice.GetRate(p.ActivePrice.Default)

		//loyaltyAmountToSpent - set as default without potential wish
		loyaltyAmountToSpent := loyaltyPrice.GetAmountToSpend(nil)
		if wishedToPay != nil {
			wishedPrice := wishedToPay.GetByType(chargeType)
			if wishedPrice != nil && wishedPrice.Currency() == loyaltyPrice.GetFinalPrice().Currency() {
				//Use the passed wishedPrice of that type
				loyaltyAmountToSpent = loyaltyPrice.GetAmountToSpend(wishedPrice.Amount())
			}
		}

		//loyaltyPriceValue - is the value of this points in the "real" currency
		loyaltyPriceValue := new(big.Float).Mul(&rateLoyaltyDefaultPrices, &loyaltyAmountToSpent)

		if loyaltyPriceValue.Cmp(big.NewFloat(0)) <= 0 {
			continue
		}
		//Add the loyalty charge and at the same time reduce the remainingValue
		remainingMainChargeValue = new(big.Float).Sub(remainingMainChargeValue, loyaltyPriceValue)
		requiredCharges[chargeType] = priceDomain.Charge{
			Price: priceDomain.NewFromBigFloat(loyaltyAmountToSpent, loyaltyPrice.GetFinalPrice().Currency()),
			Type:  chargeType,
			Value: priceDomain.NewFromBigFloat(*loyaltyPriceValue, valuedPrice.Currency()).GetPayable(),
		}
	}

	remainingMainChargePrice := priceDomain.NewFromBigFloat(*remainingMainChargeValue, valuedPrice.Currency()).GetPayable()
	requiredCharges[priceDomain.ChargeTypeMain] = priceDomain.Charge{
		Price: remainingMainChargePrice,
		Type:  priceDomain.ChargeTypeMain,
		Value: remainingMainChargePrice,
	}

	return Charges{chargesByType: requiredCharges}
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

//NewWishedToPay - factory to get new WishedToPay struct
func NewWishedToPay() WishedToPay {
	return WishedToPay{
		priceByType: make(map[string]priceDomain.Price),
	}
}

//Add - returns new WishedToPay instance with the given whish added
func (w WishedToPay) Add(ctype string, price priceDomain.Price) WishedToPay {
	w.priceByType[ctype] = price
	return w
}

//GetByType - returns the wihsed price for the given type or nil
func (w WishedToPay) GetByType(ctype string) *priceDomain.Price {
	if price, ok := w.priceByType[ctype]; ok {
		return &price
	}
	return nil
}

//GetByType - returns a charge of given type. If it was not found a Zero amount is returned and the second return value is false
func (c Charges) GetByType(ctype string) (priceDomain.Charge, bool) {
	if charge, ok := c.chargesByType[ctype]; ok {
		return charge, ok
	}
	return priceDomain.Charge{}, false
}

//GetRate - get the currency conversion rate of the current loyaltyprice Default price (not discounted) in relation to the passed value
func (l LoyaltyPriceInfo) GetRate(valuedPrice priceDomain.Price) big.Float {
	if !l.Default.IsPositive() {
		return *big.NewFloat(0)
	}
	return *new(big.Float).Quo(valuedPrice.Amount(), l.Default.Amount())
}


//GetFinalPrice - gets either the Default or the Discounted Loyaltyprice
func (l LoyaltyPriceInfo) GetFinalPrice() priceDomain.Price {
	if l.IsDiscounted && l.Discounted.IsLessThen(l.Default) {
		return l.Discounted
	}
	return l.Default
}


//GetAmountToSpend - takes the whishedamaount and evaluates min and max and returns the points that need to be payed.
func (l LoyaltyPriceInfo) GetAmountToSpend(wishedAmount *big.Float) big.Float {
	//less or equal - return min
	if wishedAmount == nil || l.MinPointsToSpent.Cmp(wishedAmount) > 0 {
		return l.MinPointsToSpent
	}
	//more then max - return max
	if l.MaxPointsToSpent.Cmp(wishedAmount) == -1 {
		return l.MaxPointsToSpent
	}
	return *wishedAmount
}
