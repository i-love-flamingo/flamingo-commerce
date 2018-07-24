package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCartItem(t *testing.T) {
	var cart = new(Cart)

	cartItem := Item{MarketplaceCode: "code1", Qty: 5}
	cart.Cartitems = append(cart.Cartitems, cartItem)

	found, nr := cart.HasItem("code1", "")
	assert.True(t, found)
	assert.Equal(t, nr, 1)

	item, err := cart.GetByLineNr(1)
	assert.NoError(t, err)
	assert.Equal(t, &cartItem, item)
}

func TestCart_HasDeliveryIntentArrival(t *testing.T) {
	var cart = new(Cart)

	cart.Cartitems = append(cart.Cartitems, getItemWithArrivalIntent())

	resultArrival := cart.HasItemWithIntent("arrival")
	assert.True(t, resultArrival)
}

func TestCart_HasDeliveryIntentDeparture(t *testing.T) {
	var cart = new(Cart)

	cart.Cartitems = append(cart.Cartitems, getItemWithDepartureIntent())

	resultDeparture := cart.HasItemWithIntent("departure")
	assert.True(t, resultDeparture)
}

func TestCart_HasNoMixedCart(t *testing.T) {
	var cart = new(Cart)

	cart.Cartitems = append(cart.Cartitems, getItemWithDepartureIntent())

	resultNoMixedCart := cart.HasMixedCart()
	assert.False(t, resultNoMixedCart)
}

func TestCart_HasMixedCart(t *testing.T) {
	var cart = new(Cart)

	cart.Cartitems = append(cart.Cartitems, getItemWithDepartureIntent())
	cart.Cartitems = append(cart.Cartitems, getItemWithArrivalIntent())

	resultMixedCart := cart.HasMixedCart()
	assert.True(t, resultMixedCart)
}

func getItemWithDepartureIntent() Item {
	deliveryIntentDeparture := DeliveryIntent{
		Method: "departure",
	}
	return Item{
		MarketplaceCode:        "codeDeparture",
		OriginalDeliveryIntent: &deliveryIntentDeparture,
		Qty: 1,
	}
}

func getItemWithArrivalIntent() Item {
	deliveryIntentArrival := DeliveryIntent{
		Method: "arrival",
	}
	return Item{
		MarketplaceCode:        "codeArrival",
		OriginalDeliveryIntent: &deliveryIntentArrival,
		Qty: 1,
	}
}
