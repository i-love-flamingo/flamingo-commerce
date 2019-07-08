package cart_test

import (
	"reflect"
	"sort"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

func TestItem_CollectDiscounts(t *testing.T) {
	tests := []struct {
		name string
		item *cart.Item
		want []*cart.AppliedDiscount
	}{
		{
			name: "no discounts on item",
			item: &cart.Item{},
			want: []*cart.AppliedDiscount{},
		},
		{
			name: "multiple different discounts on item",
			item: func() *cart.Item {
				builder := cart.ItemBuilder{}
				builder.AddDiscount(cart.ItemDiscount{
					Code:  "code-1",
					Title: "title-1",
					Type:  "type-1",
				})
				builder.AddDiscount(cart.ItemDiscount{
					Code:  "code-1",
					Title: "title-2",
					Type:  "type-1",
				})
				builder.AddDiscount(cart.ItemDiscount{
					Code:  "code-1",
					Title: "title-1",
					Type:  "type-2",
				})
				item, _ := builder.Build()
				return item
			}(),
			want: []*cart.AppliedDiscount{
				{
					Code:  "code-1",
					Title: "title-1",
					Type:  "type-1",
				},
				{
					Code:  "code-1",
					Title: "title-2",
					Type:  "type-1",
				},
				{
					Code:  "code-1",
					Title: "title-1",
					Type:  "type-2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.CollectDiscounts()
			// we need to sort result to circumvent implementation changes in order
			sort.Sort(cart.ByCode(got))
			sort.Sort(cart.ByCode(tt.want))
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
				builder.AddDiscount(cart.ItemDiscount{
					Code:  "code-1",
					Title: "title-1",
					Type:  "type-1",
				})
				item, _ := builder.Build()
				return item
			}(),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.HasAppliedDiscounts(); got != tt.want {
				t.Errorf("Item.HasDiscounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
