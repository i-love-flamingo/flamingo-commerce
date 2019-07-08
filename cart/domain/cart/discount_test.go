package cart_test

import (
	"reflect"
	"sort"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// ByCode implements sort.Interface for []AppliedDiscount based on code
	ByCode cart.AppliedDiscounts
)

// implementations for sort interface

func (a ByCode) Len() int {
	return len(a)
}

func (a ByCode) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByCode) Less(i, j int) bool {
	return a[i].CampaignCode < a[j].CampaignCode
}

// buildItemWithDiscounts helper for item building
func buildItemWithDiscounts() *cart.Item {
	builder := cart.ItemBuilder{}
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-1",
		Label:        "title-1",
		Type:         "type-1",
		Applied:      domain.NewFromFloat(10.0, "$"),
	})
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-2",
		Label:        "title-2",
		Type:         "type-1",
		Applied:      domain.NewFromFloat(15.0, "$"),
	})
	builder.AddDiscount(cart.AppliedDiscount{
		CampaignCode: "code-3",
		Label:        "title-1",
		Type:         "type-2",
		Applied:      domain.NewFromFloat(5.0, "$"),
	})
	item, _ := builder.Build()
	return item
}

// buildDeliverxWithDiscounts helper for delivery building
// Adds an item with discount twice
// This means when discounts are summed up (based on type + delivery)
// The amount should be added to the previous discount
func buildDeliveryWithDiscounts() *cart.Delivery {
	builder := cart.DeliveryBuilder{}
	builder.SetDeliveryCode("code")
	builder.AddItem(*buildItemWithDiscounts())
	builder.AddItem(*buildItemWithDiscounts())
	// add items with discounts
	delivery, _ := builder.Build()
	return delivery
}

func TestCart_CollectDiscounts(t *testing.T) {
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
					delivery := buildDeliveryWithDiscounts()
					result = append(result, *delivery)
					delivery = buildDeliveryWithDiscounts()
					result = append(result, *delivery)
					return result
				}(),
			},
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(40.0, "$"),
				},
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(60.0, "$"),
				},
				{
					CampaignCode: "code-3",
					Label:        "title-1",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(20.0, "$"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.cart.CollectDiscounts()
			// we need to sort result to circumvent implementation changes in order
			sort.Sort(ByCode(got))
			sort.Sort(ByCode(tt.want))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.CollectDiscounts() = %v, want %v", got, tt.want)
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
					delivery := buildDeliveryWithDiscounts()
					result = append(result, *delivery)
					delivery = buildDeliveryWithDiscounts()
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

func TestDelivery_CollectDiscounts(t *testing.T) {
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
			name: "delivery with items but without discounts",
			delivery: func() *cart.Delivery {
				builder := cart.DeliveryBuilder{}
				builder.AddItem(cart.Item{})
				builder.AddItem(cart.Item{})
				builder.SetDeliveryCode("code")
				delivery, _ := builder.Build()
				return delivery
			}(),
			want: cart.AppliedDiscounts{},
		},
		{
			name:     "delivery with items with discounts",
			delivery: buildDeliveryWithDiscounts(),
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(20.0, "$"),
				},
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(30.0, "$"),
				},
				{
					CampaignCode: "code-3",
					Label:        "title-1",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(10.0, "$"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.delivery.CollectDiscounts()
			// we need to sort result to circumvent implementation changes in order
			sort.Sort(ByCode(got))
			sort.Sort(ByCode(tt.want))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delivery.CollectDiscounts() = %v, want %v", got, tt.want)
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
			name: "delivery with items but without discounts",
			delivery: func() *cart.Delivery {
				builder := cart.DeliveryBuilder{}
				builder.AddItem(cart.Item{})
				builder.AddItem(cart.Item{})
				builder.SetDeliveryCode("code")
				delivery, _ := builder.Build()
				return delivery
			}(),
			want: false,
		},
		{
			name:     "delivery with items with discounts",
			delivery: buildDeliveryWithDiscounts(),
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

func TestItem_CollectDiscounts(t *testing.T) {
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
			item: buildItemWithDiscounts(),
			want: cart.AppliedDiscounts{
				{
					CampaignCode: "code-1",
					Label:        "title-1",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(10.0, "$"),
				},
				{
					CampaignCode: "code-2",
					Label:        "title-2",
					Type:         "type-1",
					Applied:      domain.NewFromFloat(15.0, "$"),
				},
				{
					CampaignCode: "code-3",
					Label:        "title-1",
					Type:         "type-2",
					Applied:      domain.NewFromFloat(5.0, "$"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.item.CollectDiscounts()
			// we need to sort result to circumvent implementation changes in order
			sort.Sort(ByCode(got))
			sort.Sort(ByCode(tt.want))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Item.CollectDiscounts() = %v, want %v", got, tt.want)
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
