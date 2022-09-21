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
	// BasicProduct interface need to be implemented by all Product Types!
	BasicProduct interface {
		// BaseData gives you basic product information
		// like attributes, title, media, marketplacecode etc
		BaseData() BasicProductData
		// TeaserData gives you basic information to typically show in list views - including a teaser prices
		TeaserData() TeaserData
		// GetSpecifications return grouped specifications - typically used for comparisons
		GetSpecifications() Specifications
		// IsSaleable indicates if that product type can be purchased in general
		IsSaleable() bool
		// SaleableData returns information required to purchase - like the definite price information
		// makes only sense to call if IsSaleable() returns true
		SaleableData() Saleable
		// Type returns product type
		Type() string
		// GetIdentifier returns an identifier for the product instance
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
		Badges           Badges

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

		// Deprecated: use Stock[x].Level instead
		StockLevel string
		Stock      []Stock

		Keywords []string
		IsNew    bool
	}

	// CategoryTeaser represents some Teaser infos for Category
	CategoryTeaser struct {
		// Code the identifier of the Category
		Code string
		// The Path (root to leaf) for this Category - separated by "/"
		Path string
		// Name is the speaking name of the category
		Name string
		// Parent is an optional link to parent teaser
		Parent *CategoryTeaser `swaggerignore:"true"`
	}

	// Saleable are properties required for being selled
	Saleable struct {
		IsSaleable      bool
		SaleableFrom    time.Time
		SaleableTo      time.Time
		ActivePrice     PriceInfo
		AvailablePrices []PriceInfo
		// LoyaltyPrices holds optional infos for products that can be paid in a loyalty program
		LoyaltyPrices []LoyaltyPriceInfo
		// LoyaltyEarnings holds optional infos about potential loyalty earnings
		LoyaltyEarnings []LoyaltyEarningInfo
	}

	// PriceInfo holds product price information
	PriceInfo struct {
		Default           priceDomain.Price
		Discounted        priceDomain.Price
		DiscountText      string
		ActiveBase        big.Float `swaggertype:"string"`
		ActiveBaseAmount  big.Float `swaggertype:"string"`
		ActiveBaseUnit    string
		IsDiscounted      bool
		CampaignRules     []string
		DenyMoreDiscounts bool
		Context           PriceContext
		TaxClass          string
	}

	// LoyaltyPriceInfo contains info used for product with
	LoyaltyPriceInfo struct {
		// Type or Name of the Loyalty program
		Type             string
		Default          priceDomain.Price
		IsDiscounted     bool
		Discounted       priceDomain.Price
		DiscountText     string
		MinPointsToSpent big.Float  `swaggertype:"string"`
		MaxPointsToSpent *big.Float `swaggertype:"string"`
		Context          PriceContext
	}

	// LoyaltyEarningInfo contains earning infos
	LoyaltyEarningInfo struct {
		Type    string
		Default priceDomain.Price
	}

	// PriceContext defines the scope in which the price was calculated
	PriceContext struct {
		DeliveryCode  string
		CustomerGroup string
		ChannelCode   string
		Locale        string
	}

	// TeaserData is the teaser-information for product previews
	TeaserData struct {
		ShortTitle       string
		ShortDescription string
		URLSlug          string
		// TeaserPrice is the price that should be shown in teasers (listview)
		TeaserPrice PriceInfo
		// TeaserPriceIsFromPrice is set to true in cases where a product might have different prices (e.g. configurable)
		TeaserPriceIsFromPrice bool
		// PreSelectedVariantSku might be set for configurables to give a hint to link to a variant of a configurable (That might be the case if a user filters for an attribute and in the teaser the variant with that attribute is shown)
		PreSelectedVariantSku string
		// Media
		Media []Media
		// The sku that should be used to link from Teasers
		MarketPlaceCode       string
		TeaserAvailablePrices []PriceInfo
		// TeaserLoyaltyPriceInfo is the loyalty price that can be used for teaser (e.g. on listing views)
		TeaserLoyaltyPriceInfo *LoyaltyPriceInfo
		// TeaserLoyaltyEarning is the teaser for the loyalty earning used in grid / list view
		TeaserLoyaltyEarningInfo *LoyaltyEarningInfo
		// Badges optional slice of badges to teaser a product
		Badges Badges
	}

	// Media holds product media information
	Media struct {
		Type      string
		MimeType  string
		Usage     string
		Title     string
		Reference string
	} // @name ProductMedia

	// Attributes describe a product attributes map
	Attributes map[string]Attribute // @name ProductAttributes

	// Attribute for product attributes
	// Example:
	// Attribute{
	//		Code:      "my-attribute",
	//		CodeLabel: "My Attribute",
	//		Label:     "My attribute value",
	//		RawValue:  2,
	//		UnitCode:  "PCS",
	//	}
	Attribute struct {
		// Code is the internal attribute identifier
		Code string
		// CodeLabel is the human-readable (perhaps localized) attribute name
		CodeLabel string
		// Label is the human-readable (perhaps localized) attribute value
		Label string
		// RawValue is the untouched original value of the attribute
		RawValue interface{}
		// UnitCode is the internal code of the attribute values measuring unit
		UnitCode string
	} // @name ProductAttribute

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

	// WishedToPay is the list of prices by type
	WishedToPay struct {
		priceByType map[string]priceDomain.Price
	}

	// Badge for a product
	// Example:
	// Badge {
	//   Code: "new",
	//   Label: "New Product",
	// }
	Badge struct {
		Code  string
		Label string
	}

	// Badges slice of Badge
	Badges []Badge

	// Stock holds data with product availability info
	Stock struct {
		InStock      bool
		Level        string
		Amount       int
		DeliveryCode string
	}
)

// Stock Level values
const (
	StockLevelOutOfStock = "out"
	StockLevelInStock    = "in"
	StockLevelLowStock   = "low"
)

// Value returns the raw value
func (at Attribute) Value() string {
	return strings.Trim(fmt.Sprintf("%v", at.RawValue), " ")
}

// IsEnabledValue returns true if the value can be seen as a toggle and is enabled
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

// IsDisabledValue returns true if the value can be seen as a disable toggle/switch value
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
	_, ok := at.RawValue.([]Attribute)
	if ok {
		return true
	}

	_, ok = at.RawValue.([]interface{})
	return ok
}

// Values builds a list of product attribute values in case the raw value is a slice
func (at Attribute) Values() []string {
	var result []string

	list, ok := at.RawValue.([]Attribute)
	if ok {
		for _, entry := range list {
			result = append(result, entry.Value())
		}
		return result
	}

	listFallback, ok := at.RawValue.([]interface{})
	if ok {
		for _, entry := range listFallback {
			result = append(result, strings.Trim(fmt.Sprintf("%v", entry), " "))
		}
	}
	return result
}

// Labels builds a list of human-readable product attribute values in case the raw value is a slice of Attribute, uses Values() as fallback
func (at Attribute) Labels() []string {
	var result []string
	list, ok := at.RawValue.([]Attribute)
	if ok {
		for _, entry := range list {
			result = append(result, entry.Label)
		}
		return result
	}

	return at.Values()
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

// HasAllAttributes returns true, if all attributes are set
func (bpd BasicProductData) HasAllAttributes(keys []string) bool {
	for _, key := range keys {
		if !bpd.HasAttribute(key) {
			return false
		}
	}
	return true
}

// Attribute get Attribute by key
func (bpd BasicProductData) Attribute(key string) Attribute {
	return bpd.Attributes[key]
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

// IsSaleableNow checks flag and time
func (p Saleable) IsSaleableNow() bool {
	if !p.IsSaleable {
		return false
	}

	// For some reasons IsZero does not always work - thats why we check for 1970
	if (p.SaleableFrom.IsZero() || p.SaleableFrom.Year() == 1970 || p.SaleableFrom.Before(time.Now())) &&
		(p.SaleableTo.IsZero() || p.SaleableTo.Year() == 1970 || p.SaleableTo.After(time.Now())) {
		return true
	}

	return false
}

// GetLoyaltyPriceByType returns the loyaltyentry that matches the type
func (p Saleable) GetLoyaltyPriceByType(ltype string) (*LoyaltyPriceInfo, bool) {
	for _, lp := range p.LoyaltyPrices {
		if lp.Type == ltype {
			return &lp, true
		}
	}
	return nil, false
}

// GetLoyaltyEarningByType returns the loyalty earning infos for a specific loyalty type
func (p Saleable) GetLoyaltyEarningByType(ltype string) (*LoyaltyEarningInfo, bool) {
	for _, le := range p.LoyaltyEarnings {
		if le.Type == ltype {
			return &le, true
		}
	}
	return nil, false
}

func (p Saleable) generateLoyaltyChargeSplit(valuedPriceToPay *priceDomain.Price, loyaltyPointsWishedToPay *WishedToPay, qty int, ignoreMin bool) priceDomain.Charges {
	if valuedPriceToPay == nil {
		finalPrice := p.ActivePrice.GetFinalPrice()
		valuedPriceToPay = &finalPrice
	}
	requiredCharges := make(map[string]priceDomain.Charge)
	remainingMainChargeValue := valuedPriceToPay.Amount()

	// getLoyaltyCharge - private func that returns the loyaltyCharge of the given type. making sure the currentlyRemainingMainChargeValue is not exceeded
	getValidLoyaltyCharge := func(loyaltyAmountWishedToSpent big.Float, loyaltyPrice LoyaltyPriceInfo, chargeType string, currentlyRemainingMainChargeValue *big.Float) priceDomain.Charge {
		loyaltyCurrency := loyaltyPrice.GetFinalPrice().Currency()
		rateLoyaltyFinalPriceToRealFinalPrice := loyaltyPrice.GetRate(p.ActivePrice.GetFinalPrice())
		maximumPossibleLoyaltyValue := big.NewFloat(0.0)
		if currentlyRemainingMainChargeValue.Cmp(big.NewFloat(0.0)) != 0 {
			maximumPossibleLoyaltyValue = new(big.Float).Quo(currentlyRemainingMainChargeValue, &rateLoyaltyFinalPriceToRealFinalPrice)
		}

		maximumPossibleLoyaltyPrice := priceDomain.NewFromBigFloat(*maximumPossibleLoyaltyValue, loyaltyCurrency).GetPayable()

		if loyaltyAmountWishedToSpent.Cmp(maximumPossibleLoyaltyValue) > 0 {
			loyaltyAmountWishedToSpent = *maximumPossibleLoyaltyValue
		}
		valuedLoyalityPrice := priceDomain.NewFromBigFloat(*new(big.Float).Mul(&rateLoyaltyFinalPriceToRealFinalPrice, &loyaltyAmountWishedToSpent), valuedPriceToPay.Currency()).GetPayable()
		if maximumPossibleLoyaltyPrice.Amount().Cmp(&loyaltyAmountWishedToSpent) == 0 {
			// If the wish equals the rounded maximum - we need to use the complete remaining value
			valuedLoyalityPrice = priceDomain.NewFromBigFloat(*currentlyRemainingMainChargeValue, valuedPriceToPay.Currency())
		}
		return priceDomain.Charge{
			Price: priceDomain.NewFromBigFloat(loyaltyAmountWishedToSpent, loyaltyCurrency).GetPayable(),
			Type:  chargeType,
			Value: valuedLoyalityPrice,
		}
	}

	for _, loyaltyPrice := range p.LoyaltyPrices {
		chargeType := loyaltyPrice.Type
		if chargeType == "" {
			continue
		}

		if !loyaltyPrice.GetFinalPrice().IsPositive() {
			continue
		}

		// loyaltyAmountToSpent - set as default without potential wish the minimum
		loyaltyAmountToSpent := loyaltyPrice.getMin(qty)

		// check if the minimum points should be ignored, if so minimum will be set to 0
		if ignoreMin {
			loyaltyAmountToSpent = *big.NewFloat(0.0)
		}

		if loyaltyPointsWishedToPay != nil {
			// if a loyaltyPointsWishedToPay is passed evaluate it within min and max and update loyaltyAmountToSpent:
			wishedPrice := loyaltyPointsWishedToPay.GetByType(chargeType)

			if wishedPrice != nil && wishedPrice.Currency() == loyaltyPrice.GetFinalPrice().Currency() {
				wishedPriceRounded := wishedPrice.GetPayable()

				// if wish is bigger than min we using the wish
				if loyaltyAmountToSpent.Cmp(wishedPriceRounded.Amount()) <= 0 {
					loyaltyAmountToSpent = *wishedPriceRounded.Amount()
				}
				// evaluate max
				max := loyaltyPrice.getMax(qty)
				if max != nil {
					// more then max - return max
					if max.Cmp(wishedPrice.Amount()) == -1 {
						loyaltyAmountToSpent = *max
					}
				}
			}
		}

		loyaltyCharge := getValidLoyaltyCharge(loyaltyAmountToSpent, loyaltyPrice, chargeType, remainingMainChargeValue)

		if !loyaltyCharge.Value.IsPositive() {
			continue
		}

		// Add the loyalty charge and at the same time reduce the remainingValue
		remainingMainChargeValue = new(big.Float).Sub(remainingMainChargeValue, loyaltyCharge.Value.Amount())
		requiredCharges[chargeType] = loyaltyCharge
	}

	remainingMainChargePrice := priceDomain.NewFromBigFloat(*remainingMainChargeValue, valuedPriceToPay.Currency()).GetPayable()
	requiredCharges[priceDomain.ChargeTypeMain] = priceDomain.Charge{
		Price: remainingMainChargePrice,
		Type:  priceDomain.ChargeTypeMain,
		Value: remainingMainChargePrice,
	}
	return *priceDomain.NewCharges(requiredCharges)
}

// GetLoyaltyChargeSplit gets the Charges that need to be paid by type:
// Type "main" is the remaining charge in the main currency and the other charges returned are the loyalty price charges that need to be paid.
// The method takes the min, max and the calculated loyalty conversion rate into account
//
// * valuedPriceToPay  Optional the price that need to be paid - if not given the products final price will be used
// * loyaltyPointsWishedToPay   Optional a list of loyaltyPrices that the (customer) wants to spend. Its used as a wish and may not be fulfilled because of min, max properties on the products loyaltyPrices
// * qty the quantity of the current item affects min max loyalty charge
func (p Saleable) GetLoyaltyChargeSplit(valuedPriceToPay *priceDomain.Price, loyaltyPointsWishedToPay *WishedToPay, qty int) priceDomain.Charges {
	return p.generateLoyaltyChargeSplit(valuedPriceToPay, loyaltyPointsWishedToPay, qty, false)
}

// GetLoyaltyChargeSplitIgnoreMin same as GetLoyaltyChargeSplit but ignoring the min points to spend
func (p Saleable) GetLoyaltyChargeSplitIgnoreMin(valuedPriceToPay *priceDomain.Price, loyaltyPointsWishedToPay *WishedToPay, qty int) priceDomain.Charges {
	return p.generateLoyaltyChargeSplit(valuedPriceToPay, loyaltyPointsWishedToPay, qty, true)
}

// getMin returns minimum points to spend - scaled by qty
func (l LoyaltyPriceInfo) getMin(qty int) big.Float {
	return *new(big.Float).Mul(&l.MinPointsToSpent, big.NewFloat(float64(qty)))
}

// getMax returns max points to spend - scaled by qty. If no max set returns nil
func (l LoyaltyPriceInfo) getMax(qty int) *big.Float {
	if !l.HasMax() {
		return nil
	}
	return new(big.Float).Mul(l.MaxPointsToSpent, big.NewFloat(float64(qty)))
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
	for _, stock := range bpd.Stock {
		if stock.InStock {
			return true
		}
	}

	if bpd.StockLevel == "" || bpd.StockLevel == StockLevelOutOfStock {
		return false
	}

	return true
}

// IsInStockForDeliveryCode returns information if current product whether in stock or not for provided delivery code
func (bpd BasicProductData) IsInStockForDeliveryCode(deliveryCode string) bool {
	for _, stock := range bpd.Stock {
		if stock.DeliveryCode == deliveryCode {
			return stock.InStock
		}
	}

	return false
}

// NewWishedToPay returns a new WishedToPay struct
func NewWishedToPay() WishedToPay {
	return WishedToPay{
		priceByType: make(map[string]priceDomain.Price),
	}
}

// Add returns new WishedToPay instance with the given wish added
func (w WishedToPay) Add(ctype string, price priceDomain.Price) WishedToPay {
	if w.priceByType == nil {
		w.priceByType = make(map[string]priceDomain.Price)
	}
	w.priceByType[ctype] = price
	return w
}

// GetByType returns the wished price for the given type or nil
func (w WishedToPay) GetByType(ctype string) *priceDomain.Price {
	if price, ok := w.priceByType[ctype]; ok {
		return &price
	}
	return nil
}

// GetRate returns the currency conversion rate of the current loyaltyprice final price in relation to the passed value
func (l LoyaltyPriceInfo) GetRate(valuedPrice priceDomain.Price) big.Float {
	if !l.GetFinalPrice().IsPositive() {
		return *big.NewFloat(0)
	}
	return *new(big.Float).Quo(valuedPrice.Amount(), l.GetFinalPrice().Amount())
}

// HasMax checks if product has a maximum (points to spend) restriction
func (l LoyaltyPriceInfo) HasMax() bool {
	return l.MaxPointsToSpent != nil
}

// GetFinalPrice gets either the Default or the Discounted Loyaltyprice
func (l LoyaltyPriceInfo) GetFinalPrice() priceDomain.Price {
	if l.IsDiscounted && l.Discounted.IsLessThen(l.Default) {
		return l.Discounted
	}
	return l.Default
}

// Split splits the given WishedToPay in payable childs
func (w WishedToPay) Split(count int) []WishedToPay {
	// init slice
	result := make([]WishedToPay, count)
	for k := range result {
		result[k] = NewWishedToPay()
	}
	// fill slice with splitted
	for chargeType, v := range w.priceByType {
		valuesSplitted, _ := v.SplitInPayables(count)
		for i, splittedValue := range valuesSplitted {
			result[i] = result[i].Add(chargeType, splittedValue)
		}
	}
	return result
}

// AttributeKeys lists all available keys
func (a Attributes) AttributeKeys() []string {
	res := make([]string, len(a))
	i := 0
	for k := range a {
		res[i] = k
		i++
	}
	return res
}

// Attributes lists all attributes
func (a Attributes) Attributes() []Attribute {
	res := make([]Attribute, len(a))
	i := 0
	for _, v := range a {
		res[i] = v
		i++
	}
	return res
}

// HasAttribute checks if an attribute is available
func (a Attributes) HasAttribute(key string) bool {
	_, exist := a[key]
	return exist
}

// Attribute returns a specified attribute
func (a Attributes) Attribute(key string) Attribute {
	attribute := a[key]
	return attribute
}

// AttributesByKey returns slice of attributes by given attribute keys
func (a Attributes) AttributesByKey(keys []string) []Attribute {
	res := make([]Attribute, 0)
	for _, key := range keys {
		if a.HasAttribute(key) {
			res = append(res, a.Attribute(key))
		}
	}

	return res
}

// CPath returns the constructed Path from this category to the root - splitted by /
func (c *CategoryTeaser) CPath() string {
	if c.Parent == nil || c.Parent == c {
		return c.Code
	}
	return c.Parent.CPath() + "/" + c.Code
}

// First of the badges, returns nil if there is no first badge
func (b Badges) First() *Badge {
	if len(b) == 0 {
		return nil
	}

	result := b[0]
	return &result
}
