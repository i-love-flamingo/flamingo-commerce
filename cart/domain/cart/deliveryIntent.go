package cart

import (
	"strings"

	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	//DeliveryIntentBuilder - Factory
	DeliveryIntentBuilder struct {
		Logger flamingo.Logger `inject:""`
	}

	//DeliveryInfoBuilder - Factory
	DeliveryInfoBuilder interface {
		BuildDeliveryInfoUpdateCommand(ctx web.Context, decoratedCart *DecoratedCart) ([]DeliveryInfoUpdateCommand, error)
	}

	DefaultDeliveryInfoBuilder struct {
	}

	//DeliveryIntent - represents the Intent for delivery
	DeliveryIntent struct {
		Method                     string
		AutodetectDeliveryLocation bool
		DeliveryLocationCode       string
		DeliveryLocationType       string
	}

	PickUpDetectionService interface {
		Detect(product productDomain.BasicProduct, request AddRequest) (locationCode string, locationType string, err error)
	}
)

//BuildDeliveryInfoUpdateCommand - default implementation to get DeliveryInfo for cart. It is simply using the DeliverIntent on the Items
func (dib *DefaultDeliveryInfoBuilder) BuildDeliveryInfoUpdateCommand(ctx web.Context, decoratedCart *DecoratedCart) ([]DeliveryInfoUpdateCommand, error) {
	var updateCommands []DeliveryInfoUpdateCommand
	for _, cartitems := range decoratedCart.Cart.GetCartItemsByOriginalDeliveryIntent() {
		if len(cartitems) < 1 {
			continue
		}
		deliveryInfo := cartitems[0].OriginalDeliveryIntent.BuildDeliveryInfo()
		itemIds := make([]string, len(cartitems))
		for _, cartitem := range cartitems {
			itemIds = append(itemIds, cartitem.ID)
		}
		updateCommands = append(updateCommands, DeliveryInfoUpdateCommand{
			DeliveryInfo:    &deliveryInfo,
			AssignedItemIds: itemIds,
		})
	}
	return updateCommands, nil
}

//BuildDeliveryIntent - gets DeliveryIntent by string representation
func (b *DeliveryIntentBuilder) BuildDeliveryIntent(representation string) DeliveryIntent {
	if representation == "" {
		b.Logger.WithField("category", "cart").WithField("subcategory", "DeliveryIntentBuilder").Warnf("Empty IntentString ")
		return DeliveryIntent{
			Method: DELIVERY_METHOD_UNSPECIFIED,
		}
	}
	if representation == DELIVERY_METHOD_DELIVERY {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_DELIVERY,
		}
	}

	if representation == "pickup_autodetect" {
		return DeliveryIntent{
			Method: DELIVERY_METHOD_PICKUP,
			AutodetectDeliveryLocation: true,
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

//BuildDeliveryInfo - gets the (initial) DeliveryInfo that is meant by this Intent
func (di *DeliveryIntent) BuildDeliveryInfo() DeliveryInfo {
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
	if di.Method == DELIVERY_METHOD_PICKUP && di.AutodetectDeliveryLocation {
		return "pickup_autodetect"
	}
	if di.Method == DELIVERY_METHOD_PICKUP {
		return di.Method + "_" + di.DeliveryLocationType + "_" + di.DeliveryLocationCode

	}
	return di.Method
}
