package cart

import (
	"sort"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// Item for Cart
	Item struct {
		// ID of the item - needs to be unique over the whole cart
		ID string
		// ExternalReference can be used by cart service implementations to separate the representation in an external
		// cart service from the unique item ID
		ExternalReference string
		// MarketplaceCode is the identifier for the product
		MarketplaceCode string
		// VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceID string

		Qty int

		AdditionalData map[string]string

		// SinglePriceGross is the gross price (incl. taxes) for a single product
		SinglePriceGross priceDomain.Price

		// SinglePriceNet is the net price (excl. taxes) for a single product
		SinglePriceNet priceDomain.Price

		// RowPriceGross is the price incl. taxes for the whole Qty of products
		RowPriceGross priceDomain.Price

		// RowPriceGrossWithDiscount is the price incl. taxes with deducted discounts for the whole Qty of products
		// This is in most cases the final price for the customer to pay
		RowPriceGrossWithDiscount priceDomain.Price

		// RowPriceGrossWithItemRelatedDiscount is the price incl. taxes with deducted item related discounts for the whole Qty of products
		RowPriceGrossWithItemRelatedDiscount priceDomain.Price

		// RowPriceNet is the price excl. taxes for the whole Qty of products for the whole Qty of products
		RowPriceNet priceDomain.Price

		// RowPriceNetWithDiscount is the discounted net price for the whole Qty of products
		RowPriceNetWithDiscount priceDomain.Price

		// RowPriceNetWithItemRelatedDiscount is the price excl. taxes with deducted item related discounts for the whole Qty of products
		RowPriceNetWithItemRelatedDiscount priceDomain.Price

		// RowTaxes is a list of all taxes applied for the given Qty of products
		RowTaxes Taxes

		// TotalTaxAmount is the sum of all applied taxes for the whole Qty of products
		TotalTaxAmount priceDomain.Price

		// AppliedDiscounts contains the details about the discounts applied to this item - they can be "itemrelated" or not
		// itemrelated would be e.g. special price, buy 3 pay 2
		// non-itemrelated would be e.g. 10% on everything
		AppliedDiscounts AppliedDiscounts

		// TotalDiscountAmount is the sum of all applied discounts (aka the savings for the customer)
		TotalDiscountAmount priceDomain.Price

		// ItemRelatedDiscountAmount is the sum of all itemrelated Discounts
		ItemRelatedDiscountAmount priceDomain.Price

		// NonItemRelatedDiscountAmount is the sum of non-itemrelated Discounts where IsItemRelated = false
		NonItemRelatedDiscountAmount priceDomain.Price
	}

	// ItemSplitter used to split an item
	ItemSplitter struct {
		itemBuilderProvider ItemBuilderProvider
		configUseGrossPrice bool
	}
)

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

// Inject dependencies
func (s *ItemSplitter) Inject(itemBuilderProvider ItemBuilderProvider, config *struct {
	UseGrossPrice bool `inject:"config:commerce.product.priceIsGross,optional"`
}) {
	s.itemBuilderProvider = itemBuilderProvider
	if config != nil {
		s.configUseGrossPrice = config.UseGrossPrice
	}
}

// SplitInSingleQtyItems the given item into multiple items with Qty 1 and make sure the sum of the items prices matches by using SplitInPayables
func (s *ItemSplitter) SplitInSingleQtyItems(givenItem Item) ([]Item, error) {
	var items []Item
	// configUseGrossPrice true then:
	// Given: SinglePriceGross / all AppliedDiscounts  / All Taxes
	// Calculated: SinglePriceNet / RowPriceGross / RowPriceNet / SinglePriceNet

	// configUseGrossPrice false then:
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
			// The split adds the moving cents to the first ones, resulting in
			// having the smallest prices at the end but since discounts are
			// negative, we need to reverse it to ensure that a split of the row
			// totals has the rounding cents at the same positions
			sort.Slice(apSplitted, func(i, j int) bool {
				return apSplitted[i].FloatAmount() > apSplitted[j].FloatAmount()
			})
			p := make([]float64, 0)
			for _, i := range apSplitted {
				p = append(p, i.FloatAmount())
			}
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
			itemBuilder.SetSinglePriceGross(givenItem.SinglePriceGross.GetPayable())
			itemBuilder.CalculatePricesAndTaxAmountsFromSinglePriceGross()
		} else {
			itemBuilder.SetSinglePriceNet(givenItem.SinglePriceNet.GetPayable())
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
