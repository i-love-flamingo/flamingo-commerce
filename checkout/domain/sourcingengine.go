package domain

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// SourcingService helps in retrieving item sources
	SourcingService interface {
		GetSourceID(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (string, error)
	}

	// SourcingEngine computes item sources
	SourcingEngine struct {
		SourcingService SourcingService          `inject:",optional"`
		Logger          flamingo.Logger          `inject:""`
		Cartservice     *application.CartService `inject:""`
	}
)

// SetSourcesForCartItems gets Sources and modifies the Cart Items
// todo move to application layer ?
func (se *SourcingEngine) SetSourcesForCartItems(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart) error {
	if se.SourcingService == nil {
		return nil
	}
	for _, decoratedDelivery := range decoratedCart.DecoratedDeliveries {
		for _, decoratedCartItem := range decoratedDelivery.DecoratedItems {
			sourceID, err := se.SourcingService.GetSourceID(ctx, session, decoratedCart, decoratedDelivery.Delivery.DeliveryInfo.Code, &decoratedCartItem)
			if err != nil {
				se.Logger.WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Error(err)
				return fmt.Errorf("checkout.application.sourcingengine error: %v", err)
			}
			se.Logger.WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Debug("SourcingEngine detected source %v for item %v", sourceID, decoratedCartItem.Item.ID)
			err = se.Cartservice.UpdateItemSourceID(ctx, session, decoratedCartItem.Item.ID, decoratedDelivery.Delivery.DeliveryInfo.Code, sourceID)
			if err != nil {
				return errors.Wrap(err, "Could not update cart item")
			}
		}
	}
	return nil
}
