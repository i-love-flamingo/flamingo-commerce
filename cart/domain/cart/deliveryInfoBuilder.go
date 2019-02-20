package cart

import (
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (

	//DeliveryInfoBuilder - Factory
	DeliveryInfoBuilder interface {
		BuildByDeliveryCode(deliveryCode string) (DeliveryInfo, error)
		//BuildDeliveryInfoUpdateCommand(ctx web.Context, decoratedCart *DecoratedCart) ([]DeliveryInfoUpdateCommand, error)
	}

	// DefaultDeliveryInfoBuilder defines the default delivery info builder used
	DefaultDeliveryInfoBuilder struct {
		logger flamingo.Logger
	}
)

// Inject dependencies
func (b *DefaultDeliveryInfoBuilder) Inject(
	logger flamingo.Logger,
) {
	b.logger = logger
}

// BuildByDeliveryCode builds a DeliveryInfo by deliveryCode
func (b *DefaultDeliveryInfoBuilder) BuildByDeliveryCode(deliverycode string) (DeliveryInfo, error) {
	if deliverycode == "" {
		b.logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder").Warn("Empty deliverycode")
		return DeliveryInfo{
			Code:   deliverycode,
			Method: DeliveryMethodUnspecified,
		}, nil
	}
	if deliverycode == DeliveryMethodDelivery {
		return DeliveryInfo{
			Code:   deliverycode,
			Method: DeliveryMethodDelivery,
		}, nil
	}

	if deliverycode == "pickup_store" {
		return DeliveryInfo{
			Code:   deliverycode,
			Method: DeliveryMethodPickup,
			DeliveryLocation: DeliveryLocation{
				Type: DeliverylocationTypeStore,
			},
		}, nil
	}

	intentParts := strings.SplitN(deliverycode, "_", 3)
	if len(intentParts) != 3 {
		b.logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder").Warn("Unknown deliverycode", deliverycode)
		return DeliveryInfo{
			Code:   deliverycode,
			Method: DeliveryMethodUnspecified,
		}, nil
	}
	if intentParts[0] == DeliveryMethodPickup || intentParts[0] == DeliveryMethodDelivery {
		if intentParts[1] == DeliverylocationTypeStore {
			return DeliveryInfo{
				Code:   deliverycode,
				Method: intentParts[0],
				DeliveryLocation: DeliveryLocation{
					Code: intentParts[2],
					Type: DeliverylocationTypeStore,
				},
			}, nil
		} else if intentParts[1] == DeliverylocationTypeCollectionpoint {
			return DeliveryInfo{
				Code:   deliverycode,
				Method: intentParts[0],
				DeliveryLocation: DeliveryLocation{
					Code: intentParts[2],
					Type: DeliverylocationTypeCollectionpoint,
				},
			}, nil
		} else {
			return DeliveryInfo{
				Code:   deliverycode,
				Method: intentParts[0],
				DeliveryLocation: DeliveryLocation{
					Code: intentParts[2],
					Type: intentParts[1],
				},
			}, nil
		}
	}
	b.logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder").Warn("Unknown IntentString", deliverycode)
	return DeliveryInfo{
		Code:   deliverycode,
		Method: DeliveryMethodUnspecified,
	}, nil
}
