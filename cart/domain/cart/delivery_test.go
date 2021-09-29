package cart_test

import (
	"fmt"
	"math/big"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"github.com/stretchr/testify/assert"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestDelivery_TaxCalculations(t *testing.T) {
	d := cart.Delivery{
		DeliveryInfo: cart.DeliveryInfo{Code: "foo"},
		Cartitems: []cart.Item{
			{ID: "1", RowTaxes: []cart.Tax{{Type: "gst", Rate: new(big.Float).SetInt64(7), Amount: priceDomain.NewFromInt(287, 100, "$")}}},
			{ID: "2", RowTaxes: []cart.Tax{{Type: "gst", Rate: new(big.Float).SetInt64(7), Amount: priceDomain.NewFromInt(175, 100, "$")}}},
		},
		ShippingItem: cart.ShippingItem{
			Title:     "shipping",
			TaxAmount: priceDomain.NewFromInt(55, 100, "$"),
		},
	}

	// check total tax amount
	assert.True(t, priceDomain.NewFromInt(462+55, 100, "$").Equal(d.SumTotalTaxAmount()), fmt.Sprintf("result wrong %f", d.SumTotalTaxAmount().FloatAmount()))

	// check row taxes
	taxes := d.SumRowTaxes()
	assert.Equal(t, 1, len(taxes), "expected gst to be merged")
	assert.Equal(t, "gst", taxes[0].Type, "expected gst as type")
	assert.Equal(t, new(big.Float).SetInt64(7), taxes[0].Rate, "expected rate to be correct")
	assertPricesWithLikelyEqual(t, priceDomain.NewFromInt(462, 100, "$"), taxes.TotalAmount(), "total tax amount wrong")
}

// assertPricesWithLikelyEqual - helper
func assertPricesWithLikelyEqual(t *testing.T, p1 priceDomain.Price, p2 priceDomain.Price, msg string) {
	t.Helper()
	assert.True(t, p1.LikelyEqual(p2), fmt.Sprintf("%v (%f != %f)", msg, p1.FloatAmount(), p2.FloatAmount()))

}

func TestDelivery_HasItems(t *testing.T) {
	delivery := cart.Delivery{Cartitems: []cart.Item{{}}}
	assert.True(t, delivery.HasItems())

	delivery = cart.Delivery{}
	assert.False(t, delivery.HasItems())
}

func TestShippingItem_Tax(t *testing.T) {
	shippingItem := cart.ShippingItem{TaxAmount: priceDomain.NewFromInt(250, 100, "$")}
	assert.Equal(t, "tax", shippingItem.Tax().Type)
	assert.Equal(t, priceDomain.NewFromInt(250, 100, "$"), shippingItem.Tax().Amount)
}

func TestDeliveryInfo_LoadAdditionalInfo(t *testing.T) {
	// todo:
	// var info map[string]interface{}
	// deliveryInfo := cart.DeliveryInfo{
	// 	AdditionalDeliveryInfos: map[string]json.RawMessage{"foo": []byte(`{"baz":"bar"}`)}}
	// deliveryInfo.LoadAdditionalInfo("foo", info)
}
func TestDeliveryInfo_GetAdditionalData(t *testing.T) {
	// todo:
}
func TestDeliveryInfo_AdditionalDataKeys(t *testing.T) {
	// todo:
}
func TestDeliveryInfo_GetAdditionalDeliveryInfo(t *testing.T) {
	// todo:
}
func TestDeliveryInfo_AdditionalDeliveryInfoKeys(t *testing.T) {
	// todo:
}
