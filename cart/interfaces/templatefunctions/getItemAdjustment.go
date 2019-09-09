package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// GetItemAdjustment is exported as a template function
	GetItemAdjustment struct {
		logger flamingo.Logger
	}

	// QuantityAdjustment is returned by the template function
	QuantityAdjustment struct {
		item    cart.Item
		prevQty int
		currQty int
	}
)

// Inject dependencies
func (tf *GetItemAdjustment) Inject(
	logger flamingo.Logger,

) {
	tf.logger = logger.WithField(flamingo.LogKeyModule, "cart").WithField(flamingo.LogKeyCategory, "getItemAdjustment")
}

// Func defines the GetItemAdjustment template function
func (tf *GetItemAdjustment) Func(ctx context.Context) interface{} {
	return func(item cart.Item, deliveryCode string) QuantityAdjustment {
		session := web.SessionFromContext(ctx)

		if adjustmentsI, found := session.Load("cart.view.adjustment.update"); found {
			if adjustments, ok := adjustmentsI.(application.QtyAdjustmentResults); ok {
				for _, a := range adjustments {
					if a.Item.ID == item.ID && a.DeliveryCode == deliveryCode {
						return QuantityAdjustment{
							item:    item,
							prevQty: item.Qty - a.RestrictionResult.RemainingDifference,
							currQty: item.Qty,
						}
					}
				}
			}
		}

		return QuantityAdjustment{
			item:    item,
			prevQty: item.Qty,
			currQty: item.Qty,
		}
	}
}
