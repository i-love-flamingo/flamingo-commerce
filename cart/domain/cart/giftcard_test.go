package cart_test

import (
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestCart_SumAppliedGiftCards(t *testing.T) {
	tests := []struct {
		name string
		cart *cart.Cart
		want domain.Price
	}{
		{
			name: "empty cart",
			cart: &cart.Cart{},
			want: domain.Price{}.GetPayable(),
		},
		{
			name: "cart with gift cards applied",
			cart: &cart.Cart{
				AppliedGiftCards: cart.AppliedGiftCards{
					{
						Applied: domain.NewFromFloat(15.0, "$"),
					},
					{
						Applied: domain.NewFromFloat(25.99, "$"),
					},
				},
			},
			want: domain.NewFromFloat(40.99, "$").GetPayable(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.cart.SumAppliedGiftCards()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.CanSumGiftCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_HasAppliedGiftCards(t *testing.T) {
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
			name: "cart with gift cards applied",
			cart: &cart.Cart{
				AppliedGiftCards: cart.AppliedGiftCards{
					{
						Applied: domain.NewFromFloat(15.0, "$"),
					},
					{
						Applied: domain.NewFromFloat(25.99, "$"),
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cart.HasAppliedGiftCards()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.HasAppliedGiftCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_SumGrandTotalWithGiftCards(t *testing.T) {
	tests := []struct {
		name string
		cart *cart.Cart
		want domain.Price
	}{
		{
			name: "empty cart",
			cart: &cart.Cart{},
			want: domain.Price{}.GetPayable(),
		},
		{
			name: "cart without discounts",
			cart: &cart.Cart{
				ID: "id-1",
				Deliveries: []cart.Delivery{
					{
						Cartitems: []cart.Item{
							{
								ID:            "test-1",
								Qty:           1,
								RowPriceGross: domain.NewFromFloat(10, "$").GetPayable(),
							},
						},
					},
				},
			},
			want: domain.NewFromFloat(10.0, "$").GetPayable(),
		},
		{
			name: "cart with gift cards applied",
			cart: &cart.Cart{
				ID: "id-1",
				Deliveries: []cart.Delivery{
					{
						Cartitems: []cart.Item{
							{
								ID:            "test-1",
								Qty:           1,
								RowPriceGross: domain.NewFromFloat(50.99, "$").GetPayable(),
							},
						},
					},
				},
				AppliedGiftCards: cart.AppliedGiftCards{
					{
						Applied: domain.NewFromFloat(15.0, "$").GetPayable(),
					},
					{
						Applied: domain.NewFromFloat(25.99, "$").GetPayable(),
					},
				},
			},
			want: domain.NewFromFloat(10.0, "$").GetPayable(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.cart.SumGrandTotalWithGiftCards()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.SumGrandTotalWithGiftCards = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_HasRemainingGiftCards(t *testing.T) {
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
			name: "cart without remaining discounts",
			cart: &cart.Cart{
				AppliedGiftCards: cart.AppliedGiftCards{
					{
						Applied: domain.NewFromFloat(1.0, "$"),
					},
					{
						Applied: domain.NewFromFloat(5.0, "$"),
					},
				},
			},
			want: false,
		},
		{
			name: "cart with remaining discounts",
			cart: &cart.Cart{
				AppliedGiftCards: cart.AppliedGiftCards{
					{
						Applied: domain.NewFromFloat(1.0, "$"),
					},
					{
						Applied: domain.NewFromFloat(5.0, "$"),
					},
					{
						Applied:   domain.NewFromFloat(15.0, "$"),
						Remaining: domain.NewFromFloat(1.0, "$"),
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cart.HasRemainingGiftCards()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.HasRemainingGiftCards = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppliedGiftCards_ByRemaining(t *testing.T) {
	tests := []struct {
		name  string
		cards cart.AppliedGiftCards
		want  cart.AppliedGiftCards
	}{
		{
			name: "no gift cards with remaining",
			cards: cart.AppliedGiftCards{
				{
					Applied: domain.NewFromFloat(1.0, "$"),
				},
				{
					Applied: domain.NewFromFloat(5.0, "$"),
				},
			},
			want: cart.AppliedGiftCards{},
		},
		{
			name: "gift cards with remaining",
			cards: cart.AppliedGiftCards{
				{
					Applied: domain.NewFromFloat(1.0, "$"),
				},
				{
					Applied:   domain.NewFromFloat(5.0, "$"),
					Remaining: domain.NewFromFloat(1.0, "$"),
				},
				{
					Applied:   domain.NewFromFloat(7.0, "$"),
					Remaining: domain.NewFromFloat(12.0, "$"),
				},
			},
			want: cart.AppliedGiftCards{
				{
					Applied:   domain.NewFromFloat(5.0, "$"),
					Remaining: domain.NewFromFloat(1.0, "$"),
				},
				{
					Applied:   domain.NewFromFloat(7.0, "$"),
					Remaining: domain.NewFromFloat(12.0, "$"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cards.ByRemaining()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppliedGiftCards.ByRemaining = %v, want %v", got, tt.want)
			}
		})
	}
}
