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
	assert.Equal(t, domain.NewFromInt(300, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
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
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
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
	assert.Equal(t, domain.NewFromInt(10, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())

	// verify first product charges
	// total items
	totalGCValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), totalGCValue.Value.FloatAmount())
	totalMainValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), totalMainValue.Value.FloatAmount())
	// shipping items
	shippingGCValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), shippingGCValue.Value.FloatAmount())
	shippingMainValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), shippingMainValue.Value.FloatAmount())
	// cart items
	itemGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(2, 1, "€").FloatAmount(), itemGCValue.Value.FloatAmount())
	itemMainValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), itemMainValue.Value.FloatAmount())

	// verify second product charges
	totalGCValue = selection.ItemSplit().TotalItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), totalGCValue.Value.FloatAmount())
	totalMainValue = selection.ItemSplit().TotalItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), totalMainValue.Value.FloatAmount())
	// shipping items
	shippingGCValue = selection.ItemSplit().ShippingItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), shippingGCValue.Value.FloatAmount())
	shippingMainValue = selection.ItemSplit().ShippingItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), shippingMainValue.Value.FloatAmount())
	// cart items
	itemGCValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(8, 1, "€").FloatAmount(), itemGCValue.Value.FloatAmount())
	itemMainValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(0, 1, "€").FloatAmount(), itemMainValue.Value.FloatAmount())
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
	assert.Equal(t, domain.NewFromInt(10, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(2, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
	// calculation:
	// cartotal: 12, giftcardtotal = 10
	// ratio: giftcardtotal / cartotal = 10 / 12
	// item 1: giftcard: 4 * ratio = 3.33, remaining: 4 - 3.33  = 0.67 (remaining giftcard: 10 - 3.33 = 6.67)
	// item 2: giftcard: 8 * ratio = 6.67, remaining: 8 - 6.67  = 1.33 (remaining giftcard: 6.67 - 6.67 = 0)
	// verify first product charges
	relativeGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(333, 100, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(67, 100, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
	// verfiy second product charges
	relativeGCValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(667, 100, "€").FloatAmount(), relativeGCValue.Value.FloatAmount())
	relativeMainValue = selection.ItemSplit().CartItems["2"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(133, 100, "€").FloatAmount(), relativeMainValue.Value.FloatAmount())
}

func Test_PayCompleteCartWithGiftcards(t *testing.T) {
	pricedItems := PricedItems{
		cartItems: make(map[string]domain.Price),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(4, 1, "€")
	pricedItems.cartItems["2"] = domain.NewFromInt(8, 1, "€")
	appliedGc := []AppliedGiftCard{
		{
			Applied: domain.NewFromInt(12, 1, "€"),
		},
	}
	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
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
	pricedItems := PricedItems{
		totalItems:    make(map[string]domain.Price),
		shippingItems: make(map[string]domain.Price),
		cartItems:     make(map[string]domain.Price),
	}
	pricedItems.cartItems["1"] = domain.NewFromInt(300099, 100, "€")
	pricedItems.shippingItems["1"] = domain.NewFromInt(88895, 100, "€")
	pricedItems.totalItems["1"] = domain.NewFromInt(1200095, 100, "€")
	appliedGc := []AppliedGiftCard{
		{
			Applied: domain.NewFromInt(50, 1, "€"),
		},
	}
	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(50, 1, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(1584089, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
	// calculation:
	// cartotal: 15888, giftcardtotal = 50
	// ratio: giftcardtotal / cartotal = 50 / 15890.89
	// total items giftcard    amount: ratio * 12000.95 = 37.76
	// shipping items giftcard amount: ratio * 888.95 = 2.8
	// cart items giftcard     amount: ratio * 3000.99 = 9.44
	// Sum giftcard: 37.76 + 2.8 + 9.44 = 50

	// verify total item charges
	totalGCValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(3776, 100, "€").FloatAmount(), totalGCValue.Value.FloatAmount())
	totalMainValue := selection.ItemSplit().TotalItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(1196319, 100, "€").FloatAmount(), totalMainValue.Value.FloatAmount())
	// verify shipping item charges
	shippingGCValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(280, 100, "€").FloatAmount(), shippingGCValue.Value.FloatAmount())
	shippingMainValue := selection.ItemSplit().ShippingItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(88615, 100, "€").FloatAmount(), shippingMainValue.Value.FloatAmount())
	// verify cart item charges
	itemGCValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeGiftCard)
	assert.Equal(t, domain.NewFromInt(944, 100, "€").FloatAmount(), itemGCValue.Value.FloatAmount())
	itemMainValue := selection.ItemSplit().CartItems["1"].ChargesByType().GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, domain.NewFromInt(299155, 100, "€").FloatAmount(), itemMainValue.Value.FloatAmount())

}
