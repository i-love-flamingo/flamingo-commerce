package cart_test

import (
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
		poi  cartDomain.PlacedOrderInfos
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
			poi: cartDomain.PlacedOrderInfos{
				cartDomain.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				cartDomain.PlacedOrderInfo{
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
			poi: cartDomain.PlacedOrderInfos{
				cartDomain.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				cartDomain.PlacedOrderInfo{
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

func TestItem_PriceCalculation(t *testing.T) {

	item := cartDomain.Item{
		SinglePrice: domain.NewFromInt(1234, 100, "EUR"),
		AppliedDiscounts: []cartDomain.ItemDiscount{
			cartDomain.ItemDiscount{
				Price:         domain.NewFromInt(100, 100, "EUR"),
				IsItemRelated: true,
			},
			cartDomain.ItemDiscount{
				Price:         domain.NewFromInt(200, 100, "EUR"),
				IsItemRelated: false,
			},
		},
		TaxAmount: domain.NewFromInt(13, 100, "EUR"),
		Qty:       10,
	}

	assert.Equal(t, item.RowTotal(), domain.NewFromInt(12340, 100, "EUR"), "rowtotal is qty * singleprice")
	assert.Equal(t, item.SinglePriceInclTax(), domain.NewFromInt(1247, 100, "EUR"), "SinglePriceInclTax is SinglePrice + tax")

	assert.Equal(t, item.RowTotalInclTax(), domain.NewFromInt(12470, 100, "EUR"), "RowTotalInclTax")

	assert.Equal(t, item.ItemRelatedDiscountAmount(), domain.NewFromInt(100, 100, "EUR"), "ItemRelatedDiscountAmount")
	assert.Equal(t, item.NonItemRelatedDiscountAmount(), domain.NewFromInt(200, 100, "EUR"), "NonItemRelatedDiscountAmount")
	assert.Equal(t, item.TotalDiscountAmount(), domain.NewFromInt(300, 100, "EUR"), "TotalDiscountAmount")

	assert.Equal(t, item.RowTotalWithDiscountInclTax(), domain.NewFromInt(12170, 100, "EUR"), "RowTotalWithDiscountInclTax")
	assert.Equal(t, item.RowTotalWithItemRelatedDiscount(), domain.NewFromInt(12240, 100, "EUR"), "RowTotalWithItemRelatedDiscount")
	assert.Equal(t, item.RowTotalWithItemRelatedDiscountInclTax(), domain.NewFromInt(12370, 100, "EUR"), "RowTotalWithItemRelatedDiscountInclTax")

}
