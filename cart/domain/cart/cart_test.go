package cart_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
)

func Test_GetDeliveryCodes(t *testing.T) {
	cart := new(cartDomain.Cart)

	deliveryHome := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "home",
		},
		Cartitems: []cartDomain.Item{},
	}
	deliveryInFlight := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "inFlight",
		},
		Cartitems: []cartDomain.Item{},
	}
	deliveryWithoutItems := cartDomain.Delivery{
		DeliveryInfo: cartDomain.DeliveryInfo{
			Code: "withoutItems",
		},
	}

	cart.Deliveries = append(cart.Deliveries, deliveryHome)
	cart.Deliveries = append(cart.Deliveries, deliveryInFlight)
	cart.Deliveries = append(cart.Deliveries, deliveryWithoutItems)

	deliveryCodes := cart.GetDeliveryCodes()

	assert.Len(t, deliveryCodes, 2)
	assert.Contains(t, deliveryCodes, "home")
	assert.Contains(t, deliveryCodes, "inFlight")
}

/*
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
*/
