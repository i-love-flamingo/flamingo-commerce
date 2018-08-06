package cart

import (
	"strings"

	"flamingo.me/flamingo/framework/flamingo"
)

type (

	//DeliveryInfoBuilder - Factory
	DeliveryInfoBuilder interface {
		BuildByDeliveryCode(deliveryCode string) (DeliveryInfo, error)
		//BuildDeliveryInfoUpdateCommand(ctx web.Context, decoratedCart *DecoratedCart) ([]DeliveryInfoUpdateCommand, error)
	}

	DefaultDeliveryInfoBuilder struct {
		Logger flamingo.Logger `inject:""`
	}
)

/*
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

*/

func (b *DefaultDeliveryInfoBuilder) BuildByDeliveryCode(deliverycode string) (DeliveryInfo, error) {
	if deliverycode == "" {
		b.Logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder").Warn("Empty deliverycode")
		return DeliveryInfo{
			Method: DELIVERY_METHOD_UNSPECIFIED,
		}, nil
	}
	if deliverycode == DELIVERY_METHOD_DELIVERY {
		return DeliveryInfo{
			Method: DELIVERY_METHOD_DELIVERY,
		}, nil
	}

	if deliverycode == "pickup_store" {
		return DeliveryInfo{
			Method: DELIVERY_METHOD_PICKUP,
			DeliveryLocation: DeliveryLocation{
				Type: DELIVERYLOCATION_TYPE_STORE,
			},
		}, nil
	}

	intentParts := strings.SplitN(deliverycode, "_", 3)
	if len(intentParts) != 3 {
		b.Logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder").Warn("Unknown deliverycode", deliverycode)
		return DeliveryInfo{
			Method: DELIVERY_METHOD_UNSPECIFIED,
		}, nil
	}
	if intentParts[0] == DELIVERY_METHOD_PICKUP {
		if intentParts[1] == DELIVERYLOCATION_TYPE_STORE {
			return DeliveryInfo{
				Method: DELIVERY_METHOD_PICKUP,
				DeliveryLocation: DeliveryLocation{
					Code: intentParts[2],
					Type: DELIVERYLOCATION_TYPE_STORE,
				},
			}, nil
		}
		if intentParts[1] == DELIVERYLOCATION_TYPE_COLLECTIONPOINT {
			return DeliveryInfo{
				Method: DELIVERY_METHOD_PICKUP,
				DeliveryLocation: DeliveryLocation{
					Code: intentParts[2],
					Type: DELIVERYLOCATION_TYPE_COLLECTIONPOINT,
				},
			}, nil
		}
	}
	b.Logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder").Warn("Unknown IntentString", deliverycode)
	return DeliveryInfo{
		Method: DELIVERY_METHOD_UNSPECIFIED,
	}, nil
}
