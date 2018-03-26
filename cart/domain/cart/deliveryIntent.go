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
func BuildDeliveryIntent(representation string) DeliveryIntent {
	if representation == DELIVERY_METHOD_DELIVERY {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_DELIVERY,
		}
	}

	//"pickup-autodetect"
	if representation == "pickup-autodetect" {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_PICKUP,
		}
	}

	intentParts := strings.SplitN(representation, "_", 2)
	if len(intentParts) != 2 {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_UNSPECIFIED,
		}
	}
	if intentParts[0] == DELIVERY_METHOD_PICKUP {
		return DeliveryIntent{
			Method:               DELIVERY_METHOD_PICKUP,
			DeliveryLocationCode: intentParts[1],
			DeliveryLocationType: DELIVERYLOCATION_TYPE_STORE,
		}
	}
	if intentParts[0] == DELIVERYLOCATION_TYPE_COLLECTIONPOINT {
		return DeliveryIntent{
			Method:               DELIVERY_METHOD_PICKUP,
			DeliveryLocationCode: intentParts[1],
			DeliveryLocationType: DELIVERYLOCATION_TYPE_COLLECTIONPOINT,
		}
	}
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
	return di.Method + "_" + di.DeliveryLocationCode
}
