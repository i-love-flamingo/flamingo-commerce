package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// GetQuantityAdjustmentDeletedItemsMessages is exported as a template function
	GetQuantityAdjustmentDeletedItemsMessages struct{}

	// GetQuantityAdjustmentUpdatedItemsMessage is exported as a template function
	GetQuantityAdjustmentUpdatedItemsMessage struct{}

	// GetQuantityAdjustmentCouponCodesRemoved is exported as a template function
	GetQuantityAdjustmentCouponCodesRemoved struct{}

	// RemoveQuantityAdjustmentMessages is exported as a template function
	RemoveQuantityAdjustmentMessages struct{}

	// QuantityAdjustment is returned by the template function
	QuantityAdjustment struct {
		Item         cart.Item
		DeliveryCode string
		PrevQty      int
		CurrQty      int
		Reason       string
	}
)

// Func defines the GetQuantityAdjustmentDeletedItemsMessages template function
func (gdm *GetQuantityAdjustmentDeletedItemsMessages) Func(ctx context.Context) interface{} {
	return func() []QuantityAdjustment {
		session := web.SessionFromContext(ctx)

		deletedAdjustments := make([]QuantityAdjustment, 0)

		if sessionAdjustments, found := session.Load("cart.view.quantity.adjustments"); found {
			if adjustments, ok := sessionAdjustments.(application.QtyAdjustmentResults); ok {
				for _, a := range adjustments {
					if a.WasDeleted {
						deletedAdjustments = append(deletedAdjustments, QuantityAdjustment{
							Item:         a.OriginalItem,
							DeliveryCode: a.DeliveryCode,
							PrevQty:      a.NewQty - a.RestrictionResult.RemainingDifference,
							CurrQty:      a.NewQty,
							Reason:       a.RestrictionResult.RestrictorName,
						})
					}
				}
			}
		}

		return deletedAdjustments
	}
}

// Func defines the GetQuantityAdjustmentUpdatedItemsMessage template function
func (gum *GetQuantityAdjustmentUpdatedItemsMessage) Func(ctx context.Context) interface{} {
	return func(item cart.Item, deliveryCode string) QuantityAdjustment {
		session := web.SessionFromContext(ctx)

		if sessionAdjustments, found := session.Load("cart.view.quantity.adjustments"); found {
			if adjustments, ok := sessionAdjustments.(application.QtyAdjustmentResults); ok {
				for _, a := range adjustments {
					if a.OriginalItem.ID == item.ID && a.DeliveryCode == deliveryCode {
						return QuantityAdjustment{
							Item:         a.OriginalItem,
							DeliveryCode: a.DeliveryCode,
							PrevQty:      a.NewQty - a.RestrictionResult.RemainingDifference,
							CurrQty:      a.NewQty,
							Reason:       a.RestrictionResult.RestrictorName,
						}
					}
				}
			}
		}

		return QuantityAdjustment{
			Item:         item,
			DeliveryCode: deliveryCode,
			PrevQty:      item.Qty,
			CurrQty:      item.Qty,
		}
	}
}

// Func defines the GetQuantityAdjustmentCouponCodesRemoved template function
func (gcd *GetQuantityAdjustmentCouponCodesRemoved) Func(ctx context.Context) interface{} {
	return func() bool {
		session := web.SessionFromContext(ctx)

		if sessionAdjustments, found := session.Load("cart.view.quantity.adjustments"); found {
			if adjustments, ok := sessionAdjustments.(application.QtyAdjustmentResults); ok {
				return adjustments.HasRemovedCouponCodes()
			}
		}

		return false
	}
}

// Func defines the RemoveQuantityAdjustmentMessages template function
func (rm *RemoveQuantityAdjustmentMessages) Func(ctx context.Context) interface{} {
	return func() bool {
		session := web.SessionFromContext(ctx)

		session.Delete("cart.view.quantity.adjustments")

		return true
	}
}
