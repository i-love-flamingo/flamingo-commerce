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

	builder := provider()
	builder.SetSinglePriceGross(priceDomain.NewFromInt(2065, 100, "€")).
		SetQty(5).AddTaxInfo("tax", big.NewFloat(7), nil).
		SetID("2").
		AddDiscount(cartDomain.AppliedDiscount{Applied: priceDomain.NewFromInt(-3172, 100, "€")}).
		CalculatePricesAndTaxAmountsFromSinglePriceGross()
	item, err := builder.Build()
	require.NoError(t, err)

	splittedItems, err := splitter.SplitInSingleQtyItems(*item)
	require.NoError(t, err)

	// 20.65 * 5 = 103.25
	assert.Equal(t, 103.25, item.RowPriceGross.FloatAmount())
	assert.Equal(t, -31.72, item.TotalDiscountAmount().FloatAmount())
	// (98.57 - 31.70) * 0.07
	assert.Equal(t, 4.68, item.TotalTaxAmount().FloatAmount())
	// TotalTaxAmount + 98.57 = 103.25
	assert.Equal(t, 98.57, item.RowPriceNet.FloatAmount())
	assert.Equal(t, 66.85, item.RowPriceNetWithDiscount().FloatAmount())

	var discount, rowGrossTotal, rowNetTotal, totalTaxAmount, totalDiscountAmount float64
	for _, splitItem := range splittedItems {
		assert.Equal(t, 1, splitItem.Qty)
		// make sure single and row price are equal:
		assert.Equal(t, splitItem.SinglePriceNet.FloatAmount(), splitItem.RowPriceNet.FloatAmount())
		assert.Equal(t, splitItem.SinglePriceGross.FloatAmount(), splitItem.RowPriceGross.FloatAmount())
		// make sure it's consistent (net+tax=gross):
		assert.Equal(t, splitItem.RowPriceGross.FloatAmount(), splitItem.RowPriceNet.ForceAdd(splitItem.TotalTaxAmount()).FloatAmount())
		rowGrossTotal = rowGrossTotal + splitItem.RowPriceGross.FloatAmount()
		rowNetTotal = rowNetTotal + splitItem.RowPriceNet.FloatAmount()
		totalTaxAmount = totalTaxAmount + splitItem.TotalTaxAmount().FloatAmount()
		rate, _ := splitItem.RowTaxes[0].Rate.Float64()
		assert.Equal(t, 7.0, rate)
		totalDiscountAmount = totalDiscountAmount + splitItem.TotalDiscountAmount().FloatAmount()

		// discount split cents should be at the end, so the next discount must be the same or smaller
		assert.GreaterOrEqual(t, discount, splitItem.TotalDiscountAmount().FloatAmount())
		discount = splitItem.TotalDiscountAmount().FloatAmount()
	}
	assert.Equal(t, item.RowPriceGross.FloatAmount(), rowGrossTotal)
	assert.Equal(t, item.RowPriceNet.FloatAmount(), rowNetTotal)
	assert.Equal(t, item.TotalTaxAmount().FloatAmount(), totalTaxAmount)
	assert.Equal(t, item.TotalDiscountAmount().FloatAmount(), totalDiscountAmount)

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
