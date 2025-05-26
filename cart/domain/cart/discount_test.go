package cart_test

import (
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/testutils"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestCart_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name string
		cart *cart.Cart
		want cart.AppliedDiscounts
	}{
		{
			name: "empty cart",
			cart: &cart.Cart{},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "cart with deliveries but without items",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-1"}}
					result = append(result, delivery)
					delivery = cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-2"}}
					result = append(result, delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "cart with deliveries with items but without discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := cart.Delivery{
						DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
						Cartitems:    []cart.Item{{}, {}},
					}
					result = append(result, delivery)
					delivery = cart.Delivery{
						DeliveryInfo: cart.DeliveryInfo{Code: "code-2"},
					}
					result = append(result, delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "cart with deliveries with items with discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithDiscounts
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithDiscounts
					result = append(result, *delivery)
					return result
				}(),
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
				{
					CampaignCode: "code-7",
					Label:        "title-7",
					Type:         "type-7",
					Applied:      domain.NewFromFloat(-40.0, "$"),
					SortOrder:    6,
					CustomAttributes: map[string]interface{}{
						"attr1": 3,
						"attr2": 1,
					},
				},
			},
		},
		{
			name: "cart with deliveries with items with duplicate discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithDuplicateDiscounts
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithDuplicateDiscounts
					result = append(result, *delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-40.0, "$"),
					SortOrder:    0,
				},
			},
		},
		{
			name: "cart with different deliveries with items with discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithDiscounts
					result = append(result, *delivery)
					delivery = testutils.BuildAlternativeDeliveryWithAlternativeDiscounts
					result = append(result, *delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-5",
					Label:        "title-5",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-30.0, "$"),
					SortOrder:    0,
				},
				{
					CampaignCode: "code-6",
					Label:        "title-6",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    1,
				},
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
				{
					CampaignCode: "code-4",
					Label:        "title-4",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    5,
				},
				{
					CampaignCode: "code-7",
					Label:        "title-7",
					Type:         "type-7",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    6,
					CustomAttributes: map[string]interface{}{
						"attr1": 3,
						"attr2": 1,
					},
				},
			},
		},
		{
			name: "cart with deliveries with items and shipping discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)

					delivery := cart.Delivery{
						DeliveryInfo: cart.DeliveryInfo{Code: "code-1"},
						Cartitems:    []cart.Item{{}, {}},
						ShippingItem: *testutils.BuildShippingItemWithDiscounts,
					}
					result = append(result, delivery)
					delivery = cart.Delivery{
						DeliveryInfo: cart.DeliveryInfo{Code: "code-2"},
						Cartitems:    []cart.Item{{}, {}},
						ShippingItem: *testutils.BuildShippingItemWithDiscounts,
					}
					result = append(result, delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    3,
				},
			},
		},
		{
			name: "cart with deliveries with items with discounts and shipping discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithoutDiscountsAndShippingDiscounts
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithoutDiscountsAndShippingDiscounts
					result = append(result, *delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.cart.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_HasDiscounts(t *testing.T) {
	tests := []struct {
		name string
		cart *cart.Cart
		want bool
	}{
		{
			name: "empty cart",
			cart: &cart.Cart{},
			want: false,
		},
		{
			name: "cart with deliveries but without items",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-1"}}
					result = append(result, delivery)
					delivery = cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-2"}}
					result = append(result, delivery)
					return result
				}(),
			},
			want: false,
		},
		{
			name: "cart with deliveries with items but without discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-1"}, Cartitems: []cart.Item{{}, {}}}
					result = append(result, delivery)
					delivery = cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code-2"}, Cartitems: []cart.Item{{}, {}}}
					result = append(result, delivery)
					return result
				}(),
			},
			want: false,
		},
		{
			name: "cart with deliveries with items with discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithDiscounts
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithDiscounts
					result = append(result, *delivery)
					return result
				}(),
			},
			want: true,
		},
		{
			name: "cart with deliveries with items with duplicate discounts and shipping discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts
					result = append(result, *delivery)
					return result
				}(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.cart.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("Cart.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelivery_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name     string
		delivery *cart.Delivery
		want     cart.AppliedDiscounts
	}{
		{
			name:     "empty delivery",
			delivery: &cart.Delivery{},
			want:     cart.AppliedDiscounts{},
		},
		{
			name:     "delivery with items but without discounts",
			delivery: testutils.BuildDeliveryWithoutDiscounts,
			want:     cart.AppliedDiscounts{},
		},
		{
			name:     "delivery with items with different discounts",
			delivery: testutils.BuildDeliveryWithDifferentDiscounts,
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-5",
					Label:        "title-5",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-15.0, "$"),
					SortOrder:    0,
				},
				{
					CampaignCode: "code-6",
					Label:        "title-6",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    1,
				},
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
				{
					CampaignCode: "code-4",
					Label:        "title-4",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    5,
				},
				{
					CampaignCode: "code-7",
					Label:        "title-7",
					Type:         "type-7",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    6,
					CustomAttributes: map[string]interface{}{
						"attr1": 3,
						"attr2": 1,
					},
				},
			},
		},
		{
			name:     "delivery with item with duplicate discounts",
			delivery: testutils.BuildDeliveryWithDuplicateDiscounts,
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-20.0, "$"),
					SortOrder:    0,
				},
			},
		},
		{
			name:     "delivery with items but without discounts and shipping discounts",
			delivery: testutils.BuildDeliveryWithoutDiscountsAndShippingDiscounts,
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    3,
				},
			},
		},
		{
			name:     "delivery with items with different discounts and shipping discounts",
			delivery: testutils.BuildDeliveryWithDifferentDiscountsAndShippingDiscounts,
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-5",
					Label:        "title-5",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-15.0, "$"),
					SortOrder:    0,
				},
				{
					CampaignCode: "code-6",
					Label:        "title-6",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    1,
				},
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-20.0, "$"),
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
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    4,
				},
				{
					CampaignCode: "code-4",
					Label:        "title-4",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    5,
				},
				{
					CampaignCode: "code-7",
					Label:        "title-7",
					Type:         "type-7",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    6,
					CustomAttributes: map[string]interface{}{
						"attr1": 3,
						"attr2": 1,
					},
				},
			},
		},
		{
			name:     "delivery with item with duplicate discounts and shipping discounts",
			delivery: testutils.BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts,
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-30.0, "$"),
					SortOrder:    0,
				},
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.delivery.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delivery.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelivery_HasDiscounts(t *testing.T) {
	tests := []struct {
		name     string
		delivery *cart.Delivery
		want     bool
	}{
		{
			name:     "empty delivery",
			delivery: &cart.Delivery{},
			want:     false,
		},
		{
			name:     "delivery with items but without discounts",
			delivery: testutils.BuildDeliveryWithoutDiscounts,
			want:     false,
		},
		{
			name:     "delivery with items with discounts",
			delivery: testutils.BuildDeliveryWithDiscounts,
			want:     true,
		},
		{
			name:     "delivery with items with duplicate discounts and shipping discounts",
			delivery: testutils.BuildDeliveryWithDuplicateDiscountsAndShippingDiscounts,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.delivery.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("Delivery.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name string
		item *cart.Item
		want cart.AppliedDiscounts
	}{
		{
			name: "no discounts on item",
			item: &cart.Item{
				AppliedDiscounts: cart.AppliedDiscounts{},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "multiple different discounts on item",
			item: testutils.BuildItemWithDiscounts,
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
				{
					CampaignCode: "code-7",
					Label:        "title-7",
					Type:         "type-7",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    6,
					CustomAttributes: map[string]interface{}{
						"attr1": 3,
						"attr2": 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.item.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Item.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_HasDiscounts(t *testing.T) {
	tests := []struct {
		name string
		item *cart.Item
		want bool
	}{
		{
			name: "no discounts on item",
			item: &cart.Item{},
			want: false,
		},
		{
			name: "multiple discounts on item",
			item: &cart.Item{AppliedDiscounts: cart.AppliedDiscounts{{
				CampaignCode: "code-1",
				Label:        "title-1",
				Type:         "type-1",
			}}},
			want: true,
		},
		{
			name: "duplicate discounts on item",
			item: testutils.BuildItemWithDuplicateDiscounts,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.item.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("Item.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShippingItem_MergeDiscounts(t *testing.T) {
	tests := []struct {
		name     string
		shipping *cart.ShippingItem
		want     cart.AppliedDiscounts
	}{
		{
			name: "no discounts on shipping",
			shipping: &cart.ShippingItem{
				AppliedDiscounts: cart.AppliedDiscounts{},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name:     "multiple discounts on shipping",
			shipping: testutils.BuildShippingItemWithDiscounts,
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.shipping.MergeDiscounts()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShippingItem.MergeDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShippingItem_HasDiscounts(t *testing.T) {
	tests := []struct {
		name     string
		shipping *cart.ShippingItem
		want     bool
	}{
		{
			name:     "no discounts on shipping",
			shipping: &cart.ShippingItem{},
			want:     false,
		},
		{
			name:     "multiple discounts on shipping",
			shipping: testutils.BuildShippingItemWithDiscounts,
			want:     true,
		},
		{
			name:     "duplicate discounts on shipping",
			shipping: testutils.BuildShippingItemWithDuplicateDiscounts,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.shipping.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("ShippingItem.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppliedDiscounts_ByCampaignCode(t *testing.T) {
	type args struct {
		campaignCode string
	}
	tests := []struct {
		name      string
		args      args
		discounts cart.AppliedDiscounts
		want      cart.AppliedDiscounts
	}{
		{
			name: "no match for filter",
			args: args{
				campaignCode: "code-3",
			},
			discounts: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
				},
				{
					CampaignCode: "code-1",
				},
				{
					CampaignCode: "code-2",
				},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "match for filter",
			args: args{
				campaignCode: "code-1",
			},
			discounts: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
				},
				{
					CampaignCode: "code-1",
				},
				{
					CampaignCode: "code-2",
				},
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
				},
				{
					CampaignCode: "code-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.discounts.ByCampaignCode(tt.args.campaignCode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppliedDiscounts.ByCampaignCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppliedDiscounts_ByType(t *testing.T) {
	type args struct {
		filterType string
	}
	tests := []struct {
		name      string
		args      args
		discounts cart.AppliedDiscounts
		want      cart.AppliedDiscounts
	}{
		{
			name: "no match for filter",
			args: args{
				filterType: "type-3",
			},
			discounts: cart.AppliedDiscounts{
				{
					Type: "type-1",
				},
				{
					Type: "type-2",
				},
				{
					Type: "type-1",
				},
			},
			want: cart.AppliedDiscounts{},
		},
		{
			name: "match for filter",
			args: args{
				filterType: "type-1",
			},
			discounts: cart.AppliedDiscounts{
				{
					Type: "type-1",
				},
				{
					Type: "type-2",
				},
				{
					Type: "type-1",
				},
			},
			want: cart.AppliedDiscounts{
				{
					Type: "type-1",
				},
				{
					Type: "type-1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.discounts.ByType(tt.args.filterType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppliedDiscounts.ByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppliedDiscounts_Sum(t *testing.T) {
	tests := []struct {
		name      string
		discounts cart.AppliedDiscounts
		want      domain.Price
		wantErr   bool
	}{
		{
			name: "sum of no discounts",
			want: domain.NewZero(""),
		},
		{
			name: "sum of multiple discounts",
			discounts: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    2,
				},
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "$"),
					SortOrder:    3,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-2",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(-12.0, "$"),
					SortOrder:    1,
				},
			},
			want: domain.NewFromFloat(-27.0, "$"),
		},
		{
			name: "sum of multiple discounts with different currencies",
			discounts: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-5.0, "$"),
					SortOrder:    0,
				},
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(-10.0, "EUR"),
					SortOrder:    0,
				},
			},
			want:    domain.NewZero(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.discounts.Sum()
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("AppliedDiscounts.Sum() gotErr %v, wantErr %v", gotErr != nil, tt.wantErr)
			}
			if !got.Equal(tt.want) {
				t.Errorf("AppliedDiscounts.Sum() = %v%v, want %v%v", got.Amount(), got.Currency(), tt.want.Amount(), tt.want.Currency())
			}
		})
	}
}
