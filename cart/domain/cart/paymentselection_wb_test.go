package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func getPaymentMethodMapping(t *testing.T) map[string]string {
	t.Helper()
	return map[string]string{
		domain.ChargeTypeMain:     "creditcard",
		domain.ChargeTypeGiftCard: "giftcard",
	}
}

func Test_CanBuildSimpleSelectionFromCard(t *testing.T) {

	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "delcode",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(199, 100, "€"),
					},
					{
						ID:            "2",
						RowPriceGross: domain.NewFromInt(299, 100, "€"),
					},
				},
				ShippingItem: ShippingItem{
					PriceNet: domain.NewFromInt(7, 1, "€"),
				},
			},
		},
	}
	selection, _ := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
}

func Test_CanBuildSimpleSelectionWithGiftCard_NoGc(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "delcode",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(199, 100, "€"),
					},
					{
						ID:            "2",
						RowPriceGross: domain.NewFromInt(299, 100, "€"),
					},
				},
				ShippingItem: ShippingItem{
					PriceNet: domain.NewFromInt(7, 1, "€"),
				},
			},
		},
		AppliedGiftCards: AppliedGiftCards{},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
}

func Test_CanBuildSimpleSelectionWithGiftCard(t *testing.T) {

	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "delcode",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(199, 100, "€"),
					},
					{
						ID:            "2",
						RowPriceGross: domain.NewFromInt(299, 100, "€"),
					},
				},
				ShippingItem: ShippingItem{
					PriceNet: domain.NewFromInt(7, 1, "€"),
				},
			},
		},
		AppliedGiftCards: AppliedGiftCards{
			{
				Code:    "code-1",
				Applied: domain.NewFromInt(100, 100, "€"),
			},
			{
				Code:    "code-2",
				Applied: domain.NewFromInt(200, 100, "€"),
			},
		},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
	want := domain.NewFromInt(199, 100, "€").FloatAmount()
	got := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Price.FloatAmount()
	assert.Equal(t, want, got)

	want = domain.NewFromInt(0, 0, "€").FloatAmount()
	got = selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Price.FloatAmount()
	assert.Equal(t, want, got)

	want = domain.NewFromInt(101, 100, "€").FloatAmount()
	got = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Price.FloatAmount()
	assert.Equal(t, want, got)

	want = domain.NewFromInt(198, 100, "€").FloatAmount()
	got = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Price.FloatAmount()
	assert.Equal(t, want, got)

}

func Test_CanBuildSimpleSelectionWithGiftCard2(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "delcode",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(199, 100, "€"),
					},
					{
						ID:            "2",
						RowPriceGross: domain.NewFromInt(299, 100, "€"),
					},
				},
				ShippingItem: ShippingItem{
					PriceNet: domain.NewFromInt(7, 1, "€"),
				},
			},
		},
		AppliedGiftCards: AppliedGiftCards{
			{
				Applied: domain.NewFromInt(1198, 100, "€"),
			},
		},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(0, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
}

func Test_CanCalculateGiftCardChargeWithRest(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "delcode",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(4, 1, "€"),
					},
					{
						ID:            "2",
						RowPriceGross: domain.NewFromInt(8, 1, "€"),
					},
				},
			},
		},
		AppliedGiftCards: AppliedGiftCards{
			{
				Applied: domain.NewFromInt(10, 1, "€"),
			},
		},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	// verfiy complete cart splits
	assert.Equal(t, domain.NewFromInt(10, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(2, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())

	// verify first product charges
	relativeGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(4, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
	// verfiy second product charges
	relativeGCValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(6, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(2, 1, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
}

func Test_PayCompleteCartWithGiftCards(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "delcode",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(4, 1, "€"),
					},
					{
						ID:            "2",
						RowPriceGross: domain.NewFromInt(8, 1, "€"),
					},
				},
			},
		},
		AppliedGiftCards: AppliedGiftCards{
			{
				Applied: domain.NewFromInt(12, 1, "€"),
			},
		},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(12, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
	// item 1 is completely paid for
	relativeGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(4, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	// item 2 is completely paid for
	relativeGCValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(8, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
}

func Test_CartWithExpensiveItems(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "1",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(300099, 100, "€"),
					},
				},
				ShippingItem: ShippingItem{
					Title:    "1",
					PriceNet: domain.NewFromInt(88895, 100, "€"),
				},
			},
		},
		Totalitems: []Totalitem{
			{
				Code:  "1",
				Title: "1",
				Price: domain.NewFromInt(1200095, 100, "€"),
			},
		},
		AppliedGiftCards: []AppliedGiftCard{
			{
				Code:    "code-1",
				Applied: domain.NewFromInt(50, 1, "€"),
			},
			{
				Code:    "code-2",
				Applied: domain.NewFromInt(50, 1, "€"),
			},
		},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(100, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(1579089, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())

	// verify total item charges
	totalGCValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(0, 0, "€").FloatAmount(), totalGCValue.Value.FloatAmount())
	totalMainValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(1200095, 100, "€").FloatAmount(), totalMainValue.Value.FloatAmount())
	// verify shipping item charges
	shippingGCValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(0, 0, "€").FloatAmount(), shippingGCValue.Value.FloatAmount())
	shippingMainValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(88895, 100, "€").FloatAmount(), shippingMainValue.Value.FloatAmount())
	// verify cart item charges
	itemGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(100, 1, "€").FloatAmount(), itemGCValue.Value.FloatAmount())
	itemMainValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(290099, 100, "€").FloatAmount(), itemMainValue.Value.FloatAmount())
}

func Test_CartWithShipping(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "1",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(150, 1, "€"),
					},
				},
				ShippingItem: ShippingItem{
					Title:    "1",
					PriceNet: domain.NewFromInt(99, 1, "€"),
				},
			},
		},
		AppliedGiftCards: []AppliedGiftCard{
			{
				Code:    "code-1",
				Applied: domain.NewFromInt(120, 1, "€"),
			},
			{
				Code:    "code-2",
				Applied: domain.NewFromInt(40, 1, "€"),
			},
		},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(160, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(89, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())

	// verify cart item charges
	itemGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(150, 1, "€").FloatAmount(), itemGCValue.Value.FloatAmount())
	itemMainValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), itemMainValue.Value.FloatAmount())

	appliedGiftCardCharges := selection.ItemSplit().CartItems["1"].ChargesByType().GetAllByType(domain.ChargeTypeGiftCard)
	assert.Len(t, appliedGiftCardCharges, 2)

	cq := domain.ChargeQualifier{Type: domain.ChargeTypeGiftCard, Reference: "code-1"}
	assert.Equal(t, 120.0, selection.ItemSplit().CartItems["1"].ChargesByType().GetByChargeQualifierForced(cq).Price.FloatAmount())

	cq = domain.ChargeQualifier{Type: domain.ChargeTypeGiftCard, Reference: "code-2"}
	assert.Equal(t, 30.0, selection.ItemSplit().CartItems["1"].ChargesByType().GetByChargeQualifierForced(cq).Price.FloatAmount())

	// verify shipping item charges
	shippingGCValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(10, 1, "€").FloatAmount(), shippingGCValue.Value.FloatAmount())
	shippingMainValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(89, 1, "€").FloatAmount(), shippingMainValue.Value.FloatAmount())

	cq = domain.ChargeQualifier{Type: domain.ChargeTypeGiftCard, Reference: "code-2"}
	assert.Equal(t, 10.0, selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByChargeQualifierForced(cq).Price.FloatAmount())
	assert.Equal(t, 120.0, cart.AppliedGiftCards[0].Applied.FloatAmount())
	assert.Equal(t, 40.0, cart.AppliedGiftCards[1].Applied.FloatAmount())
}

func Test_CreateSimplePaymentWithoutGiftCards(t *testing.T) {
	cart := Cart{
		Deliveries: []Delivery{
			{
				DeliveryInfo: DeliveryInfo{
					Code: "1",
				},
				Cartitems: []Item{
					{
						ID:            "1",
						RowPriceGross: domain.NewFromInt(50, 100, "€"),
					},
				},
				ShippingItem: ShippingItem{
					Title:    "1",
					PriceNet: domain.NewFromInt(20, 100, "€"),
				},
			},
		},
		Totalitems: []Totalitem{
			{
				Code:  "1",
				Title: "1",
				Price: domain.NewFromInt(50, 100, "€"),
			},
		},
		AppliedGiftCards: AppliedGiftCards{},
	}
	selection, err := NewDefaultPaymentSelection("gateyway", getPaymentMethodMapping(t), cart)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(120, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
}
