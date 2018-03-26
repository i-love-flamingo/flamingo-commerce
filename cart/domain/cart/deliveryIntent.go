package cart

import (
	"strings"

	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	DeliveryIntentBuilder struct {
		Logger flamingo.Logger `inject:""`
	}

	//DeliveryIntent - represents the Intent for delivery
	DeliveryIntent struct {
		Method               string
		DeliveryLocationCode string
		DeliveryLocationType string
	}

	PickUpDetectionService interface {
		Detect(product productDomain.BasicProduct, request AddRequest) (locationCode string, locationType string, err error)
	}
)

//BuildDeliveryIntent - gets DeliveryIntent by string representation
func (b *DeliveryIntentBuilder) BuildDeliveryIntent(representation string) DeliveryIntent {
	if representation == DELIVERY_METHOD_DELIVERY {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_DELIVERY,
		}
	}

	//"pickup-autodetect"
	if representation == "store-autodetect" {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_PICKUP,
		}
	}

	intentParts := strings.SplitN(representation, "_", 3)
	if len(intentParts) != 3 {
		b.Logger.WithField("category", "cart").WithField("subcategory", "DeliveryIntentBuilder").Warnf("Unknown IntentString %v", representation)
		return DeliveryIntent{
			Method: DELIVERY_METHOD_UNSPECIFIED,
		}
	}
	if intentParts[0] == DELIVERY_METHOD_PICKUP {
		if intentParts[1] == DELIVERYLOCATION_TYPE_STORE {
			return DeliveryIntent{
				Method:               DELIVERY_METHOD_PICKUP,
				DeliveryLocationCode: intentParts[2],
				DeliveryLocationType: DELIVERYLOCATION_TYPE_STORE,
			}
		}
		if intentParts[1] == DELIVERYLOCATION_TYPE_COLLECTIONPOINT {
			return DeliveryIntent{
				Method:               DELIVERY_METHOD_PICKUP,
				DeliveryLocationCode: intentParts[2],
				DeliveryLocationType: DELIVERYLOCATION_TYPE_COLLECTIONPOINT,
			}
		}
	}
	b.Logger.WithField("category", "cart").WithField("subcategory", "DeliveryIntentBuilder").Warnf("Unknown IntentString %v", representation)
	return DeliveryIntent{
		Method: DELIVERY_METHOD_UNSPECIFIED,
	}
}

//GetDeliveryInfo - gets the (initial) GetDeliveryInfo that is meant by this Intent
func (di *DeliveryIntent) GetDeliveryInfo() DeliveryInfo {
	return DeliveryInfo{
		Method: di.Method,
		DeliveryLocation: DeliveryLocation{
			Code: di.DeliveryLocationCode,
			Type: di.DeliveryLocationType,
		},
	}
}

//String - Gets String representation of DeliveryIntent
func (di *DeliveryIntent) String() string {
	if di.Method == DELIVERY_METHOD_PICKUP {
		return di.Method + "_" + di.DeliveryLocationType + "_" + di.DeliveryLocationCode
	}
	return di.Method
}
