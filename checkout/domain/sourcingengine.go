package domain

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/pkg/errors"
)

type (
	SourcingService interface {
		GetSourceId(ctx context.Context, session *web.Session, decoratedCart *cart.DecoratedCart, deliveryCode string, item *cart.DecoratedCartItem) (string, error)
	}

	SourcingEngine struct {
		SourcingService SourcingService          `inject:",optional"`
		Logger          flamingo.Logger          `inject:""`
		Cartservice     *application.CartService `inject:""`
	}
)

// SetSourcesForCartItems gets Sources and modifies the Cart Items
// todo move to application layer ?
func (se *SourcingEngine) SetSourcesForCartItems(ctx context.Context, session *web.Session, decoratedCart *cart.DecoratedCart) error {
	if se.SourcingService == nil {
		return nil
	}
	for _, decoratedDelivery := range decoratedCart.DecoratedDeliveries {
		for _, decoratedCartItem := range decoratedDelivery.DecoratedItems {
			sourceId, err := se.SourcingService.GetSourceId(ctx, session, decoratedCart, decoratedDelivery.Delivery.DeliveryInfo.Code, &decoratedCartItem)
			if err != nil {
				se.Logger.WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Error(err)
				return fmt.Errorf("checkout.application.sourcingengine error: %v", err)
			}
			se.Logger.WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Debug("SourcingEngine detected source %v for item %v", sourceId, decoratedCartItem.Item.ID)
			err = se.Cartservice.UpdateItemSourceID(ctx, session, decoratedCartItem.Item.ID, decoratedDelivery.Delivery.DeliveryInfo.Code, sourceId)
			if err != nil {
				return errors.Wrap(err, "Could not update cart item")
			}
		}
	}
	return nil
}
