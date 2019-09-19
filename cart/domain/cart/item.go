package cart

import (
	"errors"
	"fmt"
	"math/big"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// Item for Cart
	Item struct {
		// ID of the item - needs to be unique over the whole cart
		ID string
		// ExternalReference can be used by cart service implementations to separate the representation in an external
		// cart service from the unique item ID
		ExternalReference string
		MarketplaceCode   string
		// VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceID string

		Qty int

		AdditionalData map[string]string

		// SinglePriceGross brutto (gross) for single product
		SinglePriceGross priceDomain.Price

		// SinglePriceNet net price for single product
		SinglePriceNet priceDomain.Price

		// RowPriceGross
		RowPriceGross priceDomain.Price

		// RowPriceNet
		RowPriceNet priceDomain.Price

		// RowPriceGross
		RowTaxes Taxes

		// AppliedDiscounts contains the details about the discounts applied to this item - they can be "itemrelated" or not
		AppliedDiscounts AppliedDiscounts
	}

	// ItemBuilder can be used to construct an item with a fluent interface
	ItemBuilder struct {
		itemCurrency        *string
		invariantError      error
		itemInBuilding      *Item
		configUseGrossPrice bool
	}

	// ItemBuilderProvider should be used to create an item
	ItemBuilderProvider func() *ItemBuilder

	// ItemSplitter used to split an item
	ItemSplitter struct {
		itemBuilderProvider ItemBuilderProvider
		configUseGrossPrice bool
	}
)

// TotalTaxAmount - returns total tax amount as price
func (i Item) TotalTaxAmount() priceDomain.Price {
	return i.RowTaxes.TotalAmount()
}

// TotalDiscountAmount gets the savings by item
func (i Item) TotalDiscountAmount() priceDomain.Price {
	result, _ := i.NonItemRelatedDiscountAmount().Add(i.ItemRelatedDiscountAmount())

	return result
}

// ItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = True
func (i Item) ItemRelatedDiscountAmount() priceDomain.Price {
	prices := make([]priceDomain.Price, 0, len(i.AppliedDiscounts))

	for _, discount := range i.AppliedDiscounts {
		if !discount.IsItemRelated {
			continue
		}
		prices = append(prices, discount.Applied)
	}

	result, _ := priceDomain.SumAll(prices...)

	return result.GetPayable()
}

// NonItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false
func (i Item) NonItemRelatedDiscountAmount() priceDomain.Price {
	prices := make([]priceDomain.Price, 0, len(i.AppliedDiscounts))

	for _, discount := range i.AppliedDiscounts {
		if discount.IsItemRelated {
			continue
		}
		prices = append(prices, discount.Applied)
	}

	result, _ := priceDomain.SumAll(prices...)

	return result.GetPayable()
}

// RowPriceGrossWithDiscount = RowPriceGross-TotalDiscountAmount()
func (i Item) RowPriceGrossWithDiscount() priceDomain.Price {
	result, _ := i.RowPriceGross.Add(i.TotalDiscountAmount())

	return result
}

// RowPriceNetWithDiscount = RowPriceNet-TotalDiscountAmount()
func (i Item) RowPriceNetWithDiscount() priceDomain.Price {
	result, _ := i.RowPriceNet.Add(i.TotalDiscountAmount())

	return result
}

// RowPriceGrossWithItemRelatedDiscount = RowPriceGross-ItemRelatedDiscountAmount()
func (i Item) RowPriceGrossWithItemRelatedDiscount() priceDomain.Price {
	result, _ := i.RowPriceGross.Add(i.ItemRelatedDiscountAmount())

	return result
}

// RowPriceNetWithItemRelatedDiscount =RowTotal-ItemRelatedDiscountAmount
func (i Item) RowPriceNetWithItemRelatedDiscount() priceDomain.Price {
	result, _ := i.RowPriceNet.Add(i.ItemRelatedDiscountAmount())

	return result
}

// AdditionalDataKeys lists all available keys
func (i Item) AdditionalDataKeys() []string {
	additionalData := i.AdditionalData
	res := make([]string, len(additionalData))
	n := 0
	for k := range additionalData {
		res[n] = k
		n++
	}
	return res
}

// AdditionalDataValues lists all values
func (i Item) AdditionalDataValues() []string {
	res := make([]string, len(i.AdditionalData))
	n := 0
	for _, v := range i.AdditionalData {
		res[n] = v
		n++
	}
	return res
}

// HasAdditionalDataKey checks if an attribute is available
func (i Item) HasAdditionalDataKey(key string) bool {
	_, exist := i.AdditionalData[key]
	return exist
}

// GetAdditionalData returns a specified attribute
func (i Item) GetAdditionalData(key string) string {
	attribute := i.AdditionalData[key]
	return attribute
}

// Inject - called by dingo
func (f *ItemBuilder) Inject(config *struct {
	UseGrosPrice bool `inject:"config:commerce.product.priceIsGross,optional"`
}) {
	if config != nil {
		f.configUseGrossPrice = config.UseGrosPrice
	}
}

// SetID - sets the id
func (f *ItemBuilder) SetID(id string) *ItemBuilder {
	f.init()
	f.itemInBuilding.ID = id
	return f
}

// SetExternalReference - sets the ExternalReference
func (f *ItemBuilder) SetExternalReference(ref string) *ItemBuilder {
	f.init()
	f.itemInBuilding.ExternalReference = ref
	return f
}

// SetFromItem - sets the data in builder from existing item - useful to get a updated item based from existing. Its not setting Taxes (use Calculate)
func (f *ItemBuilder) SetFromItem(item Item) *ItemBuilder {
	f.init()
	f.SetProductData(item.MarketplaceCode, item.VariantMarketPlaceCode, item.ProductName)
	f.SetExternalReference(item.ExternalReference)
	f.SetID(item.ID)
	f.SetQty(item.Qty)
	f.AddDiscounts(item.AppliedDiscounts...)
	f.SetSinglePriceGross(item.SinglePriceGross)
	f.SetSinglePriceNet(item.SinglePriceNet)

	return f
}

// SetVariantMarketPlaceCode sets VariantMarketPlaceCode (only for configurable_with_variant relevant)
func (f *ItemBuilder) SetVariantMarketPlaceCode(id string) *ItemBuilder {
	f.init()
	f.itemInBuilding.VariantMarketPlaceCode = id
	return f
}

// SetSourceID - optional
func (f *ItemBuilder) SetSourceID(id string) *ItemBuilder {
	f.init()
	f.itemInBuilding.SourceID = id
	return f
}

// SetAdditionalData - optional
func (f *ItemBuilder) SetAdditionalData(d map[string]string) *ItemBuilder {
	f.init()
	f.itemInBuilding.AdditionalData = d
	return f
}

// SetQty - optional (default 1)
func (f *ItemBuilder) SetQty(q int) *ItemBuilder {
	f.init()
	f.itemInBuilding.Qty = q
	return f
}

// SetSinglePriceGross - set by gross price
func (f *ItemBuilder) SetSinglePriceGross(grossPrice priceDomain.Price) *ItemBuilder {
	f.init()
	if !grossPrice.IsPayable() {
		f.invariantError = errors.New("SetSinglePriceGross need to get payable price")
	}
	f.itemInBuilding.SinglePriceGross = grossPrice
	f.checkCurrency(&grossPrice)
	return f
}

// SetSinglePriceNet - set by net
func (f *ItemBuilder) SetSinglePriceNet(price priceDomain.Price) *ItemBuilder {
	f.init()
	if !price.IsPayable() {
		f.invariantError = errors.New("SetSinglePriceNet need to get payable price")
	}
	f.itemInBuilding.SinglePriceNet = price
	f.checkCurrency(&price)
	return f
}

// AddTaxInfo - add a tax info - at least taxRate OR taxAmount need to be given. The taxAmount can be calculated
func (f *ItemBuilder) AddTaxInfo(taxType string, taxRate *big.Float, taxAmount *priceDomain.Price) *ItemBuilder {
	f.init()
	if taxRate == nil && taxAmount == nil {
		f.invariantError = errors.New("at least taxRate or taxAmount need to be given")
	}
	tax := Tax{
		Type: taxType,
		Rate: taxRate,
	}
	if taxAmount != nil {
		if !taxAmount.IsPayable() {
			f.invariantError = errors.New("taxAmount need to be payable price")
		}
		f.checkCurrency(taxAmount)
		tax.Amount = *taxAmount
	}
	f.itemInBuilding.RowTaxes = append(f.itemInBuilding.RowTaxes, tax)
	return f
}

// AddDiscount - adds a discount
func (f *ItemBuilder) AddDiscount(discount AppliedDiscount) *ItemBuilder {
	f.init()
	if !discount.Applied.IsPayable() {
		f.invariantError = errors.New("AddDiscount need to have payable price")
	}
	if !discount.Applied.IsNegative() {
		f.invariantError = fmt.Errorf("AddDiscount need to have negative price - given %f", discount.Applied.FloatAmount())
	}
	f.checkCurrency(&discount.Applied)
	f.itemInBuilding.AppliedDiscounts = append(f.itemInBuilding.AppliedDiscounts, discount)
	return f
}

// AddDiscounts - adds a discount
func (f *ItemBuilder) AddDiscounts(discounts ...AppliedDiscount) *ItemBuilder {
	for _, discount := range discounts {
		f.AddDiscount(discount)
	}
	return f
}

// CalculatePricesAndTaxAmountsFromSinglePriceNet - Vertikal Tax Calculation - based from current SinglePriceNet, Qty and the RowTax Infos given
// Sets RowPriceNet, missing tax.Amount and RowPriceGross
func (f *ItemBuilder) CalculatePricesAndTaxAmountsFromSinglePriceNet() *ItemBuilder {
	priceNet := f.itemInBuilding.SinglePriceNet
	f.itemInBuilding.RowPriceNet = priceNet.Multiply(f.itemInBuilding.Qty)
	for k, tax := range f.itemInBuilding.RowTaxes {
		// Calculate tax amount from rate if required
		if tax.Amount.IsZero() && tax.Rate != nil {
			// set tax amount and round it
			tax.Amount = f.itemInBuilding.RowPriceNetWithDiscount().TaxFromNet(*tax.Rate).GetPayable()
			f.itemInBuilding.RowTaxes[k] = tax
		}
	}
	totalTaxAmount := f.itemInBuilding.TotalTaxAmount()
	f.itemInBuilding.RowPriceGross, _ = priceDomain.SumAll(f.itemInBuilding.RowPriceNet, totalTaxAmount)
	if f.itemInBuilding.Qty == 0 {
		f.invariantError = errors.New("Quantity is Zero")
		return f
	}
	f.itemInBuilding.SinglePriceGross, _ = priceNet.Add(totalTaxAmount.Divided(f.itemInBuilding.Qty))
	return f
}

// CalculatePricesAndTax - reads the config flag and reculculates Total and Tax
func (f *ItemBuilder) CalculatePricesAndTax() *ItemBuilder {
	if f.configUseGrossPrice {
		return f.CalculatePricesAndTaxAmountsFromSinglePriceGross()
	}
	return f.CalculatePricesAndTaxAmountsFromSinglePriceNet()
}

// CalculatePricesAndTaxAmountsFromSinglePriceGross - Vertical Tax Calculation - based from current SinglePriceNet, Qty and the RowTax Infos given
// Sets RowPriceNet, missing tax.Amount and RowPriceGross
func (f *ItemBuilder) CalculatePricesAndTaxAmountsFromSinglePriceGross() *ItemBuilder {
	priceGross := f.itemInBuilding.SinglePriceGross
	f.itemInBuilding.RowPriceGross = priceGross.Multiply(f.itemInBuilding.Qty)
	for k, tax := range f.itemInBuilding.RowTaxes {
		// Calculate tax amount from rate if required
		if tax.Amount.IsZero() && tax.Rate != nil {
			tax.Amount = f.itemInBuilding.RowPriceGrossWithDiscount().TaxFromGross(*tax.Rate).GetPayable()
			f.itemInBuilding.RowTaxes[k] = tax
		}
	}
	totalTaxAmount := f.itemInBuilding.TotalTaxAmount()
	f.itemInBuilding.RowPriceNet, _ = f.itemInBuilding.RowPriceGross.Sub(totalTaxAmount)
	f.itemInBuilding.SinglePriceNet, _ = priceGross.Sub(totalTaxAmount.Divided(f.itemInBuilding.Qty))
	return f
}

// SetProductData - set product data: MarketplaceCode,VariantMarketPlaceCode,ProductName
func (f *ItemBuilder) SetProductData(marketplace string, vc string, name string) *ItemBuilder {
	f.init()
	f.itemInBuilding.MarketplaceCode = marketplace
	f.itemInBuilding.VariantMarketPlaceCode = vc
	f.itemInBuilding.ProductName = name
	return f
}

// SetByProduct - gets a product and calculates also prices
func (f *ItemBuilder) SetByProduct(product domain.BasicProduct) *ItemBuilder {
	if !product.IsSaleable() {
		f.invariantError = errors.New("Product is not saleable")
	}

	f.init()
	f.itemInBuilding.MarketplaceCode = product.BaseData().MarketPlaceCode
	f.itemInBuilding.ProductName = product.BaseData().Title

	if configurable, ok := product.(domain.ConfigurableProductWithActiveVariant); ok {
		f.itemInBuilding.VariantMarketPlaceCode = configurable.ActiveVariant.MarketPlaceCode

	}

	if f.configUseGrossPrice {
		f.SetSinglePriceGross(product.SaleableData().ActivePrice.GetFinalPrice())
		f.CalculatePricesAndTaxAmountsFromSinglePriceGross()
	} else {
		f.SetSinglePriceNet(product.SaleableData().ActivePrice.GetFinalPrice())
		f.CalculatePricesAndTaxAmountsFromSinglePriceNet()
	}

	return f
}

func (f *ItemBuilder) checkCurrency(price *priceDomain.Price) {
	if price == nil {
		return
	}
	currency := price.Currency()
	if f.itemCurrency == nil {
		f.itemCurrency = &currency
		return
	}
	if *f.itemCurrency != currency {
		f.invariantError = fmt.Errorf("There is a currency mismatch inside the item %v and %v", currency, *f.itemCurrency)
	}
}

// Build returns build item or error if invariants do not match. Any call will also REST the ItemBuilder
func (f *ItemBuilder) Build() (*Item, error) {
	if f.itemInBuilding == nil {
		return f.reset(errors.New("Nothing in building"))
	}

	if f.invariantError != nil {
		return f.reset(f.invariantError)
	}

	if f.itemInBuilding.ID == "" {
		return f.reset(errors.New("Id Required"))
	}

	checkPrice, _ := f.itemInBuilding.RowPriceNet.Add(f.itemInBuilding.TotalTaxAmount())
	if !checkPrice.LikelyEqual(f.itemInBuilding.RowPriceGross) {
		return f.reset(fmt.Errorf("RowPriceGross (%f) need to match likely TotalTaxAmount + RowPriceNet. (%f) for item %v ", f.itemInBuilding.RowPriceGross.FloatAmount(), checkPrice.FloatAmount(), f.itemInBuilding.ID))
	}

	return f.reset(nil)
}

func (f *ItemBuilder) init() {
	if f.itemInBuilding == nil {
		f.itemInBuilding = &Item{
			Qty: 1,
		}
	}
}

func (f *ItemBuilder) reset(err error) (*Item, error) {
	item := f.itemInBuilding
	f.itemInBuilding = nil
	f.invariantError = nil
	f.itemCurrency = nil
	return item, err
}

func (s *ItemSplitter) Inject(itemBuilderProvider ItemBuilderProvider, config *struct {
	UseGrosPrice bool `inject:"config:commerce.product.priceIsGross,optional"`
}) {
	s.itemBuilderProvider = itemBuilderProvider
	if config != nil {
		s.configUseGrossPrice = config.UseGrosPrice
	}
}

//Split the given item into multiple items with Qty 1 and make sure the sum of the items prices matches by using SplitInPayables
func (s *ItemSplitter) Split(givenItem Item) ([]Item, error) {
	var items []Item
	//configUseGrossPrice true then:
	// Given: SinglePriceGross / all AppliedDiscounts  / All Taxes
	// Calculated: SinglePriceNet / RowPriceGross / RowPriceNet / SinglePriceNet

	//configUseGrossPrice false then:
	// Given: SinglePriceNez / all AppliedDiscounts  / All Taxes
	// Calculated: SinglePriceGross / RowPriceGross / RowPriceNet / SinglePriceGross
	for x := 0; x < givenItem.Qty; x++ {

		itemBuilder := s.itemBuilderProvider()
		itemBuilder.SetProductData(givenItem.MarketplaceCode, givenItem.VariantMarketPlaceCode, givenItem.ProductName)
		itemBuilder.SetExternalReference(givenItem.ExternalReference)
		itemBuilder.SetID(givenItem.ID)
		itemBuilder.SetQty(1)
		for _, ap := range givenItem.AppliedDiscounts {
			apSplitted, err := ap.Applied.SplitInPayables(givenItem.Qty)
			if err != nil {
				return nil, err
			}
			newDiscount := AppliedDiscount{
				CampaignCode:  ap.CampaignCode,
				CouponCode:    ap.CouponCode,
				Label:         ap.Label,
				Applied:       apSplitted[x],
				Type:          ap.Type,
				IsItemRelated: ap.IsItemRelated,
				SortOrder:     ap.SortOrder,
			}
			itemBuilder.AddDiscount(newDiscount)
		}
		for _, rt := range givenItem.RowTaxes {
			if rt.Amount.IsZero() {
				continue
			}
			rtSplitted, err := rt.Amount.SplitInPayables(givenItem.Qty)
			if err != nil {
				return nil, err
			}
			itemBuilder.AddTaxInfo(rt.Type, rt.Rate, &rtSplitted[x])
		}
		if s.configUseGrossPrice {
			itemBuilder.SetSinglePriceGross(givenItem.SinglePriceGross)
			itemBuilder.CalculatePricesAndTaxAmountsFromSinglePriceGross()
		} else {
			itemBuilder.SetSinglePriceNet(givenItem.SinglePriceGross)
			itemBuilder.CalculatePricesAndTaxAmountsFromSinglePriceNet()
		}
		item, err := itemBuilder.Build()
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return items, nil
}
