package decorator_test

import (
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/testutils"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestDecoratedItem_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name string
		item *decorator.DecoratedCartItem
		want cart.AppliedDiscounts
	}{
		{
			name: "no discounts on item",
			item: &decorator.DecoratedCartItem{
				Item: cart.Item{},
			},
			want: nil,
		},
		{
			name: "multiple discounts on item",
			item: &decorator.DecoratedCartItem{
				Item: *testutils.BuildItemWithDiscounts(t),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-15.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    3,
				},
				{
					CampaignCode: "code-3",
					Label:        "title-1",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedItem.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedItem_HasDiscounts(t *testing.T) {
	tests := []struct {
		name string
		item *decorator.DecoratedCartItem
		want bool
	}{
		{
			name: "no discounts on item",
			item: &decorator.DecoratedCartItem{
				Item: cart.Item{},
			},
			want: false,
		},
		{
			name: "multiple discounts on item",
			item: func() *decorator.DecoratedCartItem {
				item := cart.Item{ID: "item-1", AppliedDiscounts: []cart.AppliedDiscount{{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
				}}}
				decorated := decorator.DecoratedCartItem{
					Item: item,
				}
				return &decorated
			}(),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("DecoratedItem.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedDelivery_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name     string
		delivery *decorator.DecoratedDelivery
		want     cart.AppliedDiscounts
	}{
		{
			name: "empty delivery",
			delivery: &decorator.DecoratedDelivery{
				Delivery: cart.Delivery{},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "delivery with items but without discounts",
			delivery: func() *decorator.DecoratedDelivery {
				delivery := cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code"}, Cartitems: []cart.Item{{}, {}}}
				decorated := decorator.DecoratedDelivery{
					Delivery: delivery,
				}
				return &decorated
			}(),
			want: cart.AppliedDiscounts{},
		},
		{
			name: "delivery with items with discounts",
			delivery: &decorator.DecoratedDelivery{
				Delivery: *testutils.BuildDeliveryWithDiscounts(t),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-30.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    3,
				},
				{
					CampaignCode: "code-3",
					Label:        "title-1",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.delivery.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedDelivery.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedDelivery_HasDiscounts(t *testing.T) {
	tests := []struct {
		name     string
		delivery *decorator.DecoratedDelivery
		want     bool
	}{
		{
			name: "empty delivery",
			delivery: &decorator.DecoratedDelivery{
				Delivery: cart.Delivery{},
			},
			want: false,
		},
		{
			name: "delivery with items but without discounts",
			delivery: &decorator.DecoratedDelivery{
				Delivery: *testutils.BuildDeliveryWithoutDiscounts(t),
			},
			want: false,
		},
		{
			name: "delivery with items with discounts",
			delivery: &decorator.DecoratedDelivery{
				Delivery: *testutils.BuildDeliveryWithDiscounts(t),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.delivery.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("DecoratedDelivery.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCart_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name string
		cart *decorator.DecoratedCart
		want cart.AppliedDiscounts
	}{
		{
			name: "empty cart",
			cart: &decorator.DecoratedCart{
				Cart: cart.Cart{},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "cart with deliveries with items but without discounts",
			cart: &decorator.DecoratedCart{
				Cart: cart.Cart{
					Deliveries: func() []cart.Delivery {
						result := make([]cart.Delivery, 0)
						delivery := cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-1"}}
						result = append(result, delivery)
						return result
					}(),
				},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "cart with deliveries with items with discounts",
			cart: &decorator.DecoratedCart{
				Cart: cart.Cart{
					Deliveries: func() []cart.Delivery {
						result := make([]cart.Delivery, 0)
						delivery := testutils.BuildDeliveryWithDiscounts(t)
						result = append(result, *delivery)
						delivery = testutils.BuildDeliveryWithDiscounts(t)
						result = append(result, *delivery)
						return result
					}(),
				},
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-60.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-40.0, "$"),
					SortOrder:    3,
				},
				{
					CampaignCode: "code-3",
					Label:        "title-1",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    4,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cart.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecoratedCart.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoratedCart_HasDiscounts(t *testing.T) {
	tests := []struct {
		name string
		cart *decorator.DecoratedCart
		want bool
	}{
		{
			name: "empty cart",
			cart: &decorator.DecoratedCart{
				Cart: cart.Cart{},
			},
			want: false,
		},
		{
			name: "cart with deliveries with items with discounts",
			cart: &decorator.DecoratedCart{
				Cart: cart.Cart{
					Deliveries: func() []cart.Delivery {
						result := make([]cart.Delivery, 0)
						delivery := testutils.BuildDeliveryWithDiscounts(t)
						result = append(result, *delivery)
						delivery = testutils.BuildDeliveryWithDiscounts(t)
						result = append(result, *delivery)
						return result
					}(),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cart.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("DecoratedCart.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
