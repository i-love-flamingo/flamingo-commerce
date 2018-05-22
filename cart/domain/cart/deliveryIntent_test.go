package cart

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.aoe.com/flamingo/framework/flamingo"
)

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

	intent = builder.BuildDeliveryIntent("pickup_collection-point_locpoint_1-2_3")
	assert.Equal(t, "pickup_collection-point_locpoint_1-2_3", intent.String())
	assert.Equal(t, "locpoint_1-2_3", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERY_METHOD_PICKUP, intent.Method)
	assert.Equal(t, DELIVERYLOCATION_TYPE_COLLECTIONPOINT, intent.DeliveryLocationType)

	intent = builder.BuildDeliveryIntent("pickup_autodetect")
	assert.Equal(t, "pickup_autodetect", intent.String())
	assert.Equal(t, "", intent.DeliveryLocationCode)
	assert.Equal(t, DELIVERY_METHOD_PICKUP, intent.Method)
}
