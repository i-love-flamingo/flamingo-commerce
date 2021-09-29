package cart_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestItemSplitter_SplitGrossBased(t *testing.T) {
	// todo
	splitter := &cartDomain.ItemSplitter{}

	item := &cartDomain.Item{
		ID:                                   "2",
		ExternalReference:                    "external",
		MarketplaceCode:                      "marketplace",
		VariantMarketPlaceCode:               "variant",
		ProductName:                          "product",
		SourceID:                             "warehouse1",
		Qty:                                  5,
		SinglePriceGross:                     priceDomain.NewFromInt(2065, 100, "€"),
		SinglePriceNet:                       priceDomain.NewFromInt(1930, 100, "€"),
		RowPriceGross:                        priceDomain.NewFromInt(2065*5, 100, "€"),
		RowPriceGrossWithDiscount:            priceDomain.NewFromInt(2065*5-3172, 100, "€"),
		RowPriceGrossWithItemRelatedDiscount: priceDomain.NewFromInt(2065*5, 100, "€"),
		RowPriceNet:                          priceDomain.NewFromInt(1930*5, 100, "€"),
		RowPriceNetWithDiscount:              priceDomain.NewFromInt(6685, 100, "€"),
		RowPriceNetWithItemRelatedDiscount:   priceDomain.NewFromInt(1930*5, 100, "€"),
		RowTaxes:                             []cartDomain.Tax{{Amount: priceDomain.NewFromInt(468, 100, "€"), Type: "tax", Rate: big.NewFloat(7)}},
		AppliedDiscounts: cartDomain.AppliedDiscounts{
			cartDomain.AppliedDiscount{IsItemRelated: false, Applied: priceDomain.NewFromInt(-3172, 100, "€")},
		},
		TotalDiscountAmount:          priceDomain.NewFromInt(-3172, 100, "€"),
		ItemRelatedDiscountAmount:    priceDomain.NewFromInt(-3172, 100, "€"),
		NonItemRelatedDiscountAmount: priceDomain.NewFromInt(0, 100, "€"),
	}

	splittedItems, err := splitter.SplitInSingleQtyItems(*item)
	require.NoError(t, err)

	assert.Equal(t, 4.68, item.TotalTaxAmount().FloatAmount())

	var discount, rowGrossTotal, rowNetTotal, totalTaxAmount, totalDiscountAmount float64
	for _, splitItem := range splittedItems {
		// todo: check all meta fields, e.g. source location
		assert.Equal(t, 1, splitItem.Qty)
		// make sure single and row price are equal:
		assert.Equal(t, splitItem.SinglePriceNet.FloatAmount(), splitItem.RowPriceNet.FloatAmount())
		assert.Equal(t, splitItem.SinglePriceGross.FloatAmount(), splitItem.RowPriceGross.FloatAmount())
		// make sure it's consistent (net+tax=gross):
		assert.Equal(t, splitItem.RowPriceGrossWithDiscount.FloatAmount(), splitItem.RowPriceNetWithDiscount.ForceAdd(splitItem.TotalTaxAmount()).FloatAmount())

		rowGrossTotal = rowGrossTotal + splitItem.RowPriceGross.FloatAmount()
		rowNetTotal = rowNetTotal + splitItem.RowPriceNet.FloatAmount()
		totalTaxAmount = totalTaxAmount + splitItem.TotalTaxAmount().FloatAmount()
		rate, _ := splitItem.RowTaxes[0].Rate.Float64()
		assert.Equal(t, 7.0, rate)
		totalDiscountAmount = totalDiscountAmount + splitItem.TotalDiscountAmount.FloatAmount()

		// discount split cents should be at the end, so the next discount must be the same or smaller
		assert.GreaterOrEqual(t, discount, splitItem.TotalDiscountAmount.FloatAmount())
		discount = splitItem.TotalDiscountAmount.FloatAmount()
	}
	//todo check the sum of all fields
	assert.Equal(t, item.RowPriceGross.FloatAmount(), rowGrossTotal)
	assert.Equal(t, item.RowPriceNet.FloatAmount(), rowNetTotal)
	assert.Equal(t, item.TotalTaxAmount().FloatAmount(), totalTaxAmount)
	assert.Equal(t, item.TotalDiscountAmount.FloatAmount(), totalDiscountAmount)
}

func TestItem_AdditionalDataKeys(t *testing.T) {
	// todo
}

func TestItem_AdditionalDataValues(t *testing.T) {
	// todo
}

func TestItem_HasAdditionalDataKey(t *testing.T) {
	// todo
}
func TestItem_GetAdditionalData(t *testing.T) {
	// todo
}
