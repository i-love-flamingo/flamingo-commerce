package cart

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestDeliveryInfo_TotalCalculations(t *testing.T) {
	df := DeliveryBuilder{}
	df.SetDeliveryCode("test")

	// Add item with 4,1 $ - 10 Qty and Tax of 7
	itemf := ItemBuilder{}
	// Set net price 410
	itemf.SetSinglePriceNet(
		priceDomain.NewFromInt(410, 100, "$"),
	).SetQty(10).AddTaxInfo(
		"gst", new(big.Float).SetInt64(7), nil,
	).SetID("1")

	item, err := itemf.CalculatePricesAndTaxAmountsFromSinglePriceNet().Build()
	if err != nil {
		assert.FailNow(t, "no error excpected here", err)
	}
	if item == nil {
		t.Fatal("item is nil but no error?")
	}

	assert.Equal(t, priceDomain.NewFromInt(4387, 100, "$"), item.RowPriceGross, "item1 gross price wrong")
	df.AddItem(*item)

	// Add item with 10 $ - 5 Qty / Discount of 25 (505) and Tax of 7
	// Set net price 410
	itemf.SetSinglePriceNet(
		priceDomain.NewFromInt(1000, 100, "$"),
	).SetQty(5).AddTaxInfo(
		"gst", new(big.Float).SetInt64(7), nil,
	).SetID("1")
	itemf.AddDiscount(AppliedDiscount{
		CampaignCode:  "summercampaign",
		Applied:       priceDomain.NewFromInt(-2500, 100, "$"),
		IsItemRelated: true,
	})
	item2, err := itemf.CalculatePricesAndTaxAmountsFromSinglePriceNet().Build()
	if err != nil {
		assert.FailNow(t, "no error excpected here", err)
	}
	if item2 == nil {
		t.Fatal("item2 is nil but no error?")
	}
	// item2 gros should be tax = (5 * 10) - discount (25) * 0.07 = 1,75
	assert.Equal(
		t,
		priceDomain.NewFromInt(5175, 100, "$").FloatAmount(),
		item2.RowPriceGross.FloatAmount(),
		"item2 gross price wrong",
	)
	df.AddItem(*item2)

	df.SetShippingItem(ShippingItem{
		Title:      "shipping",
		PriceNet:   priceDomain.NewFromInt(500, 100, "$"),
		TaxAmount:  priceDomain.NewFromInt(55, 100, "$"),
		PriceGross: priceDomain.NewFromInt(555, 100, "$"),
		AppliedDiscounts: AppliedDiscounts{
			{
				Applied: priceDomain.NewFromInt(-20, 100, "$"),
			},
		},
	})

	// Check totals
	d, err := df.Build()
	assert.NoError(t, err)

	// SubTotalGross - need to be 5175 + 4387
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(9562, 100, "$"), d.SubTotalGross(), "SubTotalGross result should match 95,62")

	// assert.Equal(t, priceDomain.NewFromInt(9562, 100, "$"), d.SubTotalGross(), fmt.Sprintf("SubTotalGross result should match 95,62 but is %f", d.SubTotalGross().FloatAmount()))

	// SubTotalGross - need to be 4100 + 5000
	assert.True(t, priceDomain.NewFromInt(9100, 100, "$").Equal(d.SubTotalNet()), "SubTotalNet wrong")

	// SumTotalTaxAmount is the difference
	assert.True(t, priceDomain.NewFromInt(462+55, 100, "$").Equal(d.SumTotalTaxAmount()), fmt.Sprintf("result wrong %f", d.SumTotalTaxAmount().FloatAmount()))

	// SumTotalDiscountAmount
	assert.True(t, priceDomain.NewFromInt(-2500-20, 100, "$").Equal(d.SumTotalDiscountAmount()), fmt.Sprintf("SumTotalDiscountAmount result wrong %f", d.SumTotalDiscountAmount().FloatAmount()))

	// SubTotalNetWithDiscounts
	assert.True(t, priceDomain.NewFromInt(9100-2500, 100, "$").Equal(d.SubTotalNetWithDiscounts()), fmt.Sprintf("SubTotalNetWithDiscounts result wrong %f", d.SubTotalNetWithDiscounts().FloatAmount()))

	// SubTotalGrossWithDiscounts
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(9562-2500, 100, "$"), d.SubTotalGrossWithDiscounts(), "SubTotalGrossWithDiscount")

	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(-2500, 100, "$"), d.SumItemRelatedDiscountAmount(), "SumItemRelatedDiscountAmount")
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(-20, 100, "$"), d.SumNonItemRelatedDiscountAmount(), "SumNonItemRelatedDiscountAmount")
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(-2500-20, 100, "$"), d.SumTotalDiscountAmount(), "SumTotalDiscountAmount")

	// Taxes check
	taxes := d.SumRowTaxes()
	assert.Equal(t, 1, len(taxes), "expected one merged tax")
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(462, 100, "$"), taxes.TotalAmount(), "taxes check wrong")

	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(9562-2500+555-20, 100, "$"), d.GrandTotal(), "GrandTotal")
}

// assertPricesWithLikelyEqual - helper
func assertPricesWithLikelyEqual(t *testing.T, p1 priceDomain.Price, p2 priceDomain.Price, msg string) {
	t.Helper()
	assert.True(t, p1.LikelyEqual(p2), fmt.Sprintf("%v (%f != %f)", msg, p1.FloatAmount(), p2.FloatAmount()))

}
