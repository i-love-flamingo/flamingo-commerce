package cart_test

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/testutils"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func assertDeepClone(t testing.TB, orig, cloned reflect.Value) {
	t.Helper()
	require.Equal(t, orig.Type(), cloned.Type())
	switch orig.Kind() {
	case reflect.Struct:
		for i := 0; i < orig.NumField(); i++ {
			assertDeepClone(t, orig.Field(i), cloned.Field(i))
		}
	case reflect.Slice, reflect.Array:
		if orig.IsNil() {
			return
		}
		assert.NotEqual(t, orig.Pointer(), cloned.Pointer())
		for i := 0; i < orig.Len(); i++ {
			assertDeepClone(t, orig.Index(i), cloned.Index(i))
		}
	case reflect.Map:
		if orig.IsNil() {
			return
		}
		assert.NotEqual(t, orig.Pointer(), cloned.Pointer())
		iter := orig.MapRange()
		for iter.Next() {
			assertDeepClone(t, iter.Value(), cloned.MapIndex(iter.Key()))
		}
	case reflect.Ptr:
		if orig.IsNil() {
			return
		}
		assert.NotEqual(t, orig.Pointer(), cloned.Pointer())
		assertDeepClone(t, orig.Elem(), cloned.Elem())
	default:
	}

}

func TestCart_Clone(t *testing.T) {
	t.Parallel()
	cart := cartDomain.Cart{
		ID:             "original",
		BillingAddress: &cartDomain.Address{},
		Purchaser: &cartDomain.Person{
			Address:              &cartDomain.Address{},
			PersonalDetails:      cartDomain.PersonalDetails{},
			ExistingCustomerData: &cartDomain.ExistingCustomerData{ID: "1"},
		},
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{
					Code: "code",
					DeliveryLocation: cartDomain.DeliveryLocation{
						Address:           &cartDomain.Address{},
						UseBillingAddress: false,
						Code:              "",
					},
					AdditionalData:          map[string]string{"hello": "you"},
					AdditionalDeliveryInfos: map[string]json.RawMessage{"foo": json.RawMessage("test")},
				},
				Cartitems: []cartDomain.Item{
					{
						AdditionalData: map[string]string{"hello": "you"},
					},
				},
				ShippingItem: cartDomain.ShippingItem{
					Title: "",
				},
			},
		},
		AdditionalData: cartDomain.AdditionalData{CustomAttributes: map[string]string{"hello": "you"}},
	}

	cloned, err := cart.Clone()
	require.NoError(t, err)

	assert.True(t, reflect.DeepEqual(cart, cloned), "cloned cart should have same values")

	assertDeepClone(t, reflect.ValueOf(cart), reflect.ValueOf(cloned))

	// some alibi changes to check that it is really a clone
	cloned.AdditionalData.CustomAttributes["hello"] = "bar"
	assert.Equal(t, "you", cart.AdditionalData.CustomAttributes["hello"])
}

func TestCart_GetMainShippingEMail(t *testing.T) {
	t.Parallel()

	email := "given_email"
	cart := &cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{
					DeliveryLocation: cartDomain.DeliveryLocation{
						Address: &cartDomain.Address{
							Email: email,
						},
					},
				},
			},
		},
	}

	got := cart.GetMainShippingEMail()

	assert.Equal(t, email, got, "email should be found")

	email = ""
	cart = &cartDomain.Cart{}

	got = cart.GetMainShippingEMail()

	assert.Equal(t, email, got, "email should be empty")
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

	code := "delivery_code"
	delivery := &cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: code,
		},
	}
	cart := cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			*delivery,
		},
	}

	got, found := cart.GetDeliveryByCode(code)

	assert.True(t, found, "delivery should be found")
	assert.Equal(t, delivery, got, "delivery should not be nil")

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			*delivery,
		},
	}

	got, found = cart.GetDeliveryByCode("not_existing")

	assert.False(t, found, "delivery should not be found")
	assert.Equal(t, (*cartDomain.Delivery)(nil), got, "delivery should be nil")
}

func TestHasDeliveryForCode(t *testing.T) {
	t.Parallel()

	code := "delivery_code"
	delivery := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: code,
		},
	}

	cart := cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			delivery,
		},
	}

	found := cart.HasDeliveryForCode(code)
	assert.True(t, found, "delivery should be found")

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			delivery,
		},
	}

	found = cart.HasDeliveryForCode("not_existing")
	assert.False(t, found, "delivery should not be found")
}

func TestGetDeliveryCodes(t *testing.T) {
	t.Parallel()

	cart := &cartDomain.Cart{}

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

func TestCart_GetDeliveryByItemID(t *testing.T) {
	t.Parallel()

	id := "item_id"
	delivery := &cartDomain.Delivery{
		Cartitems: []cartDomain.Item{
			{
				ID: id,
			},
		},
	}
	cart := cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			*delivery,
		},
	}

	got, err := cart.GetDeliveryByItemID(id)

	assert.Equal(t, delivery, got)
	assert.NoError(t, err)

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{},
	}

	got, err = cart.GetDeliveryByItemID(id)

	assert.Equal(t, (*cartDomain.Delivery)(nil), got)
	assert.Error(t, err)
}

func TestCart_GetTotalQty(t *testing.T) {
	t.Parallel()

	marketplaceCode := "marketplacecode"
	variantCode := "variantcode"
	expected := 1

	cart := cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				Cartitems: []cartDomain.Item{
					{
						MarketplaceCode:        marketplaceCode,
						VariantMarketPlaceCode: variantCode,
						Qty:                    1,
					},
				},
			},
		},
	}

	got := cart.GetTotalQty(marketplaceCode, variantCode)

	assert.Equal(t, expected, got)

	expected = 2

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				Cartitems: []cartDomain.Item{
					{
						MarketplaceCode:        marketplaceCode,
						VariantMarketPlaceCode: variantCode,
						Qty:                    1,
					},
					{
						MarketplaceCode:        marketplaceCode,
						VariantMarketPlaceCode: variantCode,
						Qty:                    1,
					},
				},
			},
		},
	}

	got = cart.GetTotalQty(marketplaceCode, variantCode)

	assert.Equal(t, expected, got)

	expected = 0

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				Cartitems: []cartDomain.Item{
					{
						MarketplaceCode:        marketplaceCode,
						VariantMarketPlaceCode: variantCode,
					},
				},
			},
		},
	}

	got = cart.GetTotalQty(marketplaceCode, variantCode)

	assert.Equal(t, expected, got)
}

func TestCart_GetByExternalReference(t *testing.T) {
	t.Parallel()

	externalReference := "reference"
	expected := &cartDomain.Item{
		ExternalReference: externalReference,
	}

	cart := cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				Cartitems: []cartDomain.Item{
					*expected,
				},
			},
		},
	}

	got, err := cart.GetByExternalReference(externalReference)

	assert.Equal(t, expected, got)
	assert.NoError(t, err)

	expected = nil

	cart = cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				Cartitems: []cartDomain.Item{},
			},
		},
	}

	got, err = cart.GetByExternalReference(externalReference)

	assert.Equal(t, expected, got)
	assert.Error(t, err)
}

func TestCart_GetVoucherSavings(t *testing.T) {
	t.Parallel()

	price := domain.NewFromBigFloat(*big.NewFloat(1.0), "")
	cart := cartDomain.Cart{
		Totalitems: []cartDomain.Totalitem{
			{
				Price: price,
				Type:  cartDomain.TotalsTypeVoucher,
			},
		},
	}

	got := cart.GetVoucherSavings()
	assert.Equal(t, got, price)

	price = domain.NewFromBigFloat(*big.NewFloat(2.0), "")
	cart = cartDomain.Cart{
		Totalitems: []cartDomain.Totalitem{
			{
				Price: price,
				Type:  cartDomain.TotalsTypeVoucher,
			},
		},
	}

	got = cart.GetVoucherSavings()
	assert.Equal(t, got, price)

	cart = cartDomain.Cart{
		Totalitems: []cartDomain.Totalitem{
			{
				Price: domain.NewFromBigFloat(*big.NewFloat(2.0), "c1"),
				Type:  cartDomain.TotalsTypeVoucher,
			},
			{
				Price: domain.NewFromBigFloat(*big.NewFloat(2.0), "c2"),
				Type:  cartDomain.TotalsTypeVoucher,
			},
		},
	}

	got = cart.GetVoucherSavings()
	assert.Equal(t, got, domain.NewZero("c1"))

	cart = cartDomain.Cart{
		Totalitems: []cartDomain.Totalitem{
			{
				Price: domain.NewFromBigFloat(*big.NewFloat(-2.0), ""),
				Type:  cartDomain.TotalsTypeVoucher,
			},
		},
	}

	got = cart.GetVoucherSavings()
	assert.Equal(t, got, domain.Price{})
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

func TestCart_SumShippingGrossWithDiscounts(t *testing.T) {
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
			want: domain.NewFromFloat(7.0, "$"),
		},
		{
			name: "cart with multiple deliveries with items and shipping cost, some with discounts",
			cart: func() cartDomain.Cart {
				cart := &cartDomain.Cart{}
				cart.Deliveries = append(cart.Deliveries, *testutils.BuildDeliveryWithDifferentDiscountsAndShippingDiscounts(t))
				cart.Deliveries = append(cart.Deliveries, *testutils.BuildDeliveryWithoutDiscountsAndShippingDiscounts(t))
				return *cart
			}(),
			want: domain.NewFromFloat(14.0, "$"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cart.SumShippingGrossWithDiscounts(); !got.Equal(tt.want) {
				t.Errorf("Cart.SumShippingGrossWithDiscounts() = %v, want %v", got.Amount(), tt.want.Amount())
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

func TestAppliedCouponCodes_ContainedIn(t *testing.T) {
	type args struct {
		couponCodesToCompare cartDomain.AppliedCouponCodes
	}
	tests := []struct {
		name string
		acc  cartDomain.AppliedCouponCodes
		args args
		want bool
	}{
		{
			name: "empty coupon codes are contained",
			acc:  cartDomain.AppliedCouponCodes{},
			args: args{
				couponCodesToCompare: cartDomain.AppliedCouponCodes{},
			},
			want: true,
		},
		{
			name: "same coupon codes are contained",
			acc: cartDomain.AppliedCouponCodes{
				cartDomain.CouponCode{
					Code:             "some-code",
					CustomAttributes: nil,
				},
				cartDomain.CouponCode{
					Code:             "some-other-code",
					CustomAttributes: nil,
				},
			},
			args: args{
				couponCodesToCompare: cartDomain.AppliedCouponCodes{
					cartDomain.CouponCode{
						Code:             "some-code",
						CustomAttributes: nil,
					},
					cartDomain.CouponCode{
						Code:             "some-other-code",
						CustomAttributes: nil,
					},
				},
			},
			want: true,
		},
		{
			name: "same but inverted coupon codes are contained",
			acc: cartDomain.AppliedCouponCodes{
				cartDomain.CouponCode{
					Code:             "some-code",
					CustomAttributes: nil,
				},
				cartDomain.CouponCode{
					Code:             "some-other-code",
					CustomAttributes: nil,
				},
			},
			args: args{
				couponCodesToCompare: cartDomain.AppliedCouponCodes{
					cartDomain.CouponCode{
						Code:             "some-other-code",
						CustomAttributes: nil,
					},
					cartDomain.CouponCode{
						Code:             "some-code",
						CustomAttributes: nil,
					},
				},
			},
			want: true,
		},
		{
			name: "same but different amount of coupon codes are contained",
			acc: cartDomain.AppliedCouponCodes{
				cartDomain.CouponCode{
					Code:             "some-code",
					CustomAttributes: nil,
				},
			},
			args: args{
				couponCodesToCompare: cartDomain.AppliedCouponCodes{
					cartDomain.CouponCode{
						Code:             "some-other-code",
						CustomAttributes: nil,
					},
					cartDomain.CouponCode{
						Code:             "some-code",
						CustomAttributes: nil,
					},
				},
			},
			want: true,
		},
		{
			name: "different coupon codes are not contained",
			acc: cartDomain.AppliedCouponCodes{
				cartDomain.CouponCode{
					Code:             "some-code",
					CustomAttributes: nil,
				},
				cartDomain.CouponCode{
					Code:             "some-other-code",
					CustomAttributes: nil,
				},
			},
			args: args{
				couponCodesToCompare: cartDomain.AppliedCouponCodes{
					cartDomain.CouponCode{
						Code:             "some-code",
						CustomAttributes: nil,
					},
					cartDomain.CouponCode{
						Code:             "some-different-code",
						CustomAttributes: nil,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.acc.ContainedIn(tt.args.couponCodesToCompare); got != tt.want {
				t.Errorf("ContainedIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_ProductCountAndUniqueProductCount(t *testing.T) {
	tests := []struct {
		name                   string
		cart                   cartDomain.Cart
		wantProductCount       int
		wantProductUniqueCount int
	}{
		{
			name:                   "empty cart",
			cart:                   cartDomain.Cart{},
			wantProductCount:       0,
			wantProductUniqueCount: 0,
		},
		{
			name: "single delivery in cart",
			cart: cartDomain.Cart{
				Deliveries: []cartDomain.Delivery{
					{
						Cartitems: []cartDomain.Item{
							{
								MarketplaceCode: "product1",
							},
						},
					},
				},
			},
			wantProductCount:       1,
			wantProductUniqueCount: 1,
		},
		{
			name: "two deliveries in cart with different products",
			cart: cartDomain.Cart{
				Deliveries: []cartDomain.Delivery{
					{
						Cartitems: []cartDomain.Item{
							{
								MarketplaceCode: "product1",
							},
						},
					},
					{
						Cartitems: []cartDomain.Item{
							{
								MarketplaceCode: "product2",
							},
						},
					},
				},
			},
			wantProductCount:       2,
			wantProductUniqueCount: 2,
		},
		{
			name: "two deliveries in cart with a product in both deliveries",
			cart: cartDomain.Cart{
				Deliveries: []cartDomain.Delivery{
					{
						Cartitems: []cartDomain.Item{
							{
								MarketplaceCode: "product1",
							},
						},
					},
					{
						Cartitems: []cartDomain.Item{
							{
								MarketplaceCode: "product1",
							},
							{
								MarketplaceCode: "product2",
							},
						},
					},
				},
			},
			wantProductCount:       3,
			wantProductUniqueCount: 2,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantProductCount, tt.cart.ProductCount(), "ProductCount has wrong result, expected %#v but got %#v", tt.wantProductCount, tt.cart.ProductCount())
		assert.Equal(t, tt.wantProductUniqueCount, tt.cart.ProductCountUnique(), "ProductCountUnique has wrong result, expected %#v but got %#v", tt.wantProductUniqueCount, tt.cart.ProductCountUnique())
	}
}
