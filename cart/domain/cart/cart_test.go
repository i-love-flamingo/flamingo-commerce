package cart_test

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
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

func TestCart_GetAllPaymentRequiredItems(t *testing.T) {
	cart := &cartDomain.Cart{
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "delivery-1"},
				Cartitems: []cartDomain.Item{
					{
						ID:                        "1",
						RowPriceGrossWithDiscount: domain.NewFromInt(1234, 100, "$"),
					},
					{
						ID:                        "2",
						RowPriceGrossWithDiscount: domain.NewFromInt(4321, 100, "$"),
					},
				},
				ShippingItem: cartDomain.ShippingItem{
					PriceGrossWithDiscounts: domain.NewZero("$"),
				},
			},
			{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "delivery-2"},
				ShippingItem: cartDomain.ShippingItem{
					PriceGrossWithDiscounts: domain.NewFromInt(55, 10, "$"),
				},
			},
		},
		Totalitems: []cartDomain.Totalitem{
			{
				Code:  "item-1",
				Price: domain.NewFromInt(6789, 100, "$"),
			},
			{
				Code:  "item-2",
				Price: domain.NewFromInt(9876, 100, "$"),
			},
		},
	}

	pricedItems := cart.GetAllPaymentRequiredItems()
	assert.Len(t, pricedItems.CartItems(), cart.ProductCount())
	assert.Equal(t, 12.34, pricedItems.CartItems()["1"].FloatAmount())
	assert.Equal(t, 43.21, pricedItems.CartItems()["2"].FloatAmount())
	assert.Len(t, pricedItems.ShippingItems(), 1)
	assert.Equal(t, 5.5, pricedItems.ShippingItems()["delivery-2"].FloatAmount())
	assert.Len(t, pricedItems.TotalItems(), len(cart.Totalitems))
	assert.Equal(t, 67.89, pricedItems.TotalItems()["item-1"].FloatAmount())
	assert.Equal(t, 98.76, pricedItems.TotalItems()["item-2"].FloatAmount())
}

func TestCart_GetContactMail(t *testing.T) {
	tests := []struct {
		name          string
		cart          cartDomain.Cart
		expectedEmail string
	}{
		{name: "no mail", cart: cartDomain.Cart{}, expectedEmail: ""},
		{name: "billing email", cart: cartDomain.Cart{BillingAddress: &cartDomain.Address{Email: "foo@example.com"}}, expectedEmail: "foo@example.com"},
		{name: "shipping email", cart: cartDomain.Cart{Deliveries: []cartDomain.Delivery{{DeliveryInfo: cartDomain.DeliveryInfo{DeliveryLocation: cartDomain.DeliveryLocation{Address: &cartDomain.Address{Email: "foo@example.com"}}}}}}, expectedEmail: "foo@example.com"},
		{name: "billing and shipping email",
			cart: cartDomain.Cart{
				BillingAddress: &cartDomain.Address{Email: "billing@example.com"},
				Deliveries:     []cartDomain.Delivery{{DeliveryInfo: cartDomain.DeliveryInfo{DeliveryLocation: cartDomain.DeliveryLocation{Address: &cartDomain.Address{Email: "shipping@example.com"}}}}},
			},
			expectedEmail: "shipping@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedEmail, tt.cart.GetContactMail())
		})
	}
}

func TestCart_GetDeliveryByCodeWithoutBool(t *testing.T) {
	cart := cartDomain.Cart{Deliveries: []cartDomain.Delivery{{DeliveryInfo: cartDomain.DeliveryInfo{Code: "valid"}}}}
	assert.Nil(t, cart.GetDeliveryByCodeWithoutBool("invalid"))
	assert.NotNil(t, cart.GetDeliveryByCodeWithoutBool("valid"))
}

func TestCart_GetByItemID(t *testing.T) {
	cart := cartDomain.Cart{Deliveries: []cartDomain.Delivery{
		{
			Cartitems: []cartDomain.Item{{ID: "item-1"}},
		},
		{
			Cartitems: []cartDomain.Item{{ID: "item-2"}},
		},
	}}

	t.Run("invalid item id should return error", func(t *testing.T) {
		got, err := cart.GetByItemID("invalid")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("valid item id should return item", func(t *testing.T) {
		got, err := cart.GetByItemID("item-2")
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "item-2", got.ID)
	})
}

func TestCart_IsPaymentSelected(t *testing.T) {
	cart := cartDomain.Cart{}
	assert.False(t, cart.IsPaymentSelected())
	cart.PaymentSelection = cartDomain.DefaultPaymentSelection{}
	assert.True(t, cart.IsPaymentSelected())
}

func TestCart_HasShippingCosts(t *testing.T) {
	cart := cartDomain.Cart{}
	assert.False(t, cart.HasShippingCosts())

	cart.ShippingNet = domain.NewFromFloat(5.00, "USD")
	assert.True(t, cart.HasShippingCosts())
}

func TestCart_AllShippingTitles(t *testing.T) {
	cart := cartDomain.Cart{}
	assert.Equal(t, []string{}, cart.AllShippingTitles())

	cart.Deliveries = []cartDomain.Delivery{
		{ShippingItem: cartDomain.ShippingItem{Title: "shipping a"}},
		{ShippingItem: cartDomain.ShippingItem{Title: "shipping b"}},
	}
	assert.ElementsMatch(t, []string{"shipping a", "shipping b"}, cart.AllShippingTitles())
}

func TestCart_SumTaxes(t *testing.T) {
	cart := cartDomain.Cart{}
	assert.Len(t, cart.SumTaxes(), 0)

	cart.Deliveries = []cartDomain.Delivery{
		{
			Cartitems: []cartDomain.Item{
				{RowTaxes: []cartDomain.Tax{{Type: "gst", Amount: domain.NewFromFloat(1.00, "USD")}}},
				{RowTaxes: []cartDomain.Tax{{Type: "foo", Amount: domain.NewFromFloat(2.00, "USD")}}},
			},
			ShippingItem: cartDomain.ShippingItem{TaxAmount: domain.NewFromFloat(2.00, "USD")},
		},
		{
			Cartitems: []cartDomain.Item{
				{RowTaxes: []cartDomain.Tax{{Type: "gst", Amount: domain.NewFromFloat(1.00, "USD")}}},
				{RowTaxes: []cartDomain.Tax{{Type: "bar", Amount: domain.NewFromFloat(2.00, "USD")}}},
			},
		},
	}

	assert.Len(t, cart.SumTaxes(), 4)
}

func TestCart_SumTotalTaxAmount(t *testing.T) {
	cart := cartDomain.Cart{}
	assert.Equal(t, cart.SumTotalTaxAmount(), domain.NewZero(""))

	cart.Deliveries = []cartDomain.Delivery{
		{
			Cartitems: []cartDomain.Item{
				{RowTaxes: []cartDomain.Tax{{Type: "gst", Amount: domain.NewFromFloat(1.00, "USD")}}},
				{RowTaxes: []cartDomain.Tax{{Type: "foo", Amount: domain.NewFromFloat(2.00, "USD")}}},
			},
			ShippingItem: cartDomain.ShippingItem{TaxAmount: domain.NewFromFloat(2.00, "USD")},
		},
		{
			Cartitems: []cartDomain.Item{
				{RowTaxes: []cartDomain.Tax{{Type: "gst", Amount: domain.NewFromFloat(1.00, "USD")}}},
				{RowTaxes: []cartDomain.Tax{{Type: "bar", Amount: domain.NewFromFloat(2.00, "USD")}}},
			},
		},
	}

	assert.Equal(t, cart.SumTotalTaxAmount(), domain.NewFromFloat(8.00, "USD"))
}

func TestCart_HasAppliedCouponCode(t *testing.T) {
	cart := cartDomain.Cart{}
	assert.False(t, cart.HasAppliedCouponCode())

	cart.AppliedCouponCodes = []cartDomain.CouponCode{{Code: "summer-2020"}}
	assert.True(t, cart.HasAppliedCouponCode())
}

func TestCart_GetPaymentReference(t *testing.T) {
	cart := cartDomain.Cart{EntityID: "e", ID: "i"}
	assert.Equal(t, "i-e", cart.GetPaymentReference())

	cart.AdditionalData.ReservedOrderID = "order-id"
	assert.Equal(t, "order-id", cart.GetPaymentReference())
}

func TestCart_GrandTotalCharges(t *testing.T) {
	cart := cartDomain.Cart{GrandTotal: domain.NewFromFloat(100.00, "USD")}
	expected := domain.NewCharges(map[string]domain.Charge{
		domain.ChargeTypeMain: {
			Type:  domain.ChargeTypeMain,
			Value: domain.NewFromFloat(100.00, "USD"),
			Price: domain.NewFromFloat(100.00, "USD"),
		},
	})

	assert.Equal(t, *expected, cart.GrandTotalCharges())
}

func TestCart_GetTotalItemsByType(t *testing.T) {
	cart := cartDomain.Cart{Totalitems: []cartDomain.Totalitem{
		{Type: "valid", Price: domain.NewFromFloat(5.00, "USD")},
		{Type: "valid", Price: domain.NewFromFloat(2.00, "USD")},
		{Type: "other", Price: domain.NewFromFloat(5.00, "USD")},
	}}
	assert.Equal(t, []cartDomain.Totalitem{}, cart.GetTotalItemsByType("invalid"))
	assert.Equal(t, []cartDomain.Totalitem{
		{Type: "valid", Price: domain.NewFromFloat(5.00, "USD")},
		{Type: "valid", Price: domain.NewFromFloat(2.00, "USD")},
	}, cart.GetTotalItemsByType("valid"))
}
