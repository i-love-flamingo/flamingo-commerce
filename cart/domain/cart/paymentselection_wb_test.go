package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/price/domain"
)

func Test_CanBuildSimpleSelectionFromCard(t *testing.T) {
	pricedItems := PricedItems{
		cartItems:     make(map[string]domain.Price),
		shippingItems: make(map[string]domain.Price, 1),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(199, 100, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(299, 100, "€")
	pricedItems.shippingItems["delcode"] = domain.NewFromInt(7, 1, "€")
	selection := NewSimplePaymentSelection("gateyway", "method", pricedItems)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
}

func Test_CanBuildSimpleSelectionWithGiftCard_NoGc(t *testing.T) {
	pricedItems := PricedItems{
		cartItems:     make(map[string]domain.Price),
		shippingItems: make(map[string]domain.Price, 1),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(199, 100, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(299, 100, "€")
	pricedItems.shippingItems["delcode"] = domain.NewFromInt(7, 1, "€")
	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, nil)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
}

func Test_CanBuildSimpleSelectionWithGiftCard(t *testing.T) {
	pricedItems := PricedItems{
		cartItems:     make(map[string]domain.Price),
		shippingItems: make(map[string]domain.Price, 1),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(199, 100, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(299, 100, "€")
	pricedItems.shippingItems["delcode"] = domain.NewFromInt(7, 1, "€")

	//Apply 3 € GC
	appliedGc := []AppliedGiftCard{
		{
			Applied: domain.NewFromInt(100, 100, "€"),
		},
		{
			Applied: domain.NewFromInt(200, 100, "€"),
		},
	}

	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
	assert.Equal(t, domain.NewFromInt(300, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(898, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
}

func Test_CanBuildSimpleSelectionWithGiftCard2(t *testing.T) {
	pricedItems := PricedItems{
		cartItems:     make(map[string]domain.Price),
		shippingItems: make(map[string]domain.Price, 1),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(199, 100, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(299, 100, "€")
	pricedItems.shippingItems["delcode"] = domain.NewFromInt(7, 1, "€")

	//Apply 3 € GC
	appliedGc := []AppliedGiftCard{
		{
			Applied: domain.NewFromInt(1198, 100, "€"),
		},
	}

	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(0, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
}

func Test_CanCalculateGiftCardChargeRelativeToItem(t *testing.T) {
	pricedItems := PricedItems{
		cartItems: make(map[string]domain.Price),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(2, 1, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(8, 1, "€")

	appliedGc := []AppliedGiftCard{
		{
			Applied: domain.NewFromInt(5, 1, "€"),
		},
		{
			Applied: domain.NewFromInt(5, 1, "€"),
		},
	}
	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
	assert.NoError(t, err)
	// verfiy complete cart splits
	assert.Equal(t, domain.NewFromInt(10, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
	// verify first product charges
	relativeGCValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(2, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
	// verfiy second product charges
	relativeGCValue = selection.ItemSplit().TotalItems["2"].ChargesByType().GetByTypeForced(ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(8, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue = selection.ItemSplit().TotalItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
}

func Test_CanCalculateGiftCardChargeRelativeToItemWithRest(t *testing.T) {
	pricedItems := PricedItems{
		cartItems: make(map[string]domain.Price),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(4, 1, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(8, 1, "€")

	appliedGc := []AppliedGiftCard{
		{
			Applied: domain.NewFromInt(10, 1, "€"),
		},
	}
	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
	assert.NoError(t, err)
	// verfiy complete cart splits
	assert.Equal(t, domain.NewFromInt(10, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(02, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
	// verify first product charges
	relativeGCValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(4, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
	// verfiy second product charges
	relativeGCValue = selection.ItemSplit().TotalItems["2"].ChargesByType().GetByTypeForced(ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(6, 1, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue = selection.ItemSplit().TotalItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(2, 1, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
}
