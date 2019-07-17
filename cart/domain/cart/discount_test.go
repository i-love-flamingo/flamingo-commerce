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
					builder := cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-1")
					delivery, _ := builder.Build()
					result = append(result, *delivery)
					builder = cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-2")
					result = append(result, *delivery)
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
					builder := cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-1")
					builder.AddItem(cart.Item{})
					builder.AddItem(cart.Item{})
					delivery, _ := builder.Build()
					result = append(result, *delivery)
					builder = cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-2")
					result = append(result, *delivery)
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
					delivery := testutils.BuildDeliveryWithDiscounts(t)
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithDiscounts(t)
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
			},
		},
		{
			name: "cart with different deliveries with items with discounts",
			cart: &cart.Cart{
				Deliveries: func() []cart.Delivery {
					result := make([]cart.Delivery, 0)
					delivery := testutils.BuildDeliveryWithDiscounts(t)
					result = append(result, *delivery)
					delivery = testutils.BuildAlternativeDeliveryWithAlternativeDiscounts(t)
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
					builder := cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-1")
					delivery, _ := builder.Build()
					result = append(result, *delivery)
					builder = cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-2")
					result = append(result, *delivery)
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
					builder := cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-1")
					builder.AddItem(cart.Item{})
					builder.AddItem(cart.Item{})
					delivery, _ := builder.Build()
					result = append(result, *delivery)
					builder = cart.DeliveryBuilder{}
					builder.SetDeliveryCode("code-2")
					result = append(result, *delivery)
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
					delivery := testutils.BuildDeliveryWithDiscounts(t)
					result = append(result, *delivery)
					delivery = testutils.BuildDeliveryWithDiscounts(t)
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
			delivery: testutils.BuildDeliveryWithoutDiscounts(t),
			want:     cart.AppliedDiscounts{},
		},
		{
			name:     "delivery with items with different discounts",
			delivery: testutils.BuildDeliveryWithDifferentDiscounts(t),
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
			delivery: testutils.BuildDeliveryWithoutDiscounts(t),
			want:     false,
		},
		{
			name:     "delivery with items with discounts",
			delivery: testutils.BuildDeliveryWithDiscounts(t),
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
			item: testutils.BuildItemWithDiscounts(t),
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
			item: func() *cart.Item {
				builder := cart.ItemBuilder{}
				builder.AddDiscount(cart.AppliedDiscount{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
				})
				item, _ := builder.Build()
				return item
			}(),
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
