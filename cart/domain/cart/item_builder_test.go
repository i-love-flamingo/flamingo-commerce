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
