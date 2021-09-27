package cart

import (
	"errors"
	"fmt"
	"math/big"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// ItemBuilder can be used to construct an item with a fluent interface
	ItemBuilder struct {
		itemCurrency        *string
		invariantError      error
		itemInBuilding      *Item
		configUseGrossPrice bool
	}

	// ItemBuilderProvider should be used to create an item
	ItemBuilderProvider func() *ItemBuilder
)

// Inject dependencies
func (f *ItemBuilder) Inject(config *struct {
	UseGrossPrice bool `inject:"config:commerce.product.priceIsGross,optional"`
}) {
	if config != nil {
		f.configUseGrossPrice = config.UseGrossPrice
	}
}

// SetID sets the id
func (f *ItemBuilder) SetID(id string) *ItemBuilder {
	f.init()
	f.itemInBuilding.ID = id
	return f
}

// SetExternalReference sets the ExternalReference
func (f *ItemBuilder) SetExternalReference(ref string) *ItemBuilder {
	f.init()
	f.itemInBuilding.ExternalReference = ref
	return f
}

// SetFromItem sets the data in builder from existing item - useful to get a updated item based from existing. Its not setting Taxes (use Calculate)
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

// SetSourceID sets optional source ID
func (f *ItemBuilder) SetSourceID(id string) *ItemBuilder {
	f.init()
	f.itemInBuilding.SourceID = id
	return f
}

// SetAdditionalData sets optional additional data
func (f *ItemBuilder) SetAdditionalData(d map[string]string) *ItemBuilder {
	f.init()
	f.itemInBuilding.AdditionalData = d
	return f
}

// SetQty sets the qty (defaults to 1)
func (f *ItemBuilder) SetQty(q int) *ItemBuilder {
	f.init()
	f.itemInBuilding.Qty = q
	return f
}

// SetSinglePriceGross set by gross price
func (f *ItemBuilder) SetSinglePriceGross(grossPrice priceDomain.Price) *ItemBuilder {
	f.init()
	if !grossPrice.IsPayable() {
		f.invariantError = errors.New("SetSinglePriceGross need to get payable price")
	}
	f.itemInBuilding.SinglePriceGross = grossPrice
	f.checkCurrency(&grossPrice)
	return f
}

// SetSinglePriceNet set by net price
func (f *ItemBuilder) SetSinglePriceNet(price priceDomain.Price) *ItemBuilder {
	f.init()
	if !price.IsPayable() {
		f.invariantError = errors.New("SetSinglePriceNet need to get payable price")
	}
	f.itemInBuilding.SinglePriceNet = price
	f.checkCurrency(&price)
	return f
}

// AddTaxInfo add a tax info - at least taxRate OR taxAmount need to be given. the tax amount can be calculated
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

// AddDiscount adds a discount
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

// AddDiscounts adds a list of discounts
func (f *ItemBuilder) AddDiscounts(discounts ...AppliedDiscount) *ItemBuilder {
	for _, discount := range discounts {
		f.AddDiscount(discount)
	}
	return f
}

// SetProductData set product data: MarketplaceCode, VariantMarketPlaceCode, ProductName
func (f *ItemBuilder) SetProductData(marketplace string, vc string, name string) *ItemBuilder {
	f.init()
	f.itemInBuilding.MarketplaceCode = marketplace
	f.itemInBuilding.VariantMarketPlaceCode = vc
	f.itemInBuilding.ProductName = name
	return f
}

// SetByProduct gets a product and calculates also prices
func (f *ItemBuilder) SetByProduct(product domain.BasicProduct) *ItemBuilder {
	if !product.IsSaleable() {
		f.invariantError = errors.New("Product is not saleable")
	}

	f.init()
	f.itemInBuilding.MarketplaceCode = product.BaseData().MarketPlaceCode
	f.itemInBuilding.ProductName = product.BaseData().Title

	if configurable, ok := product.(domain.ConfigurableProductWithActiveVariant); ok {
		f.itemInBuilding.MarketplaceCode = configurable.ConfigurableBaseData().MarketPlaceCode
		f.itemInBuilding.VariantMarketPlaceCode = configurable.ActiveVariant.MarketPlaceCode
	}

	if f.configUseGrossPrice {
		f.SetSinglePriceGross(product.SaleableData().ActivePrice.GetFinalPrice())
	} else {
		f.SetSinglePriceNet(product.SaleableData().ActivePrice.GetFinalPrice())
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

	checkPrice, _ := f.itemInBuilding.RowPriceNet.Add(f.itemInBuilding.TotalTaxAmount)
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
