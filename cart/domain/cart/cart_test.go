package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.aoe.com/flamingo/framework/flamingo"
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

func TestDeliveryIntent(t *testing.T) {
	builder := DeliveryIntentBuilder{
		Logger: flamingo.NullLogger{},
	}
	intent := builder.BuildDeliveryIntent("pickup_store_location_1")
	assert.Equal(t, "pickup_store_location_1", intent.String())
	assert.Equal(t, "location_1", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERYLOCATION_TYPE_STORE, intent.DeliveryLocationType)
	assert.Equal(t, DELIVERY_METHOD_PICKUP, intent.Method)

	intent = builder.BuildDeliveryIntent("")
	assert.Equal(t, DELIVERY_METHOD_UNSPECIFIED, intent.String())
	assert.Equal(t, "", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERY_METHOD_UNSPECIFIED, intent.Method, "empty intent string should by unspecified")

	intent = builder.BuildDeliveryIntent("lkjlkj")
	assert.Equal(t, DELIVERY_METHOD_UNSPECIFIED, intent.String())
	assert.Equal(t, "", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERY_METHOD_UNSPECIFIED, intent.Method, "random unvalid intent string should by unspecified")

	intent = builder.BuildDeliveryIntent("delivery")
	assert.Equal(t, "delivery", intent.String())
	assert.Equal(t, "", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERY_METHOD_DELIVERY, intent.Method)

	intent = builder.BuildDeliveryIntent("pickup_collection_locpoint_1-2_3")
	assert.Equal(t, "pickup_collection_locpoint_1-2_3", intent.String())
	assert.Equal(t, "locpoint_1-2_3", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERY_METHOD_PICKUP, intent.Method)
	assert.Equal(t, DELIVERYLOCATION_TYPE_COLLECTIONPOINT, intent.DeliveryLocationType)
}

func TestCartDeliveryInfo(t *testing.T) {

	//Prepare a Cart with one existing DeliveryInfo
	var cart = new(Cart)
	existingDeliveryInfo := DeliveryInfo{
		Method: DELIVERY_METHOD_PICKUP,
		DeliveryLocation: DeliveryLocation{
			Code: "code1",
			Type: DELIVERYLOCATION_TYPE_COLLECTIONPOINT,
		},
	}
	cart.DeliveryInfos = append(cart.DeliveryInfos, existingDeliveryInfo)

	//1. Test if the existing deliveryInfo gets returned on same intent
	deliveryInfoReference, err := cart.GetDeliveryMethodForIntent(DeliveryIntent{
		Method:               DELIVERY_METHOD_PICKUP,
		DeliveryLocationType: DELIVERYLOCATION_TYPE_COLLECTIONPOINT,
		DeliveryLocationCode: "code1",
	})

	assert.Nil(t, err)
	assert.Equal(t, &existingDeliveryInfo, deliveryInfoReference)

	//2. Test if the existing deliveryInfo gets NOT returned on some different intent
	deliveryInfoReference, err = cart.GetDeliveryMethodForIntent(DeliveryIntent{
		Method:               DELIVERY_METHOD_PICKUP,
		DeliveryLocationType: DELIVERYLOCATION_TYPE_COLLECTIONPOINT,
		DeliveryLocationCode: "code not here",
	})
	assert.Error(t, err)

}
