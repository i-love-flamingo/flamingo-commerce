package cart_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

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

func TestAppliedGiftCard_Total(t *testing.T) {
	// Applied and Remaining with different values but same currency returns a correct total
	giftCard := cart.AppliedGiftCard{
		Applied:   domain.NewFromFloat(10.5, "$"),
		Remaining: domain.NewFromFloat(4.5, "$"),
	}
	total, err := giftCard.Total()
	assert.Nil(t, err)
	assert.Equal(t, true, domain.NewFromFloat(15, "$").Equal(total))

	// Applied of Zero and Remaining with a value and currency returns a correct total
	giftCard = cart.AppliedGiftCard{
		Applied:   domain.NewZero("$"),
		Remaining: domain.NewFromFloat(10, "$"),
	}
	total, err = giftCard.Total()
	assert.Nil(t, err)
	assert.Equal(t, true, domain.NewFromFloat(10, "$").Equal(total))

	// Applied with a value and currency and Remaining of Zero returns a correct total
	giftCard = cart.AppliedGiftCard{
		Applied:   domain.NewFromFloat(5, "$"),
		Remaining: domain.NewZero("$"),
	}
	total, err = giftCard.Total()
	assert.Nil(t, err)
	assert.Equal(t, true, domain.NewFromFloat(5, "$").Equal(total))

	// Applied and Remaining with different values and different currencies returns an error and the price of Remaining
	giftCard = cart.AppliedGiftCard{
		Applied:   domain.NewFromFloat(10.5, "$"),
		Remaining: domain.NewFromFloat(4.5, "€"),
	}
	total, err = giftCard.Total()
	assert.NotNil(t, err)
	assert.Equal(t, true, domain.NewFromFloat(4.5, "€").Equal(total))
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

func TestAppliedGiftCards_GiftCardByCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name      string
		cards     cart.AppliedGiftCards
		args      args
		wantCard  *cart.AppliedGiftCard
		wantFound bool
	}{
		{
			name:  "no gift cards so none is found",
			cards: cart.AppliedGiftCards{},
			args: args{
				"some-code",
			},
			wantCard:  nil,
			wantFound: false,
		},
		{
			name: "same code so gift card is found",
			cards: cart.AppliedGiftCards{
				{
					Code: "some-code",
				},
			},
			args: args{
				"some-code",
			},
			wantCard: &cart.AppliedGiftCard{
				Code: "some-code",
			},
			wantFound: true,
		},
		{
			name: "different code so no gift card is found",
			cards: cart.AppliedGiftCards{
				{
					Code: "some-code",
				},
			},
			args: args{
				"some-other-code",
			},
			wantCard:  nil,
			wantFound: false,
		},
		{
			name: "multiple gift cards and one code is found",
			cards: cart.AppliedGiftCards{
				{
					Code: "some-code",
				},
				{
					Code: "some-other-code",
				},
			},
			args: args{
				"some-code",
			},
			wantCard: &cart.AppliedGiftCard{
				Code: "some-code",
			},
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCard, gotFound := tt.cards.GiftCardByCode(tt.args.code)
			if !reflect.DeepEqual(gotCard, tt.wantCard) {
				t.Errorf("GiftCardByCode() gotCard = %v, want %v", gotCard, tt.wantCard)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GiftCardByCode() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
