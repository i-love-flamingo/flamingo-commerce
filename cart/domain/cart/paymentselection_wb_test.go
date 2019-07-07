package cart

import (
	"github.com/stretchr/testify/assert"
	"testing"

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
		AppliedGiftCard{
			Applied: domain.NewFromInt(100,100,"€"),
		},
		AppliedGiftCard{
			Applied: domain.NewFromInt(200,100,"€"),
		},
	}

	selection, err := NewPaymentSelectionWithGiftCard("gateyway", "method", pricedItems, appliedGc)
	assert.NoError(t, err)
	assert.Equal(t, domain.NewFromInt(1198, 100, "€").FloatAmount(), selection.TotalValue().FloatAmount())
	assert.Equal(t, domain.NewFromInt(300, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(ChargeTypeGiftCard).Value.FloatAmount())
	assert.Equal(t, domain.NewFromInt(898, 100, "€").FloatAmount(), selection.CartSplit().ChargesByType().GetByTypeForced(domain.ChargeTypeMain).Value.FloatAmount())
}
