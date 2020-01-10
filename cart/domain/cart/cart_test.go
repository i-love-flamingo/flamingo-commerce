package cart_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/testutils"

	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"

	"github.com/stretchr/testify/assert"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestCart_GetMainShippingEMail(t *testing.T) {
	t.Parallel()

	expected := "given_email"
	cart := &cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{
					DeliveryLocation: cartDomain.DeliveryLocation{
						Address: &cartDomain.Address{
							Email: expected,
						},
					},
				},
			},
		},
	}

	got := cart.GetMainShippingEMail()

	assert.Equal(t, expected, got, "email should be found")

	expected = ""
	cart = &cartDomain.Cart{}

	got = cart.GetMainShippingEMail()

	assert.Equal(t, expected, got, "email should be empty")
}

func TestCart_IsEmpty(t *testing.T) {
	t.Parallel()

	cart := &cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{},
	}
	assert.Equal(t, true, cart.IsEmpty())

	cart = &cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				Cartitems: []cartDomain.Item{
					{
						Qty: 1,
					},
				},
			},
		},
	}
	assert.Equal(t, false, cart.IsEmpty())
}

func TestCart_GetDeliveryByCode(t *testing.T) {
	t.Parallel()

	cart := cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{
					Code: "delivery_code",
				},
			},
		},
	}
	delivery, found := cart.GetDeliveryByCode("delivery_code")

	assert.True(t, found, "delivery should be found")
	assert.NotNil(t, delivery, "delivery should not be nil")

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{
					Code: "delivery_code",
				},
			},
		},
	}

	delivery, found = cart.GetDeliveryByCode("code")

	assert.False(t, found, "delivery should not be found")
	assert.Nil(t, delivery, "delivery should be nil")

	expectedDelivery := &cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "delivery_code",
		},
	}
	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			*expectedDelivery,
		},
	}

	delivery, found = cart.GetDeliveryByCode("delivery_code")

	assert.True(t, found, "delivery code should not be found")
	assert.Equal(t, expectedDelivery, delivery, "delivery should be nil")
}

func Test_GetDeliveryCodes(t *testing.T) {
	t.Parallel()

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

func TestCart_SumShippingNetWithDiscounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cart cartDomain.Cart
		want domain.Price
	}{
		{
			name: "empty cart",
			cart: cartDomain.Cart{},
			want: domain.NewZero(""),
		},
		{
			name: "cart with items with discounts but no shipping cost",
			cart: func() cartDomain.Cart {
				cart := &cartDomain.Cart{}
				cart.Deliveries = append(cart.Deliveries, *testutils.BuildDeliveryWithDifferentDiscounts(t))
				return *cart
			}(),
			want: domain.NewZero(""),
		},
		{
			name: "cart with items and shipping cost, both with discounts",
			cart: func() cartDomain.Cart {
				cart := &cartDomain.Cart{}
				cart.Deliveries = append(cart.Deliveries, *testutils.BuildDeliveryWithDifferentDiscountsAndShippingDiscounts(t))
				return *cart
			}(),
			want: domain.NewFromFloat(5.0, "$"),
		},
		{
			name: "cart with multiple deliveries with items and shipping cost, some with discounts",
			cart: func() cartDomain.Cart {
				cart := &cartDomain.Cart{}
				cart.Deliveries = append(cart.Deliveries, *testutils.BuildDeliveryWithDifferentDiscountsAndShippingDiscounts(t))
				cart.Deliveries = append(cart.Deliveries, *testutils.BuildDeliveryWithoutDiscountsAndShippingDiscounts(t))
				return *cart
			}(),
			want: domain.NewFromFloat(10.0, "$"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cart.SumShippingNetWithDiscounts(); !got.Equal(tt.want) {
				t.Errorf("Cart.SumShippingNetWithDiscount() = %v, want %v", got.Amount(), tt.want.Amount())
			}
		})
	}
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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	b := cartDomain.Builder{}

	cart, err := b.AddTotalitem(cartDomain.Totalitem{
		Title: "test",
		Price: domain.NewFromInt(100, 100, "EUR"),
	}).SetIds("id", "").Build()
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(100, 100, "EUR"), cart.GrandTotal(), "gradtotal need to match given total")

}
