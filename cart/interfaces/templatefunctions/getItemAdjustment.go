package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// GetItemAdjustment is exported as a template function
	GetItemAdjustment struct{}

	// QuantityAdjustment is returned by the template function
	QuantityAdjustment struct {
		Item    cart.Item
		PrevQty int
		CurrQty int
	}
)

// Func defines the GetItemAdjustment template function
func (gia *GetItemAdjustment) Func(ctx context.Context) interface{} {
	return func(item cart.Item, deliveryCode string) interface{} {
		session := web.SessionFromContext(ctx)

		if adjustmentsI, found := session.Load("cart.view.adjustment.update"); found {
			if adjustments, ok := adjustmentsI.(application.QtyAdjustmentResults); ok {
				for _, a := range adjustments {
					if a.Item.ID == item.ID && a.DeliveryCode == deliveryCode {
						return &QuantityAdjustment{
							Item:    item,
							PrevQty: item.Qty - a.RestrictionResult.RemainingDifference,
							CurrQty: item.Qty,
						}
					}
				}
			}
		}

		return &QuantityAdjustment{
			Item:    item,
			PrevQty: item.Qty,
			CurrQty: item.Qty,
		}
	}
}
