package domain

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

type (
	SourcingService interface {
		GetSourceId(ctx web.Context, decoratedCart *cart.DecoratedCart, item *cart.DecoratedCartItem) (string, error)
	}

	SourcingEngine struct {
		SourcingService SourcingService          `inject:",optional"`
		Logger          flamingo.Logger          `inject:""`
		Cartservice     *application.CartService `inject:""`
	}
)

// SetSourcesForCartItems gets Sources and modifies the Cart Items
// todo move to application layer ?
func (se *SourcingEngine) SetSourcesForCartItems(ctx context.Context, session *sessions.Session, decoratedCart *cart.DecoratedCart) error {
	if se.SourcingService == nil {
		return nil
	}
	for _, decoratedCartItem := range decoratedCart.DecoratedItems {
		sourceId, err := se.SourcingService.GetSourceId(web.ToContext(ctx), decoratedCart, &decoratedCartItem)
		if err != nil {
			se.Logger.WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Error(err)
			return fmt.Errorf("checkout.application.sourcingengine error: %v", err)
		}
		se.Logger.WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Debug("SourcingEngine detected source %v for item %v", sourceId, decoratedCartItem.Item.ID)
		err = se.Cartservice.UpdateItemSourceId(ctx, session, decoratedCartItem.Item.ID, sourceId)
		if err != nil {
			return errors.Wrap(err, "Could not update cart item")
		}
	}
	return nil
}
