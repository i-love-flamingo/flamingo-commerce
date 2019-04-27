package cart_test

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"testing"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"github.com/stretchr/testify/assert"
)

func Test_GetDeliveryCodes(t *testing.T) {
	cart := new(cartDomain.Cart)

	dummyItem := cartDomain.Item{}

	deliveryHome := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "home",
		},
		Cartitems: []cartDomain.Item{
			dummyItem,
		},
	}
	deliveryInFlight := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "inFlight",
		},
		Cartitems: []cartDomain.Item{
			dummyItem,
		},
	}
	deliveryWithoutItems := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "withoutItems",
		},
	}

	cart.Deliveries = append(cart.Deliveries, deliveryHome)
	cart.Deliveries = append(cart.Deliveries, deliveryInFlight)
	cart.Deliveries = append(cart.Deliveries, deliveryWithoutItems)

	deliveryCodes := cart.GetDeliveryCodes()

	assert.Len(t, deliveryCodes, 2)
	assert.Contains(t, deliveryCodes, "home")
	assert.Contains(t, deliveryCodes, "inFlight")
}

/*
func TestCart_HasNoMixedCart(t *testing.T) {
	var cart = new(Cart)

	cart.Cartitems = append(cart.Cartitems, getItemWithDepartureIntent())

	resultNoMixedCart := cart.HasMixedCart()
	assert.False(t, resultNoMixedCart)
}

func TestCart_HasMixedCart(t *testing.T) {
	var cart = new(Cart)

	cart.Cartitems = append(cart.Cartitems, getItemWithDepartureIntent())
	cart.Cartitems = append(cart.Cartitems, getItemWithArrivalIntent())

	resultMixedCart := cart.HasMixedCart()
	assert.True(t, resultMixedCart)
}
*/

func TestPlacedOrderInfos_GetOrderNumberForDeliveryCode(t *testing.T) {
	type args struct {
		deliveryCode string
	}
	tests := []struct {
		name string
		poi  placeorder.PlacedOrderInfos
		args args
		want string
	}{
		{
			name: "empty order infos, empty delivery code",
		},
		{
			name: "empty order infos, with code",
			args: args{
				deliveryCode: "delivery",
			},
			want: "",
		},
		{
			name: "delivery code not in placed orders",
			poi: placeorder.PlacedOrderInfos{
				placeorder.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				placeorder.PlacedOrderInfo{
					OrderNumber:  "2",
					DeliveryCode: "delivery_2",
				},
			},
			args: args{
				deliveryCode: "delivery",
			},
			want: "",
		},
		{
			name: "delivery code in placed orders",
			poi: placeorder.PlacedOrderInfos{
				placeorder.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				placeorder.PlacedOrderInfo{
					OrderNumber:  "2",
					DeliveryCode: "delivery_2",
				},
			},
			args: args{
				deliveryCode: "delivery_1",
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.poi.GetOrderNumberForDeliveryCode(tt.args.deliveryCode); got != tt.want {
				t.Errorf("PlacedOrderInfos.GetOrderNumberForDeliveryCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaxes_AddTax(t *testing.T) {
	taxes := cartDomain.Taxes{}
	taxes = taxes.AddTax(
		cartDomain.Tax{
			Amount: domain.NewFromInt(12, 1, "EUR"),
			Type:   "gst",
		})
	taxes = taxes.AddTax(
		cartDomain.Tax{
			Amount: domain.NewFromInt(1, 1, "EUR"),
			Type:   "duty",
		})
	total := taxes.TotalAmount()
	assert.Equal(t, domain.NewFromInt(13, 1, "EUR"), total)
}

func TestTaxes_AddTaxWithMerge(t *testing.T) {
	taxes := cartDomain.Taxes{}
	taxes = taxes.AddTax(
		cartDomain.Tax{
			Amount: domain.NewFromInt(12, 1, "EUR"),
			Type:   "gst",
		})
	taxes = taxes.AddTaxWithMerge(
		cartDomain.Tax{
			Amount: domain.NewFromInt(1, 1, "EUR"),
			Type:   "gst",
		})
	total := taxes.TotalAmount()
	assert.Equal(t, domain.NewFromInt(13, 1, "EUR"), total)

	assert.Equal(t, 1, len(taxes))
}

func TestCartBuilder_BuildAndGet(t *testing.T) {
	b := cartDomain.Builder{}

	cart, err := b.AddTotalitem(cartDomain.Totalitem{
		Title: "test",
		Price: domain.NewFromInt(100, 100, "EUR"),
	}).SetIds("id", "").Build()
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(100, 100, "EUR"), cart.GrandTotal(), "gradtotal need to match given total")

}
