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
	splitter := &cartDomain.ItemSplitter{}

	items := []*cartDomain.Item{
		{
			ID:                                   "2",
			ExternalReference:                    "external",
			MarketplaceCode:                      "item related discount",
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
				cartDomain.AppliedDiscount{IsItemRelated: true, Applied: priceDomain.NewFromInt(-3172, 100, "€")},
			},
			TotalDiscountAmount:          priceDomain.NewFromInt(-3172, 100, "€"),
			ItemRelatedDiscountAmount:    priceDomain.NewFromInt(-3172, 100, "€"),
			NonItemRelatedDiscountAmount: priceDomain.NewFromInt(0, 100, "€"),
			AdditionalData:               map[string]string{"foo": "bar"},
		},
		{
			ID:                                   "2",
			ExternalReference:                    "external",
			MarketplaceCode:                      "non item related discount",
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
			ItemRelatedDiscountAmount:    priceDomain.NewFromInt(0, 100, "€"),
			NonItemRelatedDiscountAmount: priceDomain.NewFromInt(-3172, 100, "€"),
			AdditionalData:               map[string]string{"baz": "bam"},
		},
		{
			Qty:                       2,
			MarketplaceCode:           "foo",
			SinglePriceNet:            priceDomain.NewFromFloat(42.01, "EUR"),
			SinglePriceGross:          priceDomain.NewFromFloat(49.99, "EUR"),
			RowPriceNet:               priceDomain.NewFromFloat(84.02, "EUR"),
			RowPriceGross:             priceDomain.NewFromFloat(99.98, "EUR"),
			RowPriceNetWithDiscount:   priceDomain.NewFromFloat(75.61, "EUR"),
			RowPriceGrossWithDiscount: priceDomain.NewFromFloat(89.98, "EUR"),
			RowTaxes: []cartDomain.Tax{{
				Amount: priceDomain.NewFromFloat(14.37, "EUR"),
				Rate:   big.NewFloat(19.0),
				Type:   "VAT",
			}},
			TotalDiscountAmount:                  priceDomain.NewFromFloat(-10.00, "EUR"),
			NonItemRelatedDiscountAmount:         priceDomain.NewFromFloat(-10.00, "EUR"),
			ItemRelatedDiscountAmount:            priceDomain.NewZero("EUR"),
			RowPriceNetWithItemRelatedDiscount:   priceDomain.NewFromFloat(84.02, "EUR"),
			RowPriceGrossWithItemRelatedDiscount: priceDomain.NewFromFloat(99.98, "EUR"),
			AppliedDiscounts:                     []cartDomain.AppliedDiscount{{Applied: priceDomain.NewFromFloat(-10.00, "EUR")}},
		},
		{
			Qty:                                  1,
			MarketplaceCode:                      "no need for split, no taxes, no discounts",
			SinglePriceNet:                       priceDomain.NewFromFloat(42.01, "EUR"),
			SinglePriceGross:                     priceDomain.NewFromFloat(42.01, "EUR"),
			RowPriceNet:                          priceDomain.NewFromFloat(42.01, "EUR"),
			RowPriceGross:                        priceDomain.NewFromFloat(42.01, "EUR"),
			RowPriceNetWithDiscount:              priceDomain.NewFromFloat(42.01, "EUR"),
			RowPriceGrossWithDiscount:            priceDomain.NewFromFloat(42.01, "EUR"),
			TotalDiscountAmount:                  priceDomain.NewZero("EUR"),
			NonItemRelatedDiscountAmount:         priceDomain.NewZero("EUR"),
			ItemRelatedDiscountAmount:            priceDomain.NewZero("EUR"),
			RowPriceNetWithItemRelatedDiscount:   priceDomain.NewZero("EUR"),
			RowPriceGrossWithItemRelatedDiscount: priceDomain.NewZero("EUR"),
		},
	}

	for _, item := range items {
		t.Run(item.MarketplaceCode, func(t *testing.T) {
			splitItems, err := splitter.SplitInSingleQtyItems(*item)
			require.NoError(t, err)
			assert.Len(t, splitItems, item.Qty)

			var (
				discount,
				rowGrossTotal,
				rowGrossWithDiscount,
				rowGrossWithItemDiscount,
				rowNetWithDiscounts,
				rowNetWithItemDiscount,
				itemDiscountAmount,
				nonItemDiscountAmount,
				rowNetTotal,
				totalTaxAmount,
				totalDiscountAmount float64
			)
			appliedDiscounts := make([]float64, len(item.AppliedDiscounts))
			for _, splitItem := range splitItems {
				assert.Equal(t, item.ID, splitItem.ID, "ID")
				assert.Equal(t, item.ExternalReference, splitItem.ExternalReference, "ExternalReference")
				assert.Equal(t, item.MarketplaceCode, splitItem.MarketplaceCode, "MarketplaceCode")
				assert.Equal(t, item.VariantMarketPlaceCode, splitItem.VariantMarketPlaceCode, "VariantMarketPlaceCode")
				assert.Equal(t, item.ProductName, splitItem.ProductName, "ProductName")
				assert.Equal(t, item.SourceID, splitItem.SourceID, "SourceID")
				assert.Equal(t, item.AdditionalData, splitItem.AdditionalData)
				assert.Equal(t, 1, splitItem.Qty)
				// make sure single and row price are equal:
				assert.Equal(t, splitItem.SinglePriceNet.FloatAmount(), splitItem.RowPriceNet.FloatAmount())
				assert.Equal(t, splitItem.SinglePriceGross.FloatAmount(), splitItem.RowPriceGross.FloatAmount())
				// make sure it's consistent (net+tax=gross):
				assert.Equal(t, splitItem.RowPriceGrossWithDiscount.FloatAmount(), splitItem.RowPriceNetWithDiscount.ForceAdd(splitItem.TotalTaxAmount()).FloatAmount())

				rowGrossTotal += splitItem.RowPriceGross.FloatAmount()
				rowNetTotal += splitItem.RowPriceNet.FloatAmount()
				totalTaxAmount += splitItem.TotalTaxAmount().FloatAmount()
				rowGrossWithDiscount += splitItem.RowPriceGrossWithDiscount.FloatAmount()
				rowGrossWithItemDiscount += splitItem.RowPriceGrossWithItemRelatedDiscount.FloatAmount()
				rowNetWithDiscounts += splitItem.RowPriceNetWithDiscount.FloatAmount()
				rowNetWithItemDiscount += splitItem.RowPriceNetWithItemRelatedDiscount.FloatAmount()
				itemDiscountAmount += splitItem.ItemRelatedDiscountAmount.FloatAmount()
				nonItemDiscountAmount += splitItem.NonItemRelatedDiscountAmount.FloatAmount()
				if len(splitItem.RowTaxes) > 0 {
					rate, _ := splitItem.RowTaxes[0].Rate.Float64()
					expectedRate, _ := item.RowTaxes[0].Rate.Float64()
					assert.Equal(t, expectedRate, rate)
				}
				totalDiscountAmount = totalDiscountAmount + splitItem.TotalDiscountAmount.FloatAmount()

				require.Len(t, splitItem.AppliedDiscounts, len(item.AppliedDiscounts))
				for i, appliedDiscount := range splitItem.AppliedDiscounts {
					appliedDiscounts[i] += appliedDiscount.Applied.FloatAmount()
					assert.Equal(t, item.AppliedDiscounts[i].CampaignCode, appliedDiscount.CampaignCode)
					assert.Equal(t, item.AppliedDiscounts[i].IsItemRelated, appliedDiscount.IsItemRelated)
					assert.Equal(t, item.AppliedDiscounts[i].CouponCode, appliedDiscount.CouponCode)
					assert.Equal(t, item.AppliedDiscounts[i].Type, appliedDiscount.Type)
					assert.Equal(t, item.AppliedDiscounts[i].SortOrder, appliedDiscount.SortOrder)
					assert.Equal(t, item.AppliedDiscounts[i].Label, appliedDiscount.Label)
				}

				// discount split cents should be at the end, so the next discount must be the same or smaller
				assert.GreaterOrEqual(t, discount, splitItem.TotalDiscountAmount.FloatAmount())
				discount = splitItem.TotalDiscountAmount.FloatAmount()
			}

			assert.Equal(t, item.RowPriceGrossWithDiscount.FloatAmount(), rowGrossWithDiscount)
			assert.Equal(t, item.RowPriceGrossWithItemRelatedDiscount.FloatAmount(), rowGrossWithItemDiscount)
			assert.Equal(t, item.RowPriceNetWithDiscount.FloatAmount(), rowNetWithDiscounts)
			assert.Equal(t, item.RowPriceNetWithItemRelatedDiscount.FloatAmount(), rowNetWithItemDiscount)
			assert.Equal(t, item.ItemRelatedDiscountAmount.FloatAmount(), itemDiscountAmount)
			assert.Equal(t, item.NonItemRelatedDiscountAmount.FloatAmount(), nonItemDiscountAmount)
			assert.Equal(t, item.TotalDiscountAmount.FloatAmount(), totalDiscountAmount)
			for i, appliedDiscount := range item.AppliedDiscounts {
				assert.Equal(t, appliedDiscount.Applied.FloatAmount(), appliedDiscounts[i])
			}

			assert.Equal(t, item.RowPriceGross.FloatAmount(), rowGrossTotal)
			assert.Equal(t, item.RowPriceNet.FloatAmount(), rowNetTotal)
			assert.True(t, item.TotalTaxAmount().LikelyEqual(priceDomain.NewFromFloat(totalTaxAmount, item.TotalTaxAmount().Currency())))
		})
	}
}

func TestItem_AdditionalDataKeys(t *testing.T) {
	item := cartDomain.Item{
		ID:             "2",
		AdditionalData: map[string]string{"foo": "bar", "baz": "bam"},
	}

	assert.ElementsMatch(t, []string{"foo", "baz"}, item.AdditionalDataKeys())
}

func TestItem_AdditionalDataValues(t *testing.T) {
	item := cartDomain.Item{
		ID:             "2",
		AdditionalData: map[string]string{"foo": "bar", "baz": "bam"},
	}

	assert.ElementsMatch(t, []string{"bar", "bam"}, item.AdditionalDataValues())
}

func TestItem_HasAdditionalDataKey(t *testing.T) {
	item := cartDomain.Item{
		ID:             "2",
		AdditionalData: map[string]string{"foo": "bar"},
	}
	assert.True(t, item.HasAdditionalDataKey("foo"))
	assert.False(t, item.HasAdditionalDataKey("bam"))
}
func TestItem_GetAdditionalData(t *testing.T) {
	item := cartDomain.Item{
		ID:             "2",
		AdditionalData: map[string]string{"foo": "bar"},
	}

	assert.Equal(t, "bar", item.GetAdditionalData("foo"))
	assert.Equal(t, "", item.GetAdditionalData("bar"))
}
