package cart_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestItem_PriceCalculation(t *testing.T) {

	item := cartDomain.Item{
		SinglePriceNet:   priceDomain.NewFromInt(1234, 100, "EUR"),
		SinglePriceGross: priceDomain.NewFromInt(1247, 100, "EUR"),
		AppliedDiscounts: []cartDomain.AppliedDiscount{
			{
				Applied:       priceDomain.NewFromInt(-100, 100, "EUR"),
				IsItemRelated: true,
			},
			{
				Applied:       priceDomain.NewFromInt(-200, 100, "EUR"),
				IsItemRelated: false,
			},
		},
		RowPriceNet:   priceDomain.NewFromInt(12340, 100, "EUR"),
		RowPriceGross: priceDomain.NewFromInt(12470, 100, "EUR"),
		RowTaxes: cartDomain.Taxes([]cartDomain.Tax{
			{Amount: priceDomain.NewFromInt(130, 100, "EUR"), Type: "vat"},
		}),
		Qty: 10,
	}

	assert.Equal(t, item.ItemRelatedDiscountAmount(), priceDomain.NewFromInt(-100, 100, "EUR"), "ItemRelatedDiscountAmount")
	assert.Equal(t, item.NonItemRelatedDiscountAmount(), priceDomain.NewFromInt(-200, 100, "EUR"), "NonItemRelatedDiscountAmount")
	assert.Equal(t, item.TotalDiscountAmount(), priceDomain.NewFromInt(-300, 100, "EUR"), "TotalDiscountAmount")

	assertPricesWithLikelyEqual(t, item.RowPriceGrossWithDiscount(), priceDomain.NewFromInt(12170, 100, "EUR"), "RowPriceGrossWithDiscount")
	assertPricesWithLikelyEqual(t, item.RowPriceNetWithDiscount(), priceDomain.NewFromInt(12040, 100, "EUR"), "RowPriceNetWithDiscount")
	assertPricesWithLikelyEqual(t, item.RowPriceNetWithItemRelatedDiscount(), priceDomain.NewFromInt(12240, 100, "EUR"), "RowPriceNetWithItemRelatedDiscount")

	assert.Equal(t, 1, len(item.RowTaxes))
	assertPricesWithLikelyEqual(t, item.RowTaxes.TotalAmount(), priceDomain.NewFromInt(130, 100, "EUR"), "RowTaxes")

}

func TestItemBuild_SimpleBuild(t *testing.T) {

	f := &cartDomain.ItemBuilder{}
	item, err := f.SetSinglePriceNet(priceDomain.NewFromInt(100, 100, "EUR")).SetQty(10).SetID("22").SetExternalReference("kkk").CalculatePricesAndTaxAmountsFromSinglePriceNet().Build()
	assert.NoError(t, err)
	assert.Equal(t, "22", item.ID)
	assert.Equal(t, priceDomain.NewFromInt(1000, 100, "EUR"), item.RowPriceGross)

	// with tax from net:
	item, err = f.SetSinglePriceNet(priceDomain.NewFromInt(100, 100, "EUR")).SetQty(10).SetID("22").SetExternalReference("kkk").AddTaxInfo("default", big.NewFloat(10), nil).CalculatePricesAndTaxAmountsFromSinglePriceNet().Build()
	assert.NoError(t, err)
	assert.Equal(t, "22", item.ID)
	assert.Equal(t, priceDomain.NewFromInt(1100, 100, "EUR"), item.RowPriceGross)

	// with tax from gross:
	item, err = f.SetSinglePriceGross(priceDomain.NewFromInt(110, 100, "EUR")).SetQty(10).SetID("22").SetExternalReference("kkk").AddTaxInfo("default", big.NewFloat(10), nil).CalculatePricesAndTaxAmountsFromSinglePriceGross().Build()
	assert.NoError(t, err)
	assert.Equal(t, "22", item.ID)
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(1100, 100, "EUR"), item.RowPriceGross, "RowPriceGross wrong")
	assert.Equal(t, priceDomain.NewFromInt(100, 100, "EUR"), item.TotalTaxAmount())

}

func assertPricesWithLikelyEqual(t *testing.T, p1 priceDomain.Price, p2 priceDomain.Price, msg string) {
	assert.True(t, p1.LikelyEqual(p2), fmt.Sprintf("%v (%f != %f)", msg, p1.FloatAmount(), p2.FloatAmount()))

}

func TestItemSplitter_SplitGrossBased(t *testing.T) {
	provider := func() *cartDomain.ItemBuilder {
		b := cartDomain.ItemBuilder{}
		b.Inject(&struct {
			UseGrosPrice bool `inject:"config:commerce.product.priceIsGross,optional"`
		}{
			UseGrosPrice: true,
		})
		return &b
	}
	splitter := &cartDomain.ItemSplitter{}
	splitter.Inject(provider, &struct {
		UseGrossPrice bool `inject:"config:commerce.product.priceIsGross,optional"`
	}{
		UseGrossPrice: true,
	})

	builder := provider()
	builder.SetSinglePriceGross(priceDomain.NewFromInt(2065, 100, "€")).
		SetQty(5).AddTaxInfo("tax", big.NewFloat(7), nil).
		SetID("2").
		AddDiscount(cartDomain.AppliedDiscount{Applied: priceDomain.NewFromInt(-3170, 100, "€")}).
		CalculatePricesAndTaxAmountsFromSinglePriceGross()
	item, err := builder.Build()
	require.NoError(t, err)

	splittedItems, err := splitter.SplitInSingleQtyItems(*item)
	require.NoError(t, err)

	//20.65 * 5 = 103.25
	assert.Equal(t, 103.25, item.RowPriceGross.FloatAmount())
	assert.Equal(t, -31.70, item.TotalDiscountAmount().FloatAmount())
	// (98.57 - 31.70) * 0.07
	assert.Equal(t, 4.68, item.TotalTaxAmount().FloatAmount())
	// TotalTaxAmount + 98.57 = 103.25
	assert.Equal(t, 98.57, item.RowPriceNet.FloatAmount())
	assert.Equal(t, 66.87, item.RowPriceNetWithDiscount().FloatAmount())

	var rowGrossTotal, rowNetTotal, totalTaxAmount, totalDiscountAmount float64
	for _, splitItem := range splittedItems {
		assert.Equal(t, 1, splitItem.Qty)
		//make sure single and row price are equal:
		assert.Equal(t, splitItem.SinglePriceNet.FloatAmount(), splitItem.RowPriceNet.FloatAmount())
		assert.Equal(t, splitItem.SinglePriceGross.FloatAmount(), splitItem.RowPriceGross.FloatAmount())
		//make sure its constitent (net+tax=gross):
		assert.Equal(t, splitItem.RowPriceGross.FloatAmount(), splitItem.RowPriceNet.ForceAdd(splitItem.TotalTaxAmount()).FloatAmount())
		rowGrossTotal = rowGrossTotal + splitItem.RowPriceGross.FloatAmount()
		rowNetTotal = rowNetTotal + splitItem.RowPriceNet.FloatAmount()
		totalTaxAmount = totalTaxAmount + splitItem.TotalTaxAmount().FloatAmount()
		rate, _ := splitItem.RowTaxes[0].Rate.Float64()
		assert.Equal(t, 7.0, rate)
		totalDiscountAmount = totalDiscountAmount + splitItem.TotalDiscountAmount().FloatAmount()

	}
	assert.Equal(t, item.RowPriceGross.FloatAmount(), rowGrossTotal)
	assert.Equal(t, item.RowPriceNet.FloatAmount(), rowNetTotal)
	assert.Equal(t, item.TotalTaxAmount().FloatAmount(), totalTaxAmount)
	assert.Equal(t, item.TotalDiscountAmount().FloatAmount(), totalDiscountAmount)

}
